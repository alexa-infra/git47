package main

import (
	"github.com/alexa-infra/git47/internal/configs"
	"github.com/alexa-infra/git47/internal/web"
	"log"
)

func main() {
	configs, err := configs.NewConfigs()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	httpCfg, err := configs.HTTP()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	repositories, err := configs.Repositories()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	app, err := web.NewApp(httpCfg, repositories)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	app.Start()
}
