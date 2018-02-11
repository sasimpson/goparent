package config

import (
	"fmt"
	"log"

	gorethink "gopkg.in/gorethink/gorethink.v3"
)

//DBEnv - Environment for DB settings
type DBEnv struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
	Session  gorethink.QueryExecutor
}

//GetConnection - get a connection to the db
func (dbenv *DBEnv) GetConnection() (gorethink.QueryExecutor, error) {
	if dbenv.Session != nil && dbenv.Session.IsConnected() {
		return dbenv.Session, nil
	}
	connectOpts := gorethink.ConnectOpts{
		Address:  fmt.Sprintf("%s:%d", dbenv.Host, dbenv.Port),
		Database: dbenv.Database,
		Username: dbenv.Username,
		Password: dbenv.Password,
	}

	session, err := gorethink.Connect(connectOpts)
	return session, err
}

//CreateTables - this will build out the database.
func CreateTables(env *Env) {
	session, err := env.DB.GetConnection()
	if err != nil {
		log.Fatal(err.Error())
	}
	gorethink.DBCreate("goparent").Run(session)
	gorethink.DB("goparent").TableCreate("feeding").Run(session)
	gorethink.DB("goparent").TableCreate("waste").Run(session)
	gorethink.DB("goparent").TableCreate("sleep").Run(session)
	gorethink.DB("goparent").TableCreate("users").Run(session)
	gorethink.DB("goparent").TableCreate("children").Run(session)
}
