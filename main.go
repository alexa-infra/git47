package main

import (
	"context"
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
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
		ctx := context.WithValue(r.Context(), gitRepoKey, g)
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

func main() {
	r := mux.NewRouter()
	s := r.PathPrefix("/r").Subrouter()

	handler := func(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
		return makeHandler(fn, s)
	}

	s.HandleFunc("/{repo}", handler(notImplemented)).Name("summary")
	s.PathPrefix("/{repo}/tree/{ref}").HandlerFunc(handler(gitTree)).Name("tree")
	s.PathPrefix("/{repo}/blob/{hash}").HandlerFunc(handler(gitBlob)).Name("blob")
	s.HandleFunc("/{repo}/archive/{ref}.tar.gz", handler(notImplemented)).Name("archive")
	s.HandleFunc("/{repo}/commits/{ref}", handler(gitLog)).Name("commits")
	s.HandleFunc("/{repo}/commit/{hash}", handler(gitDiff)).Name("commit")
	s.HandleFunc("/{repo}/contributors", handler(notImplemented)).Name("contributors")

	// Start server
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":1323", nil))
}

var templates = template.Must(template.ParseGlob("views/*.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) error {
	return templates.ExecuteTemplate(w, tmpl, data)
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
