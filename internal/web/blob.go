package web

import (
	"github.com/alexa-infra/git47/internal/core"
	"net/http"
	"strings"
)

func blobHandler(w http.ResponseWriter, r *http.Request) {
	ctx, ok := getRequestContext(r)
	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	baseURL, err := getBlobURL(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	path := r.URL.Path[len(baseURL):]
	path = strings.Trim(path, "/")

	if path == "" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := core.GetBlob(ctx.Ref, path)
	if err != nil {
		if err == core.ErrFileNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
