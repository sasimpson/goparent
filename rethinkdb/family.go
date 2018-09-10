package rethinkdb

import (
	"context"
	"errors"
	"time"

	"github.com/sasimpson/goparent"
	"gopkg.in/gorethink/gorethink.v3"
)

//FamilyService - structure for replicating the interface
type FamilyService struct {
	Env *goparent.Env
	DB  *DBEnv
}

//Save - Create or Update a family record
func (fs *FamilyService) Save(ctx context.Context, family *goparent.Family) error {
	err := fs.DB.GetConnection()
	if err != nil {
		return err
	}

	family.LastUpdated = time.Now()
	if family.ID == "" {
		family.CreatedAt = time.Now()
	}

	res, err := gorethink.Table("family").Insert(family, gorethink.InsertOpts{Conflict: "replace"}).RunWrite(fs.DB.Session)
	if err != nil {
		return err
	}

	if res.Inserted > 0 {
		family.ID = res.GeneratedKeys[0]
	}
	return nil
}

//Family - returns a family for an ID
func (fs *FamilyService) Family(ctx context.Context, id string) (*goparent.Family, error) {
	err := fs.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	res, err := gorethink.Table("family").Get(id).Run(fs.DB.Session)
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
	err := fs.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	res, err := gorethink.Table("children").Filter(map[string]interface{}{"familyID": family.ID}).OrderBy(gorethink.Desc("birthday")).Run(fs.DB.Session)
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
func (fs *FamilyService) AddMember(ctx context.Context, family *goparent.Family, newMember *goparent.User) error {
	//check to see if they are already in the family, we don't want to add twice
	for _, member := range family.Members {
		if member == newMember.ID {
			return errors.New("user already in that family")
		}
	}

	family.Members = append(family.Members, newMember.ID)
	family.LastUpdated = time.Now()
	err := fs.Save(ctx, family)
	if err != nil {
		return err
	}

	return nil
}

//GetAdminFamily - returns the family for which the user is the admin.
func (fs *FamilyService) GetAdminFamily(ctx context.Context, user *goparent.User) (*goparent.Family, error) {
	err := fs.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	res, err := gorethink.Table("family").Filter(map[string]interface{}{
		"admin": user.ID,
	}).Run(fs.DB.Session)
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
