package handlers

import (
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type dirData struct {
	Name string
	URL  string
}

type fileData struct {
	Name string
	URL  string
}

type treeViewData struct {
	Files      []fileData
	Dirs       []dirData
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
	if path != "" {
		parentDir := dirData{
			Name: "..",
			URL: joinURL(baseURL.Path, parentPath(path)),
		}
		data.Dirs = append(data.Dirs, parentDir)
	}
	uniqDirs := make(map[string]bool)

	blobURL, err := router.Get("blob").URLPath("repo", repoName, "ref", refName)
	if err != nil {
		return err
	}

	err = tree.Files().ForEach(func(f *object.File) error {
		if strings.Index(f.Name, "/") > 0 {
			components := strings.Split(f.Name, "/")
			folderName := components[0]
			_, ok := uniqDirs[folderName]
			if ok {
				return nil
			}
			uniqDirs[folderName] = true
			dir := dirData{
				Name: folderName,
				URL: joinURL(baseURL.Path, path, folderName),
			}
			data.Dirs = append(data.Dirs, dir)
		} else {
			file := fileData{
				Name: f.Name,
				URL: joinURL(blobURL.Path, path, f.Name),
			}
			data.Files = append(data.Files, file)
		}
		return nil
	})
	if err != nil {
		return err
	}

	template, err := env.Template.GetTemplate("git-list.html")
	if err != nil {
		return err
	}

	return template.ExecuteTemplate(w, "layout", data)
}
