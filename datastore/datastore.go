package datastore

import (
	"context"
	"fmt"
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

//GetContext - appengine requires a context from the request, but we don't
//want appengine code burried in the API.  This has been added to the interface,
//so it is abstracted out to the datastore bits.
func (db *DBEnv) GetContext(r *http.Request) context.Context {
	return appengine.NewContext(r)
}

//Error is a custom error handler for the datastore code so the source of
//errors can be tracked down, as the source can get a bit deep.
type Error struct {
	Message string
	Origin  string
	Err     error
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Origin, e.Message)
}

//NewError is the creator of the new errors
func NewError(origin string, err error) error {
	return Error{
		Err:     err,
		Origin:  origin,
		Message: err.Error(),
	}
}
