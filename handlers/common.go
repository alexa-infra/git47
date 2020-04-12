package handlers

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/gorilla/mux"
	"net/http"
	"path"
	"strings"
)

func (env *Env) Setup() {
	r := env.Router

	handler := func(fn func(*Env, http.ResponseWriter, *http.Request) error) http.HandlerFunc {
		return makeHandler(fn, env)
	}

	r.HandleFunc("/r/{repo}", handler(notImplemented)).Name("summary")
	r.PathPrefix("/r/{repo}/tree/{ref}").HandlerFunc(handler(gitTree)).Name("tree")
	r.PathPrefix("/r/{repo}/blob/{ref}").HandlerFunc(handler(gitBlob)).Name("blob")
	r.HandleFunc("/r/{repo}/archive/{ref}.tar.gz", handler(notImplemented)).Name("archive")
	r.HandleFunc("/r/{repo}/commits/{ref}", handler(gitLog)).Name("commits")
	r.HandleFunc("/r/{repo}/commit/{hash}", handler(gitDiff)).Name("commit")
	r.HandleFunc("/r/{repo}/contributors", handler(notImplemented)).Name("contributors")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(env.Static.Path))))

	env.Template.Setup()
}

func notImplemented(env *Env, w http.ResponseWriter, r *http.Request) error {
	err := fmt.Errorf("Not implemented (%s) %s", mux.CurrentRoute(r).GetName(), r.URL.Path)
	return StatusError{http.StatusNotImplemented, err}
}

var (
	errRefNotFound  = errors.New("Ref not found")
	errRepoNotFound = errors.New("Repository not found")
	errBlobNotFound = errors.New("Blob not found")
	errInvalidHash  = errors.New("Invalid hash")
)

func getRepo(env *Env, r *http.Request) (*git.Repository, error) {
	vars := mux.Vars(r)
	repo := vars["repo"]
	cfg, ok := env.Repositories[repo]

	if !ok {
		return nil, errRepoNotFound
	}

	if cfg.InMemory != nil {
		return cfg.InMemory, nil
	}

	g, err := git.PlainOpen(cfg.Path)
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
