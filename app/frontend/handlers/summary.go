package handlers

import (
	"github.com/alexa-infra/git47/app/frontend/server"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"net/http"
)

type summaryViewData struct {
	NumCommits      int
	NumBranches     int
	NumTags         int
	NumFiles        int
	NumContributors int
	*server.RequestContext
}

// GitSummary returns handler which renders summary page of a repository
func GitSummary(env *server.Env) http.HandlerFunc {
	template := env.GetTemplate("git-summary.html", TemplateHelpers())
	return env.WrapHandler(func(w http.ResponseWriter, r *http.Request){
		reqCtx, _ := server.GetRequestContext(r)
		g := reqCtx.Repo
		ref := reqCtx.Ref

		data := summaryViewData{
			RequestContext: reqCtx,
		}

		cIter, err := g.Log(&git.LogOptions{From: ref.Hash()})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		uniqUsers := make(map[string]bool)
		cIter.ForEach(func(c *object.Commit) error {
			data.NumCommits++
			uniqUsers[c.Author.Email] = true
			return nil
		})
		data.NumContributors = len(uniqUsers)

		refs, err := g.References()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		refs.ForEach(func(ref *plumbing.Reference) error {
			refName := ref.Name()
			if refName.IsBranch() {
				data.NumBranches++
			}
			if refName.IsTag() {
				data.NumTags++
			}
			return nil
		})

		tree, err := ref.Commit.Tree()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tree.Files().ForEach(func(f *object.File) error {
			data.NumFiles++
			return nil
		})
		err = env.RenderTemplate(w, template, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

// GetSummaryURL builds URL of summary page
func GetSummaryURL(rc *server.RequestContext) (string, error) {
	router := rc.Env.Router
	route := router.Get("summary")
	url, err := route.URLPath("repo", rc.Config.Name)
	if err != nil {
		return "", err
	}
	return url.Path, nil
}