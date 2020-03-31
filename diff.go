package main

import (
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/gorilla/mux"
	"net/http"
)

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

	files, err := commit.Stats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = renderTemplate(w, "git-diff.html", files)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
