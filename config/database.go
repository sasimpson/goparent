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
	Session  gorethink.QueryExecutor
}

//GetConnection - get a connection to the db
func (dbenv *DBEnv) GetConnection() (gorethink.QueryExecutor, error) {
	if dbenv.Session != nil && dbenv.Session.IsConnected() {
		return dbenv.Session, nil
	}
	session, err := gorethink.Connect(gorethink.ConnectOpts{
		Address:  fmt.Sprintf("%s:%d", dbenv.Host, dbenv.Port),
		Database: dbenv.Database,
	})
	return session, err
}

func CreateTables(env *Env) {
	session, err := env.DB.GetConnection()
	if err != nil {
		log.Fatal(err.Error())
	}
	gorethink.TableCreate("feeding").Run(session)
	gorethink.TableCreate("waste").Run(session)
	gorethink.TableCreate("sleep").Run(session)
	gorethink.TableCreate("users").Run(session)
	gorethink.TableCreate("children").Run(session)
}
