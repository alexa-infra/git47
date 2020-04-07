package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGitBlob(t *testing.T) {
	r := prepareRepository(t)

	req, _ := http.NewRequest("GET", "/r/memory/blob/5e1c309dae7f45e0f39b1bf3ac3cd9db12e7d689/foo/bar", nil)

	router := mux.NewRouter()
	makeRoutes(router)

	ctx := req.Context()
	ctx = context.WithValue(ctx, gitRepoKey, r)
	ctx = context.WithValue(ctx, routerKey, router)
	req = req.WithContext(ctx)

	req = mux.SetURLVars(req, map[string]string{"repo": "memory", "hash": "5e1c309dae7f45e0f39b1bf3ac3cd9db12e7d689"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(gitBlob)

	handler.ServeHTTP(rr, req)

	if assert.Equal(t, rr.Code, http.StatusOK, rr.Body.String()) {
		assert.Contains(t, rr.Body.String(), `Hello World`)
	}
}

func TestGitBlobNotFound(t *testing.T) {
	r := prepareRepository(t)

	req, _ := http.NewRequest("GET", "/r/memory/blob/blah/foo/bar", nil)

	router := mux.NewRouter()
	makeRoutes(router)

	ctx := req.Context()
	ctx = context.WithValue(ctx, gitRepoKey, r)
	ctx = context.WithValue(ctx, routerKey, router)
	req = req.WithContext(ctx)

	req = mux.SetURLVars(req, map[string]string{"repo": "memory", "hash": "blah"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(gitBlob)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusNotFound, rr.Body.String())
}
