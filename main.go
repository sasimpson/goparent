package main

import (
	"github.com/sasimpson/goparent/api"
	"github.com/sasimpson/goparent/config"
)

func main() {
	env := config.Env{
		DB:   config.DBEnv{Host: "localhost", Port: 28015, Database: "goparent"},
		Auth: config.Authentication{SigningKey: []byte("supersecretsquirrl")},
	}
	api.RunService(&env)
}
