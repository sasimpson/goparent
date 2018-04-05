package rethinkdb

import (
	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/config"
	gorethink "gopkg.in/gorethink/gorethink.v3"
)

type ChildService struct {
	Env *config.Env
}

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

func (cs *ChildService) Child(id string) (*goparent.Child, error) {

	session, err := cs.Env.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	res, err := gorethink.Table("children").
		Filter(
			map[string]interface{}{
				"id": id,
			},
		).Run(session)
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

func (cs *ChildService) Delete(child *goparent.Child) (int, error) {
	session, err := cs.Env.DB.GetConnection()
	if err != nil {
		return 0, err
	}

	res, err := gorethink.Table("children").
		Filter(map[string]interface{}{"id": child.ID}).Delete().RunWrite(session)
	if err != nil {
		return 0, err
	}
	return res.Deleted, nil
}
