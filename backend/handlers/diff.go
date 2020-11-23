package handlers

import (
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gorilla/mux"
	"net/http"
)

type patchData struct {
	FileName string
	URL      string
	Content  string
}

type diffViewData struct {
	Commit  *commitData
	Parents []*commitData
	TreeURL string
	Stats   object.FileStats
}

func gitDiff(ctx *Context) error {
	g := ctx.repo
	r := ctx.request

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

	data := diffViewData{Commit: newCommitData(ctx, commit), Stats: stats}

	return ctx.RenderTemplate("git-diff.html", data)
}
