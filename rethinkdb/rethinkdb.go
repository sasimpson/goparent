package rethinkdb

import (
	"fmt"
	"log"

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
