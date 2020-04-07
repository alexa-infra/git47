package main

import (
	"context"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func defaultSignature() *object.Signature {
	when, _ := time.Parse(object.DateFormat, "Thu May 04 00:03:43 2017 +0200")
	return &object.Signature{
		Name:  "foo",
		Email: "foo@foo.foo",
		When:  when,
	}
}

func prepareRepository(t *testing.T) *git.Repository {
	storage := memory.NewStorage()
	fs := memfs.New()
	r, err := git.Init(storage, fs)
	require.Nil(t, err)

	err = util.WriteFile(fs, "foo", nil, 0755)
	require.Nil(t, err)

	w, err := r.Worktree()
	require.Nil(t, err)

	_, err = w.Add("foo")
	require.Nil(t, err)

	hash, err := w.Commit("foo", &git.CommitOptions{
		Author:    defaultSignature(),
		Committer: defaultSignature(),
	})
	require.Nil(t, err)
	require.Equal(t, "17a958a4b3f7f1aa265f782cf6e01e24cd4010cf", hash.String())

	fs.MkdirAll("bar", 0755)
	err = util.WriteFile(fs, "bar/foo", []byte("Hello World"), 0755)
	require.Nil(t, err)

	_, err = w.Add("bar/foo")
	require.Nil(t, err)

	hash, err = w.Commit("foobar", &git.CommitOptions{
		Author:    defaultSignature(),
		Committer: defaultSignature(),
	})
	require.Nil(t, err)
	require.Equal(t, "60a58ae38710f264b2c00f77c82ae44419381a3f", hash.String())

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName("new-branch"),
		Create: true,
	})
	require.Nil(t, err)

	err = util.WriteFile(fs, "foo", []byte("Test"), 0755)
	require.Nil(t, err)

	_, err = w.Add("foo")
	require.Nil(t, err)

	hash, err = w.Commit("foo", &git.CommitOptions{
		Author:    defaultSignature(),
		Committer: defaultSignature(),
	})
	require.Nil(t, err)
	require.Equal(t, "377229569f4a7ae706ed3a376117dabee4cec8f8", hash.String())

	_, err = r.CreateTag("v1", hash, nil)
	require.Nil(t, err)

	return r
}

func TestTreeMaster(t *testing.T) {
	r := prepareRepository(t)

	req, _ := http.NewRequest("GET", "/r/memory/tree/master", nil)

	router := mux.NewRouter()
	makeRoutes(router)

	ctx := req.Context()
	ctx = context.WithValue(ctx, gitRepoKey, r)
	ctx = context.WithValue(ctx, routerKey, router)
	req = req.WithContext(ctx)

	req = mux.SetURLVars(req, map[string]string{"repo": "memory", "ref": "master"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(gitTree)

	handler.ServeHTTP(rr, req)

	if assert.Equal(t, rr.Code, http.StatusOK, rr.Body.String()) {
		assert.Contains(t, rr.Body.String(), `<a href="/r/memory/blob/e69de29bb2d1d6434b8b29ae775ad8c2e48c5391/foo">foo</a>`)
		assert.Contains(t, rr.Body.String(), `<a href="/r/memory/tree/master/bar">bar</a>`)
	}
}

func TestTreeMasterSubtree(t *testing.T) {
	r := prepareRepository(t)

	req, _ := http.NewRequest("GET", "/r/memory/tree/master/bar", nil)

	router := mux.NewRouter()
	makeRoutes(router)

	ctx := req.Context()
	ctx = context.WithValue(ctx, gitRepoKey, r)
	ctx = context.WithValue(ctx, routerKey, router)
	req = req.WithContext(ctx)

	req = mux.SetURLVars(req, map[string]string{"repo": "memory", "ref": "master"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(gitTree)

	handler.ServeHTTP(rr, req)

	if assert.Equal(t, rr.Code, http.StatusOK, rr.Body.String()) {
		assert.Contains(t, rr.Body.String(), `<a href="/r/memory/blob/5e1c309dae7f45e0f39b1bf3ac3cd9db12e7d689/bar/foo">foo</a>`)
	}
}

func TestTreeMasterSubtreeNotFound(t *testing.T) {
	r := prepareRepository(t)

	req, _ := http.NewRequest("GET", "/r/memory/tree/master/foobar", nil)

	router := mux.NewRouter()
	makeRoutes(router)

	ctx := req.Context()
	ctx = context.WithValue(ctx, gitRepoKey, r)
	ctx = context.WithValue(ctx, routerKey, router)
	req = req.WithContext(ctx)

	req = mux.SetURLVars(req, map[string]string{"repo": "memory", "ref": "master"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(gitTree)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusNotFound, rr.Body.String())
}
