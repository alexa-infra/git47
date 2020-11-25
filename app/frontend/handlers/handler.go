package handlers

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

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

func wrapper(fn func(*Context) error, env *Env, w http.ResponseWriter, r *http.Request) error {
	rc, err := env.getRepoConfig(r)
	if err != nil {
		return StatusError{http.StatusNotFound, err}
	}

	g, err := rc.open()
	if err != nil {
		return StatusError{http.StatusNotFound, err}
	}

	ref, err := getNamedRef(g, r)
	if err != nil {
		if err != errRefNotSet {
			return StatusError{http.StatusNotFound, err}
		}

		ref, err = getDefaultBranch(g)
		if err != nil {
			return StatusError{http.StatusNotFound, err}
		}
	}
	ctx := &Context{
		Config:   rc,
		Ref:      ref,
		Env:      env,
		response: w,
		request:  r,
		repo:     g,
	}
	return fn(ctx)
}

func NewHandler(fn func(*Context) error, env *Env, r *mux.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := wrapper(fn, env, w, r)
		if err != nil {
			status := http.StatusInternalServerError
			switch e := err.(type) {
			case Error:
				status = e.Status()
			}
			log.Printf("HTTP %d - %s", status, err)
			http.Error(w, err.Error(), status)
		}
	}
}
