package goparent

import (
	"github.com/sasimpson/goparent/api"
	"github.com/sasimpson/rateup/config"
)

//This file is specifically for running in GCP AppEngine.
func main() {
	env = config.InitConfig()
	api.RunAppEngineService(env)
}
