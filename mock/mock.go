package mock

import (
	"context"
	"net/http"
)

type DBEnv struct{}

func (db *DBEnv) GetConnection() error {
	panic("not implemented")
}

func (db *DBEnv) GetContext(*http.Request) context.Context {
	return context.Background()
}
