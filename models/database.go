package models

import "gopkg.in/gorethink/gorethink.v3"

func GetConnection() (*gorethink.Session, error) {
	session, err := gorethink.Connect(gorethink.ConnectOpts{
		Address:  "localhost:28015",
		Database: "goparent",
	})
	return session, err
}
