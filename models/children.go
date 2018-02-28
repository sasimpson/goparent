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
	return children, nil
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

//Family - structure to associate groups of parents with children
type Family struct {
	ID          string    `json:"id" gorethink:"id,omitempty"`
	Admin       string    `json:"admin" gorethink:"admin"`
	Members     []string  `json:"members" gorethink:"members"`
	CreatedAt   time.Time `json:"created_at" gorethink:"created_at"`
	LastUpdated time.Time `json:"last_updated" gorethink:"last_updated"`
}

//GetAllChildren - returns all the children for a family
func (family *Family) GetAllChildren(env *config.Env) ([]Child, error) {
	session, err := env.DB.GetConnection()
	if err != nil {
		return nil, err
	}
	res, err := gorethink.Table("children").Filter(map[string]interface{}{"familyID": family.ID}).OrderBy(gorethink.Desc("birthday")).Run(session)
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
