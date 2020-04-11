package handlers

import (
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func gitBlob(env *Env, w http.ResponseWriter, r *http.Request) error {
	g, err := getRepo(env, r)
	if err != nil {
		return StatusError{http.StatusNotFound, err}
	}

	vars := mux.Vars(r)
	hashStr := vars["hash"]

	hash := plumbing.NewHash(hashStr)
	if hash.IsZero() {
		return StatusError{http.StatusBadRequest, errInvalidHash}
	}

	blob, err := g.BlobObject(hash)
	if err != nil {
		return StatusError{http.StatusNotFound, err}
	}

	reader, err := blob.Reader()
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
