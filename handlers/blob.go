package handlers

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func gitBlob(env *Env, w http.ResponseWriter, r *http.Request) error {
	rc, err := env.getRepoConfig(r)
	if err != nil {
		return StatusError{http.StatusNotFound, err}
	}

	g, err := rc.open()
	if err != nil {
		return StatusError{http.StatusNotFound, err}
	}

	ref, err := getNamedRef(g, r)
	if err != nil {
		return StatusError{http.StatusNotFound, err}
	}

	commit := ref.Commit
	tree, err := commit.Tree()
	if err != nil {
		return err
	}

	baseURL := env.getBlobURL(rc, ref)
	path := r.URL.Path[len(baseURL):]
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
