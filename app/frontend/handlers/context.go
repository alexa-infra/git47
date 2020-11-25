package handlers

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"net/http"
)

type Context struct {
	response http.ResponseWriter
	request  *http.Request
	*Env
	Config *RepoConfig
	Ref    *NamedReference
	repo   *git.Repository
}

func (c *Context) GetTreeURL(path ...string) string {
	route := c.Env.Router.Get("tree")
	repoName := c.Config.Name
	refName := c.Ref.Name

	url, err := route.URLPath("repo", repoName, "ref", refName)
	if err != nil {
		return ""
	}
	return joinURL(url.Path, path...)
}

func (c *Context) GetBlobURL(path ...string) string {
	route := c.Env.Router.Get("blob")
	repoName := c.Config.Name
	refName := c.Ref.Name

	url, err := route.URLPath("repo", repoName, "ref", refName)
	if err != nil {
		return ""
	}
	return joinURL(url.Path, path...)
}

func (c *Context) GetCommitURL(commit *object.Commit) string {
	route := c.Env.Router.Get("commit")
	repoName := c.Config.Name

	url, err := route.URLPath("repo", repoName, "hash", commit.Hash.String())
	if err != nil {
		return ""
	}
	return url.Path
}

func (c *Context) GetLogURL() string {
	route := c.Env.Router.Get("commits")
	url, err := route.URLPath("repo", c.Config.Name, "ref", c.Ref.Name)
	if err != nil {
		return ""
	}
	return url.Path
}

func (c *Context) GetSummaryURL() string {
	route := c.Env.Router.Get("summary")
	url, err := route.URLPath("repo", c.Config.Name)
	if err != nil {
		return ""
	}
	return url.Path
}

func (c *Context) RenderTemplate(name string, data interface{}) error {
	config := c.Env.Template
	template, err := config.GetTemplate(name)
	if err != nil {
		return err
	}
	return template.ExecuteTemplate(c.response, "layout", data)
}
