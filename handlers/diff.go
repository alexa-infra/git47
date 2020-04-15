package handlers

import (
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/gorilla/mux"
	"net/http"
)

func gitDiff(env *Env, w http.ResponseWriter, r *http.Request) error {
	rc, err := env.getRepoConfig(r)
	if err != nil {
		return StatusError{http.StatusNotFound, err}
	}

	g, err := rc.open()
	if err != nil {
		return StatusError{http.StatusNotFound, err}
	}

	vars := mux.Vars(r)
	hashStr := vars["hash"]

	hash := plumbing.NewHash(hashStr)
	if hash.IsZero() {
		return StatusError{http.StatusBadRequest, errInvalidHash}
	}

	commit, err := g.CommitObject(hash)
	if err != nil {
		return err
	}

	stats, err := commit.Stats()
	if err != nil {
		return err
	}

	template, err := env.Template.GetTemplate("git-diff.html")
	if err != nil {
		return err
	}

	return template.ExecuteTemplate(w, "layout", stats)
}
