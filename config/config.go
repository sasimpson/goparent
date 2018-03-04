package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

//Env - container for all environment configuraitons
type Env struct {
	Service Service
	DB      DBEnv
	Auth    Authentication
}

//Service - structure for service configurations
type Service struct {
	Host string
	Port int
}

//Authentication - structure for authentication configurations
type Authentication struct {
	SigningKey []byte
}

//InitConfig - setup and read configuration for the service
func InitConfig() *Env {
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

	return &Env{
		Service: Service{
			Host: viper.GetString("service.host"),
			Port: viper.GetInt("service.port")},
		DB: DBEnv{
			Host:     viper.GetString("rethinkdb.host"),
			Port:     viper.GetInt("rethinkdb.port"),
			Database: viper.GetString("rethinkdb.name"),
			Username: viper.GetString("rethinkdb.username"),
			Password: viper.GetString("rethinkdb.password")},
		Auth: Authentication{
			SigningKey: []byte(viper.GetString("auth.signingkey"))},
	}
}
