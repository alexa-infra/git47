package handlers

import (
	"github.com/alexa-infra/git47/app/frontend/server"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"path"
	"strings"
	"html/template"
)

func NotImplemented(w http.ResponseWriter, r *http.Request) {
	err := fmt.Errorf("Not implemented (%s) %s", mux.CurrentRoute(r).GetName(), r.URL.Path)
	http.Error(w, err.Error(), http.StatusNotImplemented)
}

var (
	errRefNotSet    = errors.New("Ref not set")
	errRefNotFound  = errors.New("Ref not found")
	errRepoNotFound = errors.New("Repository not found")
	errBlobNotFound = errors.New("Blob not found")
	errInvalidHash  = errors.New("Invalid hash")
)


func parentPath(path string) string {
	if strings.Index(path, "/") > -1 {
		return path[:strings.LastIndex(path, "/")]
	}
	return ""
}

func joinURL(base string, paths ...string) string {
	p := path.Join(paths...)
	if p == "" {
		return strings.TrimRight(base, "/")
	}
	return fmt.Sprintf("%s/%s", strings.TrimRight(base, "/"), strings.TrimLeft(p, "/"))
}

func joinPath(paths ...string) string {
	return path.Join(paths...)
}

func RegisterHandlers(env *server.Env, r *mux.Router) {
	summary := GitSummary(env)
	tree := GitTree(env)
	blob := GitBlob(env)
	commits := GitLog(env)
	diff := GitDiff(env)

	r.Handle("/r/{repo}", summary).Name("summary")
	r.Handle("/r/{repo}/", summary).Name("summary2")
	r.Handle("/r/{repo}/summary/{ref}", summary).Name("summary_ref")
	r.Handle("/r/{repo}/summary/{ref}/", summary).Name("summary_ref2")
	r.PathPrefix("/r/{repo}/tree/{ref}").Handler(tree).Name("tree")
	r.PathPrefix("/r/{repo}/blob/{ref}").Handler(blob).Name("blob")
	r.HandleFunc("/r/{repo}/archive/{ref}.tar.gz", NotImplemented).Name("archive")
	r.Handle("/r/{repo}/commits/{ref}", commits).Name("commits")
	r.Handle("/r/{repo}/commits/{ref}/", commits).Name("commits2")
	r.Handle("/r/{repo}/commit/{hash}", diff).Name("commit")
	r.Handle("/r/{repo}/commit/{hash}/", diff).Name("commit2")
	r.HandleFunc("/r/{repo}/branches", NotImplemented).Name("branches")
	r.HandleFunc("/r/{repo}/tags", NotImplemented).Name("tags")
	r.HandleFunc("/r/{repo}/contributors", NotImplemented).Name("contributors")
}

func TemplateHelpers() template.FuncMap {
	return template.FuncMap{
		"GetSummaryURL": GetSummaryURL,
		"GetTreeURL": GetTreeURL,
		"GetLogURL": GetLogURL,
		"GetCommitURL": GetCommitURL,
		"GetBlobURL": GetBlobURL,
	}
}
