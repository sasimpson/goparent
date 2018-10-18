package datastore

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/datastore"
)

//DBEnv -
type DBEnv struct {
	Client *datastore.Client
}

//GetConnection -
func (db *DBEnv) GetConnection() error {
	if db.Client != nil {
		return nil
	}

	ctx := db.GetContext(nil)
	dsClient, err := datastore.NewClient(ctx, "my-project")
	if err != nil {
		return err
	}
	db.Client = dsClient

	return nil
}

//GetContext - appengine requires a context from the request, but we don't
//want appengine code burried in the API.  This has been added to the interface,
//so it is abstracted out to the datastore bits.
func (db *DBEnv) GetContext(r *http.Request) context.Context {
	return context.Background()
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
