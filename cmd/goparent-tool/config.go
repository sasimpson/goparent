package main

import (
	"fmt"
	"log"

	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/rethinkdb"
	"github.com/spf13/viper"
)

//InitRethinkDBConfig - setup and read configuration for the service
func InitRethinkDBConfig() (*goparent.Env, *rethinkdb.DBEnv) {
	//set defaults
	viper.SetDefault("service.host", "localhost")
	viper.SetDefault("service.port", "8000")
	viper.SetDefault("rethinkdb.host", "localhost")
	viper.SetDefault("rethinkdb.port", 28015)
	viper.SetDefault("rethinkdb.name", "goparent")
	viper.SetDefault("auth.signingkey", "supersecretsquirrl")

	//parse configs if they exist
	viper.SetConfigName("goparent")
	viper.AddConfigPath("/etc/config/")
	viper.AddConfigPath("$HOME/.goparent")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("no config read in: %s", err))
	}
	log.Println("config used:", viper.ConfigFileUsed())

	return &goparent.Env{
			Service: goparent.Service{
				Host: viper.GetString("service.host"),
				Port: viper.GetInt("service.port")},
			Auth: goparent.Authentication{
				SigningKey: []byte(viper.GetString("auth.signingkey"))},
		}, &rethinkdb.DBEnv{
			Host:     viper.GetString("rethinkdb.host"),
			Port:     viper.GetInt("rethinkdb.port"),
			Database: viper.GetString("rethinkdb.name"),
			Username: viper.GetString("rethinkdb.username"),
			Password: viper.GetString("rethinkdb.password")}
}
