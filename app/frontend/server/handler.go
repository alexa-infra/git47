package server

import (
	"github.com/gorilla/mux"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"net/http"
	"errors"
	"context"
)

type key int
var reqContextKey key

var (
	errRefNotSet    = errors.New("Ref not set")
	errRefNotFound  = errors.New("Ref not found")
	errRepoNotFound = errors.New("Repository not found")
	errBlobNotFound = errors.New("Blob not found")
	errInvalidHash  = errors.New("Invalid hash")
)

type namedReference struct {
	Name   string
	Kind   string
	Commit *object.Commit
}

func (ref *namedReference) Hash() plumbing.Hash {
	return ref.Commit.Hash
}

func getRepoConfig(env *Env, r *http.Request) (*repoConfig, error) {
	vars := mux.Vars(r)
	repo := vars["repo"]
	cfg, ok := env.repositories[repo]

	if !ok {
		return nil, errRepoNotFound
	}
	return cfg, nil
}

func getNamedRef(g *git.Repository, r *http.Request) (*namedReference, error) {
	vars := mux.Vars(r)
	ref := vars["ref"]

	if ref == "" {
		return nil, errRefNotSet
	}

	branch, err := g.Reference(plumbing.NewBranchReferenceName(ref), false)
	if err == nil {
		hash := branch.Hash()

		commit, err := g.CommitObject(hash)
		if err != nil {
			return nil, err
		}

		return &namedReference{
			Name:   ref,
			Kind:   "branch",
			Commit: commit,
		}, nil
	}
	tag, err := g.Reference(plumbing.NewTagReferenceName(ref), false)
	if err == nil {
		hash := tag.Hash()

		commit, err := g.CommitObject(hash)
		if err != nil {
			return nil, err
		}

		return &namedReference{
			Name:   ref,
			Kind:   "tag",
			Commit: commit,
		}, nil
	}

	hash := plumbing.NewHash(ref)
	if !hash.IsZero() {
		commit, err := g.CommitObject(hash)
		if err == nil {
			return &namedReference{
				Name:   ref,
				Kind:   "commit",
				Commit: commit,
			}, nil
		}
	}

	return nil, errRefNotFound
}

func getDefaultBranch(g *git.Repository) (*namedReference, error) {
	branch, err := g.Reference(plumbing.NewBranchReferenceName("master"), false)
	if err != nil {
		return nil, err
	}

	commit, err := g.CommitObject(branch.Hash())
	if err != nil {
		return nil, err
	}

	ref := &namedReference{
		Name:   "master",
		Kind:   "branch",
		Commit: commit,
	}
	return ref, nil
}

// RequestContext carries current repository/branch extracted from request URL
type RequestContext struct {
	Config *repoConfig
	Repo *git.Repository
	Ref *namedReference
	Env *Env
}

// WrapHandler is used to add RequestContext to current request
func (env *Env) WrapHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	handler := http.HandlerFunc(fn)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r, err := setVariables(env, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func setVariables(env *Env, r *http.Request) (*http.Request, error) {
	rc, err := getRepoConfig(env, r)
	if err != nil {
		return nil, err
	}
	g, err := rc.open()
	if err != nil {
		return nil, err
	}
	ref, err := getNamedRef(g, r)
	if err != nil {
		if err != errRefNotSet {
			return nil, err
		}

		ref, err = getDefaultBranch(g)
		if err != nil {
			return nil, err
		}
	}
	reqCtx := &RequestContext{
		Config: rc,
		Repo: g,
		Ref: ref,
		Env: env,
	}
	ctx := context.WithValue(r.Context(), reqContextKey, reqCtx)
	return r.WithContext(ctx), nil
}

// GetRequestContext returns currently presented RequestContext
func GetRequestContext(r *http.Request) (*RequestContext, bool) {
	ctx := r.Context()
	reqCtx, ok := ctx.Value(reqContextKey).(*RequestContext)
	return reqCtx, ok
}
