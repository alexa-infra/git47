package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type key int

const (
	gitRepoKey key = iota
	routerKey
)

func makeHandler(fn func(http.ResponseWriter, *http.Request), router *mux.Router) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		g, err := getRepo(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, gitRepoKey, g)
		ctx = context.WithValue(ctx, routerKey, router)
		r = r.WithContext(ctx)
		fn(w, r)
	}
}

func getRepoVar(r *http.Request) *git.Repository {
	if rv := r.Context().Value(gitRepoKey); rv != nil {
		return rv.(*git.Repository)
	}
	return nil
}

func getRouter(r *http.Request) *mux.Router {
	if rv := r.Context().Value(routerKey); rv != nil {
		return rv.(*mux.Router)
	}
	return nil
}

func MakeRoutes(r *mux.Router) {
	handler := func(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
		return makeHandler(fn, r)
	}

	r.HandleFunc("/r/{repo}", handler(notImplemented)).Name("summary")
	r.PathPrefix("/r/{repo}/tree/{ref}").HandlerFunc(handler(gitTree)).Name("tree")
	r.PathPrefix("/r/{repo}/blob/{hash}").HandlerFunc(handler(gitBlob)).Name("blob")
	r.HandleFunc("/r/{repo}/archive/{ref}.tar.gz", handler(notImplemented)).Name("archive")
	r.HandleFunc("/r/{repo}/commits/{ref}", handler(gitLog)).Name("commits")
	r.HandleFunc("/r/{repo}/commit/{hash}", handler(gitDiff)).Name("commit")
	r.HandleFunc("/r/{repo}/contributors", handler(notImplemented)).Name("contributors")

	staticDir := "./static"
	if rootPath := os.Getenv("APPROOT"); rootPath != "" {
		staticDir = path.Join(rootPath, "static")
	}
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
}

func notImplemented(w http.ResponseWriter, r *http.Request) {
	log.Printf("Not implemented (%s) %s", mux.CurrentRoute(r).GetName(), r.URL.Path)
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

var repositories = map[string]string{
	"friday": "/home/alexey/projects/friday/.git",
	"git47":  "/home/alexey/projects/go-playground/git47/.git",
}

var (
	errRefNotFound  = errors.New("Ref not found")
	errRepoNotFound = errors.New("Repository not found")
	errBlobNotFound = errors.New("Blob not found")
)

func getRepo(r *http.Request) (*git.Repository, error) {
	vars := mux.Vars(r)
	repo := vars["repo"]
	path, ok := repositories[repo]

	if !ok {
		return nil, errRepoNotFound
	}

	g, err := git.PlainOpen(path)
	return g, err
}

func getRef(r *http.Request, g *git.Repository) (plumbing.Hash, error) {
	vars := mux.Vars(r)
	ref := vars["ref"]

	if ref == "" {
		return plumbing.ZeroHash, errRefNotFound
	}

	branch, err := g.Reference(plumbing.NewBranchReferenceName(ref), false)
	if err == nil {
		return branch.Hash(), nil
	}
	tag, err := g.Reference(plumbing.NewTagReferenceName(ref), false)
	if err == nil {
		return tag.Hash(), nil
	}

	hash := plumbing.NewHash(ref)
	if !hash.IsZero() {
		_, err := g.CommitObject(hash)
		if err == nil {
			return hash, nil
		}
	}

	return plumbing.ZeroHash, errRefNotFound
}

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

func parseTemplate(filenames ...string) *template.Template {
	if rootPath:= os.Getenv("APPROOT"); rootPath != "" {
		newFilenames := []string{}
		for _, filename := range filenames {
			newFilename := path.Join(rootPath, filename)
			newFilenames = append(newFilenames, newFilename)
		}
		filenames = newFilenames
	}

	return template.Must(template.ParseFiles(filenames...))
}
