package main

import (
	"github.com/alexa-infra/git47/app/frontend/handlers"
	"github.com/alexa-infra/git47/app/frontend/middleware"
	"github.com/alexa-infra/git47/app/frontend/server"
	"net/http"
)

func BuildPipeline(env *server.Env) {
	r := env.Router
	r.Use(middleware.Logging)

	handlers.RegisterHandlers(env, r)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(env.StaticPath))))

}
