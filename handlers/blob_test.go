package handlers

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGitBlob(t *testing.T) {
	env := makeTestEnv(t)

	req, _ := http.NewRequest("GET", "/r/memory/blob/5e1c309dae7f45e0f39b1bf3ac3cd9db12e7d689/foo/bar", nil)
	req = mux.SetURLVars(req, map[string]string{"repo": "memory", "hash": "5e1c309dae7f45e0f39b1bf3ac3cd9db12e7d689"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(makeHandler(gitBlob, env))

	handler.ServeHTTP(rr, req)

	if assert.Equal(t, rr.Code, http.StatusOK, rr.Body.String()) {
		assert.Contains(t, rr.Body.String(), `Hello World`)
	}
}

func TestGitBlobNotFound(t *testing.T) {
	env := makeTestEnv(t)

	req, _ := http.NewRequest("GET", "/r/memory/blob/blah/foo/bar", nil)
	req = mux.SetURLVars(req, map[string]string{"repo": "memory", "hash": "blah"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(makeHandler(gitBlob, env))

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusBadRequest, rr.Body.String())
}
