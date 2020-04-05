package main

import (
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

var gitDiffTemplate = template.Must(template.ParseFiles("views/base.html", "views/git-diff.html"))

func gitDiff(w http.ResponseWriter, r *http.Request) {
	g := getRepoVar(r)

	vars := mux.Vars(r)
	hashStr := vars["hash"]

	hash := plumbing.NewHash(hashStr)
	if hash.IsZero() {
		http.Error(w, "Invalid hash", http.StatusInternalServerError)
		return
	}

	commit, err := g.CommitObject(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stats, err := commit.Stats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = gitDiffTemplate.ExecuteTemplate(w, "layout", stats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
