package handlers

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strings"
)

func gitBlob(env *Env, w http.ResponseWriter, r *http.Request) error {
	g, err := getRepo(env, r)
	if err != nil {
		return StatusError{http.StatusNotFound, err}
	}

	router := env.Router
	vars := mux.Vars(r)
	repoName := vars["repo"]
	refName := vars["ref"]

	ref, err := getRef(r, g)
	if err != nil {
		return StatusError{http.StatusNotFound, err}
	}

	commit, err := g.CommitObject(ref)
	if err != nil {
		return err
	}

	tree, err := commit.Tree()
	if err != nil {
		return err
	}

	baseURL, err := router.Get("tree").URLPath("repo", repoName, "ref", refName)
	if err != nil {
		return err
	}

	path := r.URL.Path[len(baseURL.Path):]
	path = strings.Trim(path, "/")

	if path == "" {
		return StatusError{http.StatusNotFound, err}
	}

	file, err := tree.File(path)
	if err != nil {
		return StatusError{http.StatusNotFound, err}
	}

	reader, err := file.Reader()
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}
