package web

import (
	"context"
	"errors"
	"github.com/alexa-infra/git47/internal/core"
	"github.com/go-git/go-git/v5"
	"github.com/gorilla/mux"
	"net/http"
)

type key int

var reqContextKey key

var (
	errRepoNotFound = errors.New("Repository not found")
	errRefNotSet    = errors.New("Ref not set")
)

type requestContext struct {
	Ref        core.NamedReference
	Router     *mux.Router
	RepoConfig core.RepoConfig
}

func getRequestContext(r *http.Request) (*requestContext, bool) {
	ctx := r.Context()
	reqCtx, ok := ctx.Value(reqContextKey).(*requestContext)
	return reqCtx, ok
}

func withRequestContext(r *http.Request, reqCtx *requestContext) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, reqContextKey, reqCtx)
	return r.WithContext(ctx)
}

func getRepoConfig(repositories core.RepoMap, r *http.Request) (core.RepoConfig, error) {
	vars := mux.Vars(r)
	repo := vars["repo"]
	cfg, ok := repositories[repo]

	if !ok {
		return core.RepoConfig{}, errRepoNotFound
	}
	return cfg, nil
}

func getNamedRef(g *git.Repository, r *http.Request) (core.NamedReference, error) {
	vars := mux.Vars(r)
	ref := vars["ref"]

	if ref == "" {
		return core.NamedReference{}, errRefNotSet
	}

	return core.GetNamedRef(g, ref)
}

func setVariables(router *mux.Router, repositories core.RepoMap, r *http.Request) (*http.Request, error) {
	rc, err := getRepoConfig(repositories, r)
	if err != nil {
		return nil, err
	}
	g, err := rc.Open()
	if err != nil {
		return nil, err
	}
	ref, err := getNamedRef(g, r)
	if err != nil {
		if err != errRefNotSet {
			return nil, err
		}

		ref, err = core.GetDefaultBranch(g)
		if err != nil {
			return nil, err
		}
	}
	reqCtx := &requestContext{
		Router:     router,
		RepoConfig: rc,
		Ref:        ref,
	}
	newReq := withRequestContext(r, reqCtx)
	return newReq, nil
}

func wrapHandler(router *mux.Router, repositories core.RepoMap) func(func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
		handler := http.HandlerFunc(fn)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r, err := setVariables(router, repositories, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			handler.ServeHTTP(w, r)
		})
	}
}
