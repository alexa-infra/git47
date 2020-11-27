package main

import (
	"github.com/alexa-infra/git47/app/frontend/handlers"
	"github.com/alexa-infra/git47/app/frontend/middleware"
	"github.com/alexa-infra/git47/app/frontend/server"
	"net/http"
)

func buildPipeline(env *server.Env) {
	r := env.Router
	r.Use(middleware.Logging)

	s := r.PathPrefix("/r/").Subrouter()
	handlers.RegisterHandlers(env, s)

	staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir(env.StaticPath)))
	r.PathPrefix("/static/").Handler(staticHandler)

}
