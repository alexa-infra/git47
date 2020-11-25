package main

import (
	"github.com/alexa-infra/git47/app/frontend/server"
	"github.com/alexa-infra/git47/app/frontend/handlers"
	"github.com/alexa-infra/git47/app/frontend/middleware"
	"net/http"
	"path"
	"strings"
	"fmt"
)

func BuildPipeline(env *server.Env) {
	r := env.Router
	r.Use(middleware.Logging)
	setVars := middleware.Variables(env)

	env.Helpers["GetSummaryURL"] = handlers.GetSummaryURL
	env.Helpers["GetTreeURL"] = handlers.GetTreeURL
	env.Helpers["GetLogURL"] = handlers.GetLogURL
	env.Helpers["GetCommitURL"] = handlers.GetCommitURL
	env.Helpers["GetBlobURL"] = handlers.GetBlobURL

	summary := setVars(handlers.GitSummary(env))
	tree := setVars(handlers.GitTree(env))
	blob := setVars(handlers.GitBlob())
	commits := setVars(handlers.GitLog(env))
	diff := setVars(handlers.GitDiff(env))

	r.Handle("/r/{repo}", summary).Name("summary")
	r.Handle("/r/{repo}/", summary).Name("summary2")
	r.Handle("/r/{repo}/summary/{ref}", summary).Name("summary_ref")
	r.Handle("/r/{repo}/summary/{ref}/", summary).Name("summary_ref2")
	r.PathPrefix("/r/{repo}/tree/{ref}").Handler(tree).Name("tree")
	r.PathPrefix("/r/{repo}/blob/{ref}").Handler(blob).Name("blob")
	r.HandleFunc("/r/{repo}/archive/{ref}.tar.gz", handlers.NotImplemented).Name("archive")
	r.Handle("/r/{repo}/commits/{ref}", commits).Name("commits")
	r.Handle("/r/{repo}/commits/{ref}/", commits).Name("commits2")
	r.Handle("/r/{repo}/commit/{hash}", diff).Name("commit")
	r.Handle("/r/{repo}/commit/{hash}/", diff).Name("commit2")
	r.HandleFunc("/r/{repo}/branches", handlers.NotImplemented).Name("branches")
	r.HandleFunc("/r/{repo}/tags", handlers.NotImplemented).Name("tags")
	r.HandleFunc("/r/{repo}/contributors", handlers.NotImplemented).Name("contributors")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(env.StaticPath))))

}

func joinURL(base string, paths ...string) string {
	p := path.Join(paths...)
	if p == "" {
		return strings.TrimRight(base, "/")
	}
	return fmt.Sprintf("%s/%s", strings.TrimRight(base, "/"), strings.TrimLeft(p, "/"))
}
