package rethinkdb

import (
	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/config"
	gorethink "gopkg.in/gorethink/gorethink.v3"
)

//ChildService - service for implementing the interface
type ChildService struct {
	Env *config.Env
}

//Save - Create or update child record
func (cs *ChildService) Save(child *goparent.Child) error {
	session, err := cs.Env.DB.GetConnection()
	if err != nil {
		return err
	}

	res, err := gorethink.Table("children").Insert(child, gorethink.InsertOpts{Conflict: "replace"}).RunWrite(session)
	if err != nil {
		return err
	}

	if res.Inserted > 0 {
		child.ID = res.GeneratedKeys[0]
	}
	return nil
}

//Child - return a child for an ID
func (cs *ChildService) Child(id string) (*goparent.Child, error) {

	session, err := cs.Env.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	res, err := gorethink.Table("children").Get(id).Run(session)
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
func (cs *ChildService) Delete(child *goparent.Child) (int, error) {
	session, err := cs.Env.DB.GetConnection()
	if err != nil {
		return 0, err
	}

	res, err := gorethink.Table("children").Get(child.ID).Delete().RunWrite(session)
	if err != nil {
		return 0, err
	}
	return res.Deleted, nil
}
