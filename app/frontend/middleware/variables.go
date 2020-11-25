package middleware

import (
	"github.com/alexa-infra/git47/app/frontend/server"
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

type NamedReference struct {
	Name   string
	Kind   string
	Commit *object.Commit
}

func (ref *NamedReference) Hash() plumbing.Hash {
	return ref.Commit.Hash
}

func getRepoConfig(env *server.Env, r *http.Request) (*server.RepoConfig, error) {
	vars := mux.Vars(r)
	repo := vars["repo"]
	cfg, ok := env.Repositories[repo]

	if !ok {
		return nil, errRepoNotFound
	}
	return cfg, nil
}

func getNamedRef(g *git.Repository, r *http.Request) (*NamedReference, error) {
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

		return &NamedReference{
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

		return &NamedReference{
			Name:   ref,
			Kind:   "tag",
			Commit: commit,
		}, nil
	}

	hash := plumbing.NewHash(ref)
	if !hash.IsZero() {
		commit, err := g.CommitObject(hash)
		if err == nil {
			return &NamedReference{
				Name:   ref,
				Kind:   "commit",
				Commit: commit,
			}, nil
		}
	}

	return nil, errRefNotFound
}

func getDefaultBranch(g *git.Repository) (*NamedReference, error) {
	branch, err := g.Reference(plumbing.NewBranchReferenceName("master"), false)
	if err != nil {
		return nil, err
	}

	commit, err := g.CommitObject(branch.Hash())
	if err != nil {
		return nil, err
	}

	ref := &NamedReference{
		Name:   "master",
		Kind:   "branch",
		Commit: commit,
	}
	return ref, nil
}

type RequestContext struct {
	Config *server.RepoConfig
	Repo *git.Repository
	Ref *NamedReference
	Env *server.Env
}

func Variables(env *server.Env) (func(http.Handler) http.Handler) {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rc, err := getRepoConfig(env, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			g, err := rc.Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ref, err := getNamedRef(g, r)
			if err != nil {
				if err != errRefNotSet {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}

				ref, err = getDefaultBranch(g)
				if err != nil {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
			}
			reqCtx := &RequestContext{
				Config: rc,
				Repo: g,
				Ref: ref,
				Env: env,
			}
			ctx := context.WithValue(r.Context(), reqContextKey, reqCtx)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func GetRequestContext(r *http.Request) (*RequestContext, bool) {
	ctx := r.Context()
	reqCtx, ok := ctx.Value(reqContextKey).(*RequestContext)
	return reqCtx, ok
}
