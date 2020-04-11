package handlers

import (
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type treeViewData struct {
	Files      map[string]string
	Dirs       map[string]string
	ParentPath string
}

func gitTree(env *Env, w http.ResponseWriter, r *http.Request) error {
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

	if path != "" {
		tree, err = tree.Tree(path)
		if err != nil {
			return StatusError{http.StatusNotFound, err}
		}
	}

	var data treeViewData
	data.Dirs = make(map[string]string)
	data.Files = make(map[string]string)
	if path != "" {
		data.ParentPath = joinURL(baseURL.Path, parentPath(path))
	}

	blobRoute := router.Get("blob")
	tree.Files().ForEach(func(f *object.File) error {
		if strings.Index(f.Name, "/") > 0 {
			components := strings.Split(f.Name, "/")
			folderName := components[0]
			data.Dirs[folderName] = joinURL(baseURL.Path, path, folderName)
		} else {
			blobURL, err := blobRoute.URLPath("repo", repoName, "hash", f.Hash.String())
			if err != nil {
				return err
			}
			data.Files[f.Name] = joinURL(blobURL.Path, path, f.Name)
		}
		return nil
	})

	template, err := env.Template.GetTemplate("git-list.html")
	if err != nil {
		return err
	}

	return template.ExecuteTemplate(w, "layout", data)
}
