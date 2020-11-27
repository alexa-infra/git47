package server

import (
	"github.com/go-git/go-git/v5"
	"github.com/gorilla/mux"
	"net/http"
	"log"
)

type EnvConfig struct {
	StaticPath string
	TemplatePath string
}

type Env struct {
	Router *mux.Router
	Repositories repoMap
	EnvConfig
}

func NewEnv(config EnvConfig) *Env {
	return &Env{
		Router: mux.NewRouter(),
		Repositories: make(repoMap),
		EnvConfig: config,
	}
}

func (env *Env) AddRepo(name string, path string) {
	env.Repositories[name] = &RepoConfig{
		Name: name,
		Path: path,
	}
}

func (env *Env) AddRepoInMemory(name string, repo *git.Repository) {
	env.Repositories[name] = &RepoConfig{
		Name: name,
		inMemory: repo,
	}
}

func (env *Env) Start() {
	http.Handle("/", env.Router)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
