package main

import (
	"github.com/alexa-infra/git47/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	handlers.MakeRoutes(r)

	// Start server
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":1323", nil))
}
