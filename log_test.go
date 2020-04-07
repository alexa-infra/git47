package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGitLog(t *testing.T) {
	r := prepareRepository(t)

	req, _ := http.NewRequest("GET", "/r/memory/commits/master", nil)

	router := mux.NewRouter()
	makeRoutes(router)

	ctx := req.Context()
	ctx = context.WithValue(ctx, gitRepoKey, r)
	ctx = context.WithValue(ctx, routerKey, router)
	req = req.WithContext(ctx)

	req = mux.SetURLVars(req, map[string]string{"repo": "memory", "ref": "master"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(gitLog)

	handler.ServeHTTP(rr, req)

	if assert.Equal(t, rr.Code, http.StatusOK, rr.Body.String()) {
		assert.Contains(t, rr.Body.String(), `<a href="/r/memory/commit/17a958a4b3f7f1aa265f782cf6e01e24cd4010cf">foo (17a958a4b3f7f1aa265f782cf6e01e24cd4010cf)</a>`)
		assert.Contains(t, rr.Body.String(), `<a href="/r/memory/commit/60a58ae38710f264b2c00f77c82ae44419381a3f">foobar (60a58ae38710f264b2c00f77c82ae44419381a3f)</a>`)
	}
}

func TestGitLogNext(t *testing.T) {
	r := prepareRepository(t)

	req, _ := http.NewRequest("GET", "/r/memory/commits/master?next=17a958a4b3f7f1aa265f782cf6e01e24cd4010cf", nil)

	router := mux.NewRouter()
	makeRoutes(router)

	ctx := req.Context()
	ctx = context.WithValue(ctx, gitRepoKey, r)
	ctx = context.WithValue(ctx, routerKey, router)
	req = req.WithContext(ctx)

	req = mux.SetURLVars(req, map[string]string{"repo": "memory", "ref": "master"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(gitLog)

	handler.ServeHTTP(rr, req)

	if assert.Equal(t, rr.Code, http.StatusOK, rr.Body.String()) {
		assert.Contains(t, rr.Body.String(), `<a href="/r/memory/commit/17a958a4b3f7f1aa265f782cf6e01e24cd4010cf">foo (17a958a4b3f7f1aa265f782cf6e01e24cd4010cf)</a>`)
		assert.NotContains(t, rr.Body.String(), `<a href="/r/memory/commit/60a58ae38710f264b2c00f77c82ae44419381a3f">foobar (60a58ae38710f264b2c00f77c82ae44419381a3f)</a>`)
	}
}

func TestGitLogNextNotFound(t *testing.T) {
	r := prepareRepository(t)

	req, _ := http.NewRequest("GET", "/r/memory/commits/master?next=blah", nil)

	router := mux.NewRouter()
	makeRoutes(router)

	ctx := req.Context()
	ctx = context.WithValue(ctx, gitRepoKey, r)
	ctx = context.WithValue(ctx, routerKey, router)
	req = req.WithContext(ctx)

	req = mux.SetURLVars(req, map[string]string{"repo": "memory", "ref": "master"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(gitLog)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusNotFound, rr.Body.String())
}
