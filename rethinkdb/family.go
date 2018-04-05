package rethinkdb

import (
	"errors"
	"time"

	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/config"
	"gopkg.in/gorethink/gorethink.v3"
)

type FamilyService struct {
	Env *config.Env
}

func (fs *FamilyService) Save(family *goparent.Family) error {
	session, err := fs.Env.DB.GetConnection()
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

func (fs *FamilyService) Family(id string) (*goparent.Family, error) {
	session, err := fs.Env.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	res, err := gorethink.Table("family").Get(id).Run(session)
	if err != nil {
		return nil, err
	}
	var family goparent.Family
	err = res.One(&family)
	if err != nil {
		return nil, err
	}

	return &family, nil
}

//Children - returns all the children for a family
func (fs *FamilyService) Children(family *goparent.Family) ([]*goparent.Child, error) {
	session, err := fs.Env.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	res, err := gorethink.Table("children").Filter(map[string]interface{}{"familyID": family.ID}).OrderBy(gorethink.Desc("birthday")).Run(session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var rows []*goparent.Child
	err = res.All(&rows)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

//AddMember - this will add a passed in user to a family.
func (fs *FamilyService) AddMember(family *goparent.Family, newMember *goparent.User) error {
	//check to see if they are already in the family, we don't want to add twice
	for _, member := range family.Members {
		if member == newMember.ID {
			return errors.New("user already in that family")
		}
	}

	family.Members = append(family.Members, newMember.ID)
	family.LastUpdated = time.Now()
	err := fs.Save(family)
	if err != nil {
		return err
	}

	// newMember.CurrentFamily = family.ID
	// err = newMember.Save(env)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (fs *FamilyService) GetAdminFamily(user *goparent.User) (*goparent.Family, error) {
	session, err := fs.Env.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	res, err := gorethink.Table("family").Filter(map[string]interface{}{
		"admin": user.ID,
	}).Run(session)
	if err != nil {
		return nil, err
	}

	var family goparent.Family
	err = res.One(family)
	if err != nil {
		return nil, err
	}

	return &family, nil
}
