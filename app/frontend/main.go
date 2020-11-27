package main

import (
	"github.com/alexa-infra/git47/app/frontend/server"
)

func main() {
	env := server.NewEnv(server.EnvConfig{
		StaticPath:   "./static",
		TemplatePath: "./templates",
	})
	env.AddRepo("inforia", "/home/alexey/projects/inforia/main/.git")
	env.AddRepo("git47", "/home/alexey/projects/go-playground/git47/.git")
	BuildPipeline(env)
	env.Start()
}
