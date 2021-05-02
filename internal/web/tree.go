package web

import (
	"github.com/alexa-infra/git47/internal/core"
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
	*RequestContext
}

// GitTree returns handler which renders file-tree of ref/path
func GitTree(w http.ResponseWriter, r *http.Request) {
	ctx, _ := GetRequestContext(r)
	ref := ctx.Ref

	baseURL, err := GetTreeURL(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	path := r.URL.Path[len(baseURL):]
	path = strings.Trim(path, "/")

	tree, err := core.GetTree(ctx.Ref, path)
	if err != nil {
		if err == core.ErrFileNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := treeViewData{
		RequestContext: ctx,
		Path:           path,
	}

	newFileData := func(name string) fileData {
		lastCommit, _ := core.GetLastCommit(ref, path, name)
		url, _ := GetBlobURL(ctx, path, name)
		return fileData{
			Name:   name,
			URL:    url,
			Commit: newCommitData(ctx, lastCommit),
			Kind:   "File",
		}
	}

	newFolderData := func(name string) fileData {
		lastCommit, _ := core.GetLastCommit(ref, path, name)
		url, _ := GetTreeURL(ctx, path, name)
		return fileData{
			Name:   name,
			URL:    url,
			Commit: newCommitData(ctx, lastCommit),
			Kind:   "Folder",
		}
	}

	for _, x := range tree.Files {
		data.Files = append(data.Files, newFileData(x))
	}

	for _, x := range tree.Dirs {
		data.Files = append(data.Files, newFolderData(x))
	}

	if path != "" {
		lastCommit, err := core.GetLastCommit(ref, path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data.LastCommit = newCommitData(ctx, lastCommit)
	} else {
		data.LastCommit = newCommitData(ctx, ref.Commit)
	}

	err = RenderTemplate(w, "git-tree.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
