package main

import (
	"github.com/sasimpson/goparent/api"
	"github.com/sasimpson/goparent/config"
)

var env *config.Env

func main() {
	env = config.InitConfig()
	api.RunService(env)
}
