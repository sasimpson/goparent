package main

import (
	"github.com/sasimpson/goparent/api"
	"github.com/sasimpson/goparent/config"
)

func main() {
	env := config.InitConfig()
	config.CreateTables(env)
	api.RunService(env)
}
