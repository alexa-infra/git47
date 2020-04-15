package handlers

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
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

type NamedReference struct {
	Name   string
	Kind   string
	Commit *object.Commit
}

func (ref *NamedReference) Hash() plumbing.Hash {
	return ref.Commit.Hash
}

func getNamedRef(g *git.Repository, r *http.Request) (*NamedReference, error) {
	vars := mux.Vars(r)
	ref := vars["ref"]

	if ref == "" {
		return nil, errRefNotFound
	}

	branch, err := g.Reference(plumbing.NewBranchReferenceName(ref), false)
	if err == nil {
		hash := branch.Hash()

		commit, err := g.CommitObject(hash)
		if err != nil {
			return nil, err
		}

		return &NamedReference{
			Name:   ref,
			Kind:   "branch",
			Commit: commit,
		}, nil
	}
	tag, err := g.Reference(plumbing.NewTagReferenceName(ref), false)
	if err == nil {
		hash := tag.Hash()

		commit, err := g.CommitObject(hash)
		if err != nil {
			return nil, err
		}

		return &NamedReference{
			Name:   ref,
			Kind:   "tag",
			Commit: commit,
		}, nil
	}

	hash := plumbing.NewHash(ref)
	if !hash.IsZero() {
		commit, err := g.CommitObject(hash)
		if err == nil {
			return &NamedReference{
				Name:   ref,
				Kind:   "commit",
				Commit: commit,
			}, nil
		}
	}

	return nil, errRefNotFound
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

func joinPath(paths ...string) string {
	return path.Join(paths...)
}
