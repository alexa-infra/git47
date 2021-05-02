package configs

import (
	"github.com/alexa-infra/git47/internal/core"
	"github.com/alexa-infra/git47/internal/web"
	"os"
)

func getenv(name string, defaultValue string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}
	return value
}

func hasenv(name string) bool {
	_, ok := os.LookupEnv(name)
	return ok
}

type Configs struct{}

func (cfg *Configs) Repositories() (core.RepoMap, error) {
	return core.RepoMap{
		"inforia": {
			Name: "inforia",
			Path: "/home/alexey/projects/inforia/sql/.git",
		},
		"git47": {
			Name: "git47",
			Path: "/home/alexey/projects/go-playground/git47/.git",
		},
	}, nil
}

func (cfg *Configs) HTTP() (*web.Config, error) {
	return &web.Config{
		Host:    getenv("HOST", "0.0.0.0"),
		Port:    getenv("PORT", "8080"),
		Logging: !hasenv("DISABLE_LOGS"),
	}, nil
}

func NewConfigs() (*Configs, error) {
	return &Configs{}, nil
}
