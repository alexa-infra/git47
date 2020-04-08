package handlers

import (
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func gitBlob(w http.ResponseWriter, r *http.Request) {
	g := getRepoVar(r)

	vars := mux.Vars(r)
	hashStr := vars["hash"]

	hash := plumbing.NewHash(hashStr)
	if hash.IsZero() {
		http.Error(w, errBlobNotFound.Error(), http.StatusNotFound)
		return
	}

	blob, err := g.BlobObject(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reader, err := blob.Reader()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
