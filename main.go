package main

import (
	. "github.com/alexa-infra/git47/handlers"
	mw "github.com/alexa-infra/git47/middleware"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.Use(mw.Logging)

	env := &Env{
		Router: r,
		Template: &TemplateConfig{
			Path: "./templates",
		},
		Repositories: RepoMap{
			"friday": &RepoConfig{
				Path: "/home/alexey/projects/friday/.git",
			},
			"git47": &RepoConfig{
				Path: "/home/alexey/projects/go-playground/git47/.git",
			},
		},
		Static: &StaticConfig{
			Path: "./static",
		},
	}
	env.Setup()

	// Start server
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":1323", nil))
}
