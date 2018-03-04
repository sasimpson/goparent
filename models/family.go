package models

import (
	"errors"
	"time"

	"github.com/sasimpson/goparent/config"
	"gopkg.in/gorethink/gorethink.v3"
)

//Family - structure to associate groups of parents with children
type Family struct {
	ID          string    `json:"id" gorethink:"id,omitempty"`
	Admin       string    `json:"admin" gorethink:"admin"`
	Members     []string  `json:"members" gorethink:"members"`
	CreatedAt   time.Time `json:"created_at" gorethink:"created_at"`
	LastUpdated time.Time `json:"last_updated" gorethink:"last_updated"`
}

func (family *Family) Save(env *config.Env) error {
	session, err := env.DB.GetConnection()
	if err != nil {
		return err
	}

	family.LastUpdated = time.Now()
	if family.ID == "" {
		family.CreatedAt = time.Now()
	}

	res, err := gorethink.Table("family").Insert(family, gorethink.InsertOpts{Conflict: "replace"}).RunWrite(session)
	if err != nil {
		return err
	}

	if res.Inserted > 0 {
		family.ID = res.GeneratedKeys[0]
	}
	return nil
}

func (family *Family) GetFamily(env *config.Env, id string) error {
	session, err := env.DB.GetConnection()
	if err != nil {
		return err
	}

	res, err := gorethink.Table("family").Get(id).Run(session)
	if err != nil {
		return err
	}

	err = res.One(family)
	if err != nil {
		return err
	}

	return nil
}

func (family *Family) GetAdminFamily(env *config.Env, id string) error {
	session, err := env.DB.GetConnection()
	if err != nil {
		return err
	}

	res, err := gorethink.Table("family").Filter(map[string]interface{}{
		"admin": id,
	}).Run(session)
	if err != nil {
		return err
	}

	err = res.One(family)
	if err != nil {
		return err
	}

	return nil
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

//AddMember - this will add a passed in user to a family.
func (family *Family) AddMember(env *config.Env, newMember *User) error {
	//check to see if they are already in the family, we don't want to add twice
	for _, member := range family.Members {
		if member == newMember.ID {
			return errors.New("user already in that family")
		}
	}

	family.Members = append(family.Members, newMember.ID)
	family.LastUpdated = time.Now()
	err := family.Save(env)
	if err != nil {
		return err
	}

	newMember.CurrentFamily = family.ID
	err = newMember.Save(env)
	if err != nil {
		return err
	}

	return nil
}
