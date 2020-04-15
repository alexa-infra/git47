package handlers

import (
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
	Files []*fileData
	Dirs  []*fileData
	*RepoConfig
	*NamedReference
}

func getFileLastCommit(g *git.Repository, ref *object.Commit, paths ...string) (*object.Commit, error) {
	path := joinPath(paths...)
	filterFile := func(p string) bool {
		return p == path
	}
	cIter, err := g.Log(&git.LogOptions{From: ref.Hash, PathFilter: filterFile})
	if err != nil {
		return nil, err
	}
	return cIter.Next()
}

func getFolderLastCommit(g *git.Repository, ref *object.Commit, paths ...string) (*object.Commit, error) {
	path := joinPath(paths...)
	filterFolder := func(p string) bool {
		return strings.HasPrefix(p, path)
	}
	cIter, err := g.Log(&git.LogOptions{From: ref.Hash, PathFilter: filterFolder})
	if err != nil {
		return nil, err
	}
	return cIter.Next()
}

func gitTree(env *Env, w http.ResponseWriter, r *http.Request) error {
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

	commit, err := g.CommitObject(ref.Hash)
	if err != nil {
		return err
	}

	tree, err := commit.Tree()
	if err != nil {
		return err
	}

	baseURL := env.getTreeURL(rc, ref)
	path := r.URL.Path[len(baseURL):]
	path = strings.Trim(path, "/")

	if path != "" {
		tree, err = tree.Tree(path)
		if err != nil {
			return StatusError{http.StatusNotFound, err}
		}
	}

	data := treeViewData{RepoConfig: rc, NamedReference: ref}
	if path != "" {
		data.Dirs = append(data.Dirs, &fileData{
			Name: "..",
			Kind: "Folder",
			URL:  env.getTreeURL(rc, ref, parentPath(path)),
		})
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
			lastCommit, err := getFolderLastCommit(g, commit, path, folderName)
			if err != nil {
				return err
			}
			cd := newCommitData(lastCommit)
			cd.URL = env.getCommitURL(rc, lastCommit)
			uniqDirs[folderName] = true
			data.Dirs = append(data.Dirs, &fileData{
				Name: folderName,
				URL:  env.getTreeURL(rc, ref, path, folderName),
				Commit: cd,
				Kind: "Folder",
			})
		} else {
			lastCommit, err := getFileLastCommit(g, commit, path, f.Name)
			if err != nil {
				return err
			}
			cd := newCommitData(lastCommit)
			cd.URL = env.getCommitURL(rc, lastCommit)
			data.Files = append(data.Files, &fileData{
				Name: f.Name,
				URL:  env.getBlobURL(rc, ref, path, f.Name),
				Commit: cd,
				Kind: "File",
			})
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
