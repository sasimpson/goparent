package rethinkdb

import (
	"context"

	"github.com/sasimpson/goparent"
	gorethink "gopkg.in/gorethink/gorethink.v3"
)

//ChildService - service for implementing the interface
type ChildService struct {
	Env *goparent.Env
	DB  *DBEnv
}

//Save - Create or update child record
func (cs *ChildService) Save(ctx context.Context, child *goparent.Child) error {
	err := cs.DB.GetConnection()
	if err != nil {
		return err
	}

	res, err := gorethink.Table("children").Insert(child, gorethink.InsertOpts{Conflict: "replace"}).RunWrite(cs.DB.Session)
	if err != nil {
		return err
	}

	if res.Inserted > 0 {
		child.ID = res.GeneratedKeys[0]
	}
	return nil
}

//Child - return a child for an ID
func (cs *ChildService) Child(ctx context.Context, id string) (*goparent.Child, error) {

	err := cs.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	res, err := gorethink.Table("children").Get(id).Run(cs.DB.Session)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	var child goparent.Child
	err = res.One(&child)
	if err != nil {
		return nil, err
	}

	return &child, nil
}

//Delete - delete a passed child record from the datastore
func (cs *ChildService) Delete(ctx context.Context, child *goparent.Child) (int, error) {
	err := cs.DB.GetConnection()
	if err != nil {
		return 0, err
	}

	res, err := gorethink.Table("children").Get(child.ID).Delete().RunWrite(cs.DB.Session)
	if err != nil {
		return 0, err
	}
	return res.Deleted, nil
}
