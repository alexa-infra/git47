package server

import (
	"github.com/go-git/go-git/v5"
	"github.com/gorilla/mux"
	"net/http"
	"log"
)

// EnvConfig is high-level config of server Env
type EnvConfig struct {
	StaticPath string
	TemplatePath string
}

// Env is server environment, a glue for all services
type Env struct {
	Router *mux.Router
	repositories repoMap
	EnvConfig
}

// NewEnv creates new server environment
func NewEnv(config EnvConfig) *Env {
	return &Env{
		Router: mux.NewRouter(),
		repositories: make(repoMap),
		EnvConfig: config,
	}
}

// AddRepo adds new git repository to the environment
func (env *Env) AddRepo(name string, path string) {
	env.repositories[name] = &repoConfig{
		Name: name,
		Path: path,
	}
}

// AddRepoInMemory adds new in-memory repository to the environment
func (env *Env) AddRepoInMemory(name string, repo *git.Repository) {
	env.repositories[name] = &repoConfig{
		Name: name,
		inMemory: repo,
	}
}

// Start runs listen and serve loop
func (env *Env) Start() {
	http.Handle("/", env.Router)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
