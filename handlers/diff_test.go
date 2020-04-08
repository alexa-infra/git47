package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGitDiff(t *testing.T) {
	r := prepareRepository(t)

	req, _ := http.NewRequest("GET", "/r/memory/commit/60a58ae38710f264b2c00f77c82ae44419381a3f", nil)

	router := mux.NewRouter()
	MakeRoutes(router)

	ctx := req.Context()
	ctx = context.WithValue(ctx, gitRepoKey, r)
	ctx = context.WithValue(ctx, routerKey, router)
	req = req.WithContext(ctx)

	req = mux.SetURLVars(req, map[string]string{"repo": "memory", "hash": "60a58ae38710f264b2c00f77c82ae44419381a3f"})

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(gitDiff)

	handler.ServeHTTP(rr, req)

	if assert.Equal(t, rr.Code, http.StatusOK, rr.Body.String()) {
		assert.Contains(t, rr.Body.String(), `bar/foo | 1`)
	}
}
