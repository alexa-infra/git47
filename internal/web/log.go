package web

import (
	"github.com/alexa-infra/git47/internal/core"
	"net/http"
)

type commitsViewData struct {
	Commits []commitData
	*RequestContext
}

// GitLog returns handler which renders a list of commits
func GitLog(w http.ResponseWriter, r *http.Request) {
	ctx, ok := GetRequestContext(r)
	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	query := r.URL.Query()
	next := query.Get("next")
	nextRef, err := core.GetCommitRef(ctx.Ref.Repository, next)
	if err != nil && next != "" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	commits, err := core.GetLog(ctx.Ref, nextRef)
	if err != nil {
		if err == core.ErrRefNotFound {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := commitsViewData{RequestContext: ctx}
	for _, c := range commits {
		data.Commits = append(data.Commits, *newCommitData(ctx, c))
	}

	err = RenderTemplate(w, "git-commits.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
