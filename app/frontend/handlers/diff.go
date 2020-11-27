package handlers

import (
	"github.com/alexa-infra/git47/app/frontend/server"
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

// GitDiff returns handler which renders commit diff
func GitDiff(env *server.Env) http.HandlerFunc {
	template := env.GetTemplate("git-diff.html", TemplateHelpers())
	return env.WrapHandler(func (w http.ResponseWriter, r *http.Request) {
		ctx, _ := server.GetRequestContext(r)
		g := ctx.Repo

		vars := mux.Vars(r)
		hashStr := vars["hash"]

		hash := plumbing.NewHash(hashStr)
		if hash.IsZero() {
			http.Error(w, errInvalidHash.Error(), http.StatusBadRequest)
			return
		}

		commit, err := g.CommitObject(hash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		stats, err := commit.Stats()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := diffViewData{Commit: newCommitData(ctx, commit), Stats: stats}

		err = env.RenderTemplate(w, template, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

// GetCommitURL builds URL of commit diff
func GetCommitURL(rc *server.RequestContext, commit *object.Commit) (string, error) {
	router := rc.Env.Router
	route := router.Get("commit")
	url, err := route.URLPath("repo", rc.Config.Name, "hash", commit.Hash.String())
	if err != nil {
		return "", err
	}
	return url.Path, nil
}
