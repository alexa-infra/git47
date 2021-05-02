package web

import (
	"github.com/go-git/go-git/v5/plumbing/object"
	"html/template"
)

// GetBlobURL builds URL of blob page
func GetBlobURL(rc *RequestContext, path ...string) (string, error) {
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

// GetSummaryURL builds URL of summary page
func GetSummaryURL(rc *RequestContext) (string, error) {
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

// GetCommitURL builds URL of commit diff
func GetCommitURL(rc *RequestContext, commit *object.Commit) (string, error) {
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

// GetLogURL builds URL of commits page
func GetLogURL(rc *RequestContext) (string, error) {
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

// GetTreeURL builds URL of git tree by path
func GetTreeURL(rc *RequestContext, path ...string) (string, error) {
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

// GetBranchesURL builds URL of branches list page
func GetBranchesURL(rc *RequestContext) (string, error) {
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

// GetTagsURL builds URL of tags list page
func GetTagsURL(rc *RequestContext) (string, error) {
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

// GetContributorsURL builds URL of branches list page
func GetContributorsURL(rc *RequestContext) (string, error) {
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

// TemplateHelpers returns a list of helper functions used in templates
func TemplateHelpers() template.FuncMap {
	return template.FuncMap{
		"GetSummaryURL":      GetSummaryURL,
		"GetBlobURL":         GetBlobURL,
		"GetTreeURL":         GetTreeURL,
		"GetLogURL":          GetLogURL,
		"GetBranchesURL":     GetBranchesURL,
		"GetTagsURL":         GetTagsURL,
		"GetContributorsURL": GetContributorsURL,
		"GetCommitURL":       GetCommitURL,
	}
}
