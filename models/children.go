package models

import (
	"time"

	"github.com/sasimpson/goparent/config"
	gorethink "gopkg.in/gorethink/gorethink.v3"
)

type Child struct {
	ID       string    `json:"id" gorethink:"id,omitempty"`
	Name     string    `json:"name" gorethink:"name"`
	ParentID string    `json:"parentID" gorethink:"parentID"`
	Birthday time.Time `json:"birthday" gorethink:"birthday"`
}

//Save - save the structure to the datastore
func (child *Child) Save(env *config.Env) error {
	session, err := env.DB.GetConnection()
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

//GetAllChildren - get all records for a user from the datastore
func GetAllChildren(env *config.Env, user *User) ([]Child, error) {
	session, err := env.DB.GetConnection()
	if err != nil {
		return nil, err
	}
	res, err := gorethink.Table("children").Filter(map[string]interface{}{"parentID": user.ID}).OrderBy(gorethink.Desc("birthday")).Run(session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var rows []Child
	err = res.All(&rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (child *Child) GetChild(env *config.Env, user *User, childID string) error {
	session, err := env.DB.GetConnection()
	if err != nil {
		return err
	}

	res, err := gorethink.Table("children").Filter(map[string]interface{}{"parentID": user.ID, "id": childID}).Run(session)
	if err != nil {
		return err
	}
	defer res.Close()

	err = res.One(&child)
	if err != nil {
		return err
	}

	return nil
}

func (child *Child) DeleteChild(env *config.Env, user *User) (int, error) {
	session, err := env.DB.GetConnection()
	if err != nil {
		return 0, err
	}

	res, err := gorethink.Table("children").Filter(map[string]interface{}{"parentID": user.ID, "id": child.ID}).Delete().RunWrite(session)
	if err != nil {
		return 0, err
	}
	return res.Deleted, nil
}
