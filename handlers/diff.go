package handlers

import (
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gorilla/mux"
	"net/http"
)

type patchData struct {
	FileName string
	URL string
	Content string
}

type diffViewData struct {
	Commit *commitData
	Parents []*commitData
	TreeURL string
	Stats object.FileStats
}

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

	data := diffViewData{Commit: newCommitData(env, rc, commit), Stats: stats}

	template, err := env.Template.GetTemplate("git-diff.html")
	if err != nil {
		return err
	}

	return template.ExecuteTemplate(w, "layout", data)
}
