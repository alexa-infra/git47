package handlers

import (
	"github.com/alexa-infra/git47/app/frontend/server"
	"io/ioutil"
	"net/http"
	"strings"
)

// GitBlob returns handler, which renders single blob to http output (not formatted)
func GitBlob(env *server.Env) http.HandlerFunc {
	return env.WrapHandler(func (w http.ResponseWriter, r *http.Request) {
		ctx, _ := server.GetRequestContext(r)
		ref := ctx.Ref

		commit := ref.Commit
		tree, err := commit.Tree()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		baseURL, err := GetBlobURL(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		path := r.URL.Path[len(baseURL):]
		path = strings.Trim(path, "/")

		if path == "" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		file, err := tree.File(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		reader, err := file.Reader()
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
	})
}

// GetBlobURL builds URL of blob page
func GetBlobURL(rc *server.RequestContext, path ...string) (string, error) {
	router := rc.Env.Router
	route := router.Get("blob")
	url, err := route.URLPath("repo", rc.Config.Name, "ref", rc.Ref.Name)
	if err != nil {
		return "", err
	}
	return joinURL(url.Path, path...), nil
}
