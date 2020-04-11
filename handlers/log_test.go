package handlers

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGitLog(t *testing.T) {
	env := makeTestEnv(t)

	req, _ := http.NewRequest("GET", "/r/memory/commits/master", nil)
	req = mux.SetURLVars(req, map[string]string{"repo": "memory", "ref": "master"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(makeHandler(gitLog, env))

	handler.ServeHTTP(rr, req)

	if assert.Equal(t, rr.Code, http.StatusOK, rr.Body.String()) {
		assert.Contains(t, rr.Body.String(), `<a href="/r/memory/commit/17a958a4b3f7f1aa265f782cf6e01e24cd4010cf">foo (17a958a4b3f7f1aa265f782cf6e01e24cd4010cf)</a>`)
		assert.Contains(t, rr.Body.String(), `<a href="/r/memory/commit/60a58ae38710f264b2c00f77c82ae44419381a3f">foobar (60a58ae38710f264b2c00f77c82ae44419381a3f)</a>`)
	}
}

func TestGitLogNext(t *testing.T) {
	env := makeTestEnv(t)

	req, _ := http.NewRequest("GET", "/r/memory/commits/master?next=17a958a4b3f7f1aa265f782cf6e01e24cd4010cf", nil)
	req = mux.SetURLVars(req, map[string]string{"repo": "memory", "ref": "master"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(makeHandler(gitLog, env))

	handler.ServeHTTP(rr, req)

	if assert.Equal(t, rr.Code, http.StatusOK, rr.Body.String()) {
		assert.Contains(t, rr.Body.String(), `<a href="/r/memory/commit/17a958a4b3f7f1aa265f782cf6e01e24cd4010cf">foo (17a958a4b3f7f1aa265f782cf6e01e24cd4010cf)</a>`)
		assert.NotContains(t, rr.Body.String(), `<a href="/r/memory/commit/60a58ae38710f264b2c00f77c82ae44419381a3f">foobar (60a58ae38710f264b2c00f77c82ae44419381a3f)</a>`)
	}
}

func TestGitLogNextNotFound(t *testing.T) {
	env := makeTestEnv(t)

	req, _ := http.NewRequest("GET", "/r/memory/commits/master?next=blah", nil)
	req = mux.SetURLVars(req, map[string]string{"repo": "memory", "ref": "master"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(makeHandler(gitLog, env))

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusBadRequest, rr.Body.String())
}
