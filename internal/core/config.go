package core

import (
	"github.com/go-git/go-git/v5"
)

type RepoConfig struct {
	Name        string
	Path        string
	Description string
	InMemory    *git.Repository
}

func (rc RepoConfig) Open() (*git.Repository, error) {
	if rc.InMemory != nil {
		return rc.InMemory, nil
	}

	return git.PlainOpen(rc.Path)
}

type RepoMap map[string]RepoConfig
