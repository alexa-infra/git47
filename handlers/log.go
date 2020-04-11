package handlers

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type commitData struct {
	URL     string
	Hash    string
	Message string
}

type commitsViewData struct {
	Commits []commitData
}

func gitLog(env *Env, w http.ResponseWriter, r *http.Request) error {
	g, err := getRepo(env, r)
	if err != nil {
		return StatusError{http.StatusNotFound, err}
	}

	ref, err := getRef(r, g)
	if err != nil {
		return StatusError{http.StatusNotFound, err}
	}

	query := r.URL.Query()
	next := query.Get("next")

	if next != "" {
		ref = plumbing.NewHash(next)
		if ref.IsZero() {
			return StatusError{http.StatusBadRequest, errInvalidHash}
		}
	}

	// ... retrieves the commit history
	cIter, err := g.Log(&git.LogOptions{From: ref})
	if err != nil {
		return err
	}

	router := env.Router
	commitRoute := router.Get("commit")
	vars := mux.Vars(r)
	repoName := vars["repo"]

	var data commitsViewData

	// ... just iterates over the commits
	for i := 0; i < 20; i++ {
		c, err := cIter.Next()
		if err != nil {
			break
		}
		commitURL, err := commitRoute.URLPath("repo", repoName, "hash", c.Hash.String())
		if err != nil {
			break
		}
		data.Commits = append(data.Commits, commitData{
			URL:     commitURL.Path,
			Message: strings.Trim(c.Message, "\n"),
			Hash:    c.Hash.String(),
		})
	}

	template, err := env.Template.GetTemplate("git-commits.html")
	if err != nil {
		return err
	}

	return template.ExecuteTemplate(w, "layout", data)
}
