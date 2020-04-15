package handlers

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gorilla/mux"
	"net/http"
)

type StaticConfig struct {
	Path string
}

type RepoConfig struct {
	Name        string
	Path        string
	Description string
	InMemory    *git.Repository
}

func (rc *RepoConfig) open() (*git.Repository, error) {
	if rc.InMemory != nil {
		return rc.InMemory, nil
	}

	return git.PlainOpen(rc.Path)
}

type RepoMap map[string]*RepoConfig

type Env struct {
	Template     *TemplateConfig
	Router       *mux.Router
	Repositories RepoMap
	Static       *StaticConfig
}

func (env *Env) getRepoConfig(r *http.Request) (*RepoConfig, error) {
	vars := mux.Vars(r)
	repo := vars["repo"]
	cfg, ok := env.Repositories[repo]

	if !ok {
		return nil, errRepoNotFound
	}
	return cfg, nil
}

func (env *Env) getCommitURL(rc *RepoConfig, c *object.Commit) string {
	router := env.Router
	commitRoute := router.Get("commit")
	repoName := rc.Name

	url, err := commitRoute.URLPath("repo", repoName, "hash", c.Hash.String())
	if err != nil {
		return ""
	}
	return url.Path
}
