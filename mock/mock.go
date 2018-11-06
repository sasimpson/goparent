package mock

import (
	"context"
	"net/http"
)

//DBEnv mock
type DBEnv struct{}

//GetConnection interface impl mock
func (db *DBEnv) GetConnection() error {
	panic("not implemented")
}

//GetContext interface impl mock
func (db *DBEnv) GetContext(r *http.Request) context.Context {
	return r.Context()
}
