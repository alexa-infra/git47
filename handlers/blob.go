package handlers

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func gitBlob(ctx *Context) error {
	r := ctx.request
	w := ctx.response
	ref := ctx.Ref

	commit := ref.Commit
	tree, err := commit.Tree()
	if err != nil {
		return err
	}

	baseURL := ctx.GetBlobURL()
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
