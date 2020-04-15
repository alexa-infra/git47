package main

import (
	. "github.com/alexa-infra/git47/handlers"
	mw "github.com/alexa-infra/git47/middleware"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"flag"
)

func main() {
	var noCache bool
	flag.BoolVar(&noCache, "nocache", false, "disable template cache")
	flag.Parse()

	r := mux.NewRouter()
	r.Use(mw.Logging)

	env := &Env{
		Router: r,
		Template: &TemplateConfig{
			Path: "./templates",
			UseCache: !noCache,
		},
		Repositories: RepoMap{
			"friday": &RepoConfig{
				Name: "friday",
				Path: "/home/alexey/projects/friday/.git",
			},
			"git47": &RepoConfig{
				Name: "git47",
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
	log.Fatal(http.ListenAndServe(":8080", nil))
}
