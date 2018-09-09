package datastore

import (
	"context"
	"net/http"

	"google.golang.org/appengine"
)

//DBEnv -
type DBEnv struct {
}

//GetConnection -
func (db *DBEnv) GetConnection() error {
	panic("not implemented")
}

func (db *DBEnv) GetContext(r *http.Request) context.Context {
	return appengine.NewContext(r)
}
