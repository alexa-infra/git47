package handlers

import (
	"github.com/go-git/go-git/v5"
	"github.com/gorilla/mux"
)

type StaticConfig struct {
	Path string
}

type RepoConfig struct {
	Path        string
	Description string
	InMemory    *git.Repository
}

type RepoMap map[string]*RepoConfig

type Env struct {
	Template     *TemplateConfig
	Router       *mux.Router
	Repositories RepoMap
	Static       *StaticConfig
}
