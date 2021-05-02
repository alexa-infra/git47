package web

import (
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGitBlob(t *testing.T) {
	req, _ := http.NewRequest("GET", "/r/memory/blob/60a58ae38710f264b2c00f77c82ae44419381a3f/bar/foo", nil)
	req = mux.SetURLVars(req, map[string]string{"repo": "memory", "ref": "60a58ae38710f264b2c00f77c82ae44419381a3f"})

	rr := httptest.NewRecorder()
	router := makeTestRouter(t)
	router.ServeHTTP(rr, req)

	if assert.Equal(t, http.StatusOK, rr.Code, rr.Body.String()) {
		assert.Contains(t, rr.Body.String(), `Hello World`)
	}
}

func TestGitBlobNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/r/memory/blob/5e1c309dae7f45e0f39b1bf3ac3cd9db12e7d689/bar/foo", nil)
	req = mux.SetURLVars(req, map[string]string{"repo": "memory", "ref": "5e1c309dae7f45e0f39b1bf3ac3cd9db12e7d689"})

	rr := httptest.NewRecorder()
	router := makeTestRouter(t)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code, rr.Body.String())
}

func TestGitBlobDir(t *testing.T) {
	req, _ := http.NewRequest("GET", "/r/memory/blob/60a58ae38710f264b2c00f77c82ae44419381a3f/bar", nil)
	req = mux.SetURLVars(req, map[string]string{"repo": "memory", "ref": "60a58ae38710f264b2c00f77c82ae44419381a3f"})

	rr := httptest.NewRecorder()
	router := makeTestRouter(t)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code, rr.Body.String())
}
