package web

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGitLog(t *testing.T) {
	req, _ := http.NewRequest("GET", "/r/memory/commits/master", nil)
	req = mux.SetURLVars(req, map[string]string{"repo": "memory", "ref": "master"})

	rr := httptest.NewRecorder()
	router := makeTestRouter(t)
	router.ServeHTTP(rr, req)

	if assert.Equal(t, http.StatusOK, rr.Code, rr.Body.String()) {
		assert.Contains(t, rr.Body.String(), `<a href="/r/memory/commit/17a958a4b3f7f1aa265f782cf6e01e24cd4010cf">foo (17a958a)</a>`)
		assert.Contains(t, rr.Body.String(), `<a href="/r/memory/commit/60a58ae38710f264b2c00f77c82ae44419381a3f">foobar (60a58ae)</a>`)
	}
}

func TestGitLogNext(t *testing.T) {
	req, _ := http.NewRequest("GET", "/r/memory/commits/master?next=17a958a4b3f7f1aa265f782cf6e01e24cd4010cf", nil)
	req = mux.SetURLVars(req, map[string]string{"repo": "memory", "ref": "master"})

	rr := httptest.NewRecorder()
	router := makeTestRouter(t)
	router.ServeHTTP(rr, req)

	if assert.Equal(t, http.StatusOK, rr.Code, rr.Body.String()) {
		assert.Contains(t, rr.Body.String(), `<a href="/r/memory/commit/17a958a4b3f7f1aa265f782cf6e01e24cd4010cf">foo (17a958a)</a>`)
		assert.NotContains(t, rr.Body.String(), `<a href="/r/memory/commit/60a58ae38710f264b2c00f77c82ae44419381a3f">foobar (60a58ae)</a>`)
	}
}

func TestGitLogNextNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/r/memory/commits/master?next=blah", nil)
	req = mux.SetURLVars(req, map[string]string{"repo": "memory", "ref": "master"})

	rr := httptest.NewRecorder()
	router := makeTestRouter(t)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, rr.Body.String())
}
