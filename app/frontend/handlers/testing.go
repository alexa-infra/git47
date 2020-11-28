package handlers

import (
	"github.com/alexa-infra/git47/app/frontend/server"
	"testing"
)

var tEnv *server.Env

func makeTestEnv(t *testing.T) *server.Env {
	if tEnv != nil {
		return tEnv
	}

	repo := server.MakeTestRepository(t)

	env := server.NewEnv(server.EnvConfig{
		StaticPath: "../static",
		TemplatePath: "../../../templates",
	})
	env.AddRepoInMemory("memory", repo)

	r := env.Router
	s := r.PathPrefix("/r/").Subrouter()
	RegisterHandlers(env, s)
	tEnv = env
	return env
}
