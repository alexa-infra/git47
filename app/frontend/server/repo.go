package server

import (
	"github.com/go-git/go-git/v5"
)

type RepoConfig struct {
	Name        string
	Path        string
	Description string
	inMemory    *git.Repository
}

func (rc *RepoConfig) Open() (*git.Repository, error) {
	if rc.inMemory != nil {
		return rc.inMemory, nil
	}

	return git.PlainOpen(rc.Path)
}

type repoMap map[string]*RepoConfig
