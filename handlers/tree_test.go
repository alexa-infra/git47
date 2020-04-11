package handlers

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTreeMaster(t *testing.T) {
	env := makeTestEnv(t)

	req, _ := http.NewRequest("GET", "/r/memory/tree/master", nil)
	req = mux.SetURLVars(req, map[string]string{"repo": "memory", "ref": "master"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(makeHandler(gitTree, env))

	handler.ServeHTTP(rr, req)

	if assert.Equal(t, rr.Code, http.StatusOK, rr.Body.String()) {
		assert.Contains(t, rr.Body.String(), `<a href="/r/memory/blob/e69de29bb2d1d6434b8b29ae775ad8c2e48c5391/foo">foo</a>`)
		assert.Contains(t, rr.Body.String(), `<a href="/r/memory/tree/master/bar">bar</a>`)
	}
}

func TestTreeMasterSubtree(t *testing.T) {
	env := makeTestEnv(t)

	req, _ := http.NewRequest("GET", "/r/memory/tree/master/bar", nil)
	req = mux.SetURLVars(req, map[string]string{"repo": "memory", "ref": "master"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(makeHandler(gitTree, env))

	handler.ServeHTTP(rr, req)

	if assert.Equal(t, rr.Code, http.StatusOK, rr.Body.String()) {
		assert.Contains(t, rr.Body.String(), `<a href="/r/memory/blob/5e1c309dae7f45e0f39b1bf3ac3cd9db12e7d689/bar/foo">foo</a>`)
	}
}

func TestTreeMasterSubtreeNotFound(t *testing.T) {
	env := makeTestEnv(t)

	req, _ := http.NewRequest("GET", "/r/memory/tree/master/foobar", nil)
	req = mux.SetURLVars(req, map[string]string{"repo": "memory", "ref": "master"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(makeHandler(gitTree, env))

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusNotFound, rr.Body.String())
}
