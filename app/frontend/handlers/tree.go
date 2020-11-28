package handlers

import (
	"github.com/alexa-infra/git47/app/frontend/server"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"net/http"
	"strings"
)

type fileData struct {
	Name   string
	URL    string
	Kind   string
	Commit *commitData
}

type treeViewData struct {
	Files      []fileData
	Dirs       []fileData
	Path       string
	LastCommit *commitData
	*server.RequestContext
}

func getLastCommit(g *git.Repository, ref *object.Commit, paths ...string) (*object.Commit, error) {
	path := joinPath(paths...)
	cIter, err := g.Log(&git.LogOptions{From: ref.Hash})
	if err != nil {
		return nil, err
	}
	defer cIter.Close()

	var treeEntry *object.TreeEntry = nil
	var lastCommit *object.Commit = nil
	for {
		commit, err := cIter.Next()
		if err != nil {
			if lastCommit != nil {
				return lastCommit, nil
			}
			return nil, err
		}
		tree, err := commit.Tree()
		if err != nil {
			return nil, err
		}
		entry, err := tree.FindEntry(path)
		if err != nil {
			if lastCommit != nil {
				return lastCommit, nil
			}
			return nil, err
		}
		if treeEntry == nil {
			treeEntry = entry
		} else if *treeEntry != *entry {
			return lastCommit, nil
		}
		lastCommit = commit
	}
}

// GitTree returns handler which renders file-tree of ref/path
func GitTree(env *server.Env) http.HandlerFunc {
	template := env.GetTemplate("git-tree.html", TemplateHelpers())
	return env.WrapHandler(func (w http.ResponseWriter, r *http.Request) {
		ctx, _ := server.GetRequestContext(r)
		g := ctx.Repo
		ref := ctx.Ref
		commit := ref.Commit

		tree, err := commit.Tree()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		baseURL, err := GetTreeURL(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		path := r.URL.Path[len(baseURL):]
		path = strings.Trim(path, "/")

		data := treeViewData{
			RequestContext: ctx,
			Path:           path,
		}

		newFileData := func(commit *object.Commit, name string) fileData {
			url, _ := GetBlobURL(ctx, path, name)
			return fileData{
				Name:   name,
				URL:    url,
				Commit: newCommitData(ctx, commit),
				Kind:   "File",
			}
		}

		newFolderData := func(commit *object.Commit, name string) fileData {
			url, _ := GetTreeURL(ctx, path, name)
			return fileData{
				Name:   name,
				URL:    url,
				Commit: newCommitData(ctx, commit),
				Kind:   "Folder",
			}
		}

		if path != "" {
			tree, err = tree.Tree(path)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			data.Dirs = append(data.Dirs, newFolderData(nil, ".."))

			lastCommit, err := getLastCommit(g, commit, path)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			data.LastCommit = newCommitData(ctx, lastCommit)
		} else {
			data.LastCommit = newCommitData(ctx, commit)
		}

		uniqDirs := make(map[string]bool)
		err = tree.Files().ForEach(func(f *object.File) error {
			if strings.Index(f.Name, "/") > 0 {
				components := strings.Split(f.Name, "/")
				folderName := components[0]
				_, ok := uniqDirs[folderName]
				if ok {
					return nil
				}
				uniqDirs[folderName] = true

				lastCommit, err := getLastCommit(g, commit, path, folderName)
				if err != nil {
					return err
				}
				data.Dirs = append(data.Dirs, newFolderData(lastCommit, folderName))
			} else {
				lastCommit, err := getLastCommit(g, commit, path, f.Name)
				if err != nil {
					return err
				}
				data.Files = append(data.Files, newFileData(lastCommit, f.Name))
			}
			return nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = env.RenderTemplate(w, template, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

// GetTreeURL builds URL of git tree by path
func GetTreeURL(rc *server.RequestContext, path ...string) (string, error) {
	router := rc.Env.Router
	route := router.Get("tree")
	url, err := route.URLPath("repo", rc.Config.Name, "ref", rc.Ref.Name)
	if err != nil {
		return "", err
	}
	return joinURL(url.Path, path...), nil
}
