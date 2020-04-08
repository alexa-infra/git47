package handlers

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

var gitLogTemplate = parseTemplate("templates/base.html", "templates/git-commits.html")

type commitData struct {
	URL     string
	Hash    string
	Message string
}

type commitsViewData struct {
	Commits []commitData
}

func gitLog(w http.ResponseWriter, r *http.Request) {
	g := getRepoVar(r)

	ref, err := getRef(r, g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	query := r.URL.Query()
	next := query.Get("next")

	if next != "" {
		ref = plumbing.NewHash(next)
		if ref.IsZero() {
			http.Error(w, errRefNotFound.Error(), http.StatusNotFound)
			return
		}
	}

	// ... retrieves the commit history
	cIter, err := g.Log(&git.LogOptions{From: ref})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	router := getRouter(r)
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

	err = gitLogTemplate.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
