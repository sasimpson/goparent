package rethinkdb

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/sasimpson/goparent"
	"github.com/spf13/viper"
	"gopkg.in/gorethink/gorethink.v3"
)

//DBEnv - stores the connection parameters for the rethinkdb instance
type DBEnv struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
	Session  gorethink.QueryExecutor
}

//GetConnection - get a connection to the db
func (dbenv *DBEnv) GetConnection() error {
	if dbenv.Session != nil && dbenv.Session.IsConnected() {
		return nil
	}
	connectOpts := gorethink.ConnectOpts{
		Address:  fmt.Sprintf("%s:%d", dbenv.Host, dbenv.Port),
		Database: dbenv.Database,
		Username: dbenv.Username,
		Password: dbenv.Password,
	}

	session, err := gorethink.Connect(connectOpts)
	if err != nil {
		return err
	}
	dbenv.Session = session
	return nil
}

//GetContext returns the request context to satisfy the interface needs
func (dbenv *DBEnv) GetContext(r *http.Request) context.Context {
	return r.Context()
}

//CreateTables - this will build out the database.
func CreateTables(dbenv *DBEnv) {
	err := dbenv.GetConnection()
	if err != nil {
		log.Fatal(err.Error())
	}

	gorethink.DBCreate("goparent").Run(dbenv.Session)
	gorethink.DB("goparent").TableCreate("feeding").Run(dbenv.Session)
	gorethink.DB("goparent").TableCreate("waste").Run(dbenv.Session)
	gorethink.DB("goparent").TableCreate("sleep").Run(dbenv.Session)
	gorethink.DB("goparent").TableCreate("users").Run(dbenv.Session)
	gorethink.DB("goparent").TableCreate("children").Run(dbenv.Session)
	gorethink.DB("goparent").TableCreate("invites").Run(dbenv.Session)
	gorethink.DB("goparent").TableCreate("family").Run(dbenv.Session)
}

//InitConfig - setup and read configuration for the service
func InitRethinkDBConfig() (*goparent.Env, *DBEnv) {
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
			DB: &DBEnv{
				Host:     viper.GetString("rethinkdb.host"),
				Port:     viper.GetInt("rethinkdb.port"),
				Database: viper.GetString("rethinkdb.name"),
				Username: viper.GetString("rethinkdb.username"),
				Password: viper.GetString("rethinkdb.password")},
			Auth: goparent.Authentication{
				SigningKey: []byte(viper.GetString("auth.signingkey"))},
		}, &DBEnv{
			Host:     viper.GetString("rethinkdb.host"),
			Port:     viper.GetInt("rethinkdb.port"),
			Database: viper.GetString("rethinkdb.name"),
			Username: viper.GetString("rethinkdb.username"),
			Password: viper.GetString("rethinkdb.password")}
}
