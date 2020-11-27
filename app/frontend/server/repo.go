package server

import (
	"github.com/go-git/go-git/v5"
)

type repoConfig struct {
	Name        string
	Path        string
	Description string
	inMemory    *git.Repository
}

func (rc *repoConfig) open() (*git.Repository, error) {
	if rc.inMemory != nil {
		return rc.inMemory, nil
	}

	return git.PlainOpen(rc.Path)
}

type repoMap map[string]*repoConfig
