package web

import (
	"github.com/alexa-infra/git47/internal/core"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gorilla/mux"
	"net/http"
)

type diffViewData struct {
	Commit  *commitData
	Parents []*commitData
	TreeURL string
	Stats   object.FileStats
}

func diffHandler(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getRequestContext(r)

	vars := mux.Vars(r)
	hashStr := vars["hash"]
	ref, err := core.GetCommitRef(ctx.Ref.Repository, hashStr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	diff, err := core.GetDiff(ref)

	data := diffViewData{Commit: newCommitData(ctx, ref.Commit), Stats: diff}

	err = renderTemplate(w, "git-diff.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
