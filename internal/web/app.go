package web

import (
	"fmt"
	"github.com/alexa-infra/git47/internal/core"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type App struct {
	Router *mux.Router
	Config
}

func (app *App) Start() {
	http.Handle("/", app.Router)
	addr := fmt.Sprintf("%s:%s", app.Host, app.Port)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func NewApp(config *Config, repositories core.RepoMap) (*App, error) {
	router, err := NewRouter(config, repositories)
	if err != nil {
		return nil, err
	}
	return &App{
		Config: *config,
		Router: router,
	}, nil
}
