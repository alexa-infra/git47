package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"net/http"
	"strings"
)

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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// ... retrieves the commit history
	cIter, err := g.Log(&git.LogOptions{From: ref})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ... just iterates over the commits
	var items []string
	for i := 0; i < 20; i++ {
		c, err := cIter.Next()
		if err != nil {
			break
		}
		items = append(items, fmt.Sprintf("%s %s", c.Hash, strings.Trim(c.Message, "\n")))
	}

	err = renderTemplate(w, "git-diff.html", items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
