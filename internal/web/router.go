package web

import (
	"embed"
	"github.com/alexa-infra/git47/internal/core"
	"github.com/alexa-infra/git47/internal/web/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

//go:embed static/* static/css/* static/webfonts/*
var static embed.FS

func NewRouter(cfg *Config, repositories core.RepoMap) (*mux.Router, error) {
	r := mux.NewRouter()
	if cfg.Logging {
		r.Use(middleware.Logging)
	}
	s := r.PathPrefix("/r/").Subrouter()

	wrapper := WrapHandler(r, repositories)
	blob := wrapper(GitBlob)
	summary := wrapper(GitSummary)
	commits := wrapper(GitLog)
	diff := wrapper(GitDiff)
	tree := wrapper(GitTree)

	s.Handle("/{repo}", summary).Name("summary")
	s.Handle("/{repo}/", summary).Name("summary2")
	s.Handle("/{repo}/summary/{ref}", summary).Name("summary_ref")
	s.Handle("/{repo}/summary/{ref}/", summary).Name("summary_ref2")
	s.PathPrefix("/{repo}/tree/{ref}").Handler(tree).Name("tree")
	s.PathPrefix("/{repo}/blob/{ref}").Handler(blob).Name("blob")
	s.HandleFunc("/{repo}/archive/{ref}.tar.gz", NotImplemented).Name("archive")
	s.Handle("/{repo}/commits/{ref}", commits).Name("commits")
	s.Handle("/{repo}/commits/{ref}/", commits).Name("commits2")
	s.Handle("/{repo}/commit/{hash}", diff).Name("commit")
	s.Handle("/{repo}/commit/{hash}/", diff).Name("commit2")
	s.HandleFunc("/{repo}/branches", NotImplemented).Name("branches")
	s.HandleFunc("/{repo}/tags", NotImplemented).Name("tags")
	s.HandleFunc("/{repo}/contributors", NotImplemented).Name("contributors")

	staticHandler := http.FileServer(http.FS(static))
	r.PathPrefix("/static/").Handler(staticHandler)
	return r, nil
}
