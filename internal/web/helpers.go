package web

import (
	"github.com/go-git/go-git/v5/plumbing/object"
	"html/template"
)

func getBlobURL(rc *requestContext, path ...string) (string, error) {
	router := rc.Router
	route := router.Get("blob")
	if route == nil {
		return "", nil
	}
	url, err := route.URLPath("repo", rc.RepoConfig.Name, "ref", rc.Ref.Name)
	if err != nil {
		return "", err
	}
	return joinURL(url.Path, path...), nil
}

func getSummaryURL(rc *requestContext) (string, error) {
	router := rc.Router
	route := router.Get("summary")
	if route == nil {
		return "", nil
	}
	url, err := route.URLPath("repo", rc.RepoConfig.Name)
	if err != nil {
		return "", err
	}
	return url.Path, nil
}

func getCommitURL(rc *requestContext, commit *object.Commit) (string, error) {
	router := rc.Router
	route := router.Get("commit")
	if route == nil {
		return "", nil
	}
	url, err := route.URLPath("repo", rc.RepoConfig.Name, "hash", commit.Hash.String())
	if err != nil {
		return "", err
	}
	return url.Path, nil
}

func getLogURL(rc *requestContext) (string, error) {
	router := rc.Router
	route := router.Get("commits")
	if route == nil {
		return "", nil
	}
	url, err := route.URLPath("repo", rc.RepoConfig.Name, "ref", rc.Ref.Name)
	if err != nil {
		return "", err
	}
	return url.Path, nil
}

func getTreeURL(rc *requestContext, path ...string) (string, error) {
	router := rc.Router
	route := router.Get("tree")
	if route == nil {
		return "", nil
	}
	url, err := route.URLPath("repo", rc.RepoConfig.Name, "ref", rc.Ref.Name)
	if err != nil {
		return "", err
	}
	return joinURL(url.Path, path...), nil
}

func getBranchesURL(rc *requestContext) (string, error) {
	router := rc.Router
	route := router.Get("branches")
	if route == nil {
		return "", nil
	}
	url, err := route.URLPath("repo", rc.RepoConfig.Name)
	if err != nil {
		return "", err
	}
	return url.Path, nil
}

func getTagsURL(rc *requestContext) (string, error) {
	router := rc.Router
	route := router.Get("tags")
	if route == nil {
		return "", nil
	}
	url, err := route.URLPath("repo", rc.RepoConfig.Name)
	if err != nil {
		return "", err
	}
	return url.Path, nil
}

func getContributorsURL(rc *requestContext) (string, error) {
	router := rc.Router
	route := router.Get("contributors")
	if route == nil {
		return "", nil
	}
	url, err := route.URLPath("repo", rc.RepoConfig.Name)
	if err != nil {
		return "", err
	}
	return url.Path, nil
}

func templateHelpers() template.FuncMap {
	return template.FuncMap{
		"GetSummaryURL":      getSummaryURL,
		"GetBlobURL":         getBlobURL,
		"GetTreeURL":         getTreeURL,
		"GetLogURL":          getLogURL,
		"GetBranchesURL":     getBranchesURL,
		"GetTagsURL":         getTagsURL,
		"GetContributorsURL": getContributorsURL,
		"GetCommitURL":       getCommitURL,
	}
}
