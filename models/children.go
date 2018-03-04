package models

import (
	"time"

	"github.com/sasimpson/goparent/config"
	gorethink "gopkg.in/gorethink/gorethink.v3"
)

//Child - structure for child data
type Child struct {
	ID       string    `json:"id" gorethink:"id,omitempty"`
	Name     string    `json:"name" gorethink:"name"`
	ParentID string    `json:"parentID" gorethink:"parentID"`
	FamilyID string    `json:"familyID" gorethink:"familyID"`
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

//GetAllChildren - get all children records for a user from the datastore
func GetAllChildren(env *config.Env, user *User) ([]Child, error) {
	family, err := user.GetFamily(env)
	if err != nil {
		return nil, err
	}
	children, err := family.GetAllChildren(env)
	return children, err
}

//GetChild - data for a child based on the user (parent) and the child id
func (child *Child) GetChild(env *config.Env, user *User, childID string) error {
	family, err := user.GetFamily(env)
	if err != nil {
		return err
	}

	session, err := env.DB.GetConnection()
	if err != nil {
		return err
	}

	res, err := gorethink.Table("children").
		Filter(
			map[string]interface{}{
				"familyID": family.ID,
				"id":       childID,
			},
		).Run(session)
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

//DeleteChild - delete a child based on the user and the child
func (child *Child) DeleteChild(env *config.Env, user *User) (int, error) {
	family, err := user.GetFamily(env)
	if err != nil {
		return 0, err
	}

	session, err := env.DB.GetConnection()
	if err != nil {
		return 0, err
	}

	res, err := gorethink.Table("children").
		Filter(map[string]interface{}{"familyID": family.ID, "id": child.ID}).Delete().RunWrite(session)
	if err != nil {
		return 0, err
	}
	return res.Deleted, nil
}
