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
	Files      []fileData
	Dirs       []fileData
	Path       string
	LastCommit *commitData
	*Context
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

func gitTree(ctx *Context) error {
	g := ctx.repo
	commit := ctx.Ref.Commit
	r := ctx.request

	tree, err := commit.Tree()
	if err != nil {
		return err
	}

	baseURL := ctx.GetTreeURL()
	path := r.URL.Path[len(baseURL):]
	path = strings.Trim(path, "/")

	data := treeViewData{Context: ctx}

	newFileData := func(commit *object.Commit, name string) fileData {
		return fileData{
			Name:   name,
			URL:    ctx.GetBlobURL(path, name),
			Commit: newCommitData(ctx, commit),
			Kind:   "File",
		}
	}

	newFolderData := func(commit *object.Commit, name string) fileData {
		return fileData{
			Name:   name,
			URL:    ctx.GetTreeURL(path, name),
			Commit: newCommitData(ctx, commit),
			Kind:   "Folder",
		}
	}

	if path != "" {
		tree, err = tree.Tree(path)
		if err != nil {
			return StatusError{http.StatusNotFound, err}
		}

		data.Dirs = append(data.Dirs, newFolderData(nil, ".."))

		lastCommit, err := getLastCommit(g, commit, path)
		if err != nil {
			return err
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
		return err
	}

	return ctx.RenderTemplate("git-tree.html", data)
}
