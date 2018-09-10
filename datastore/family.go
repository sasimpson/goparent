package datastore

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sasimpson/goparent"
	"google.golang.org/appengine/datastore"
)

//FamilyService -
type FamilyService struct {
	Env *goparent.Env
}

//FamilyKind - constant string for all family entities in datastore
const FamilyKind = "Family"

var (
	//ErrNoFamilyFound is when no family is found
	ErrNoFamilyFound = errors.New("no family found with user as admin")
)

//Save -
func (s *FamilyService) Save(ctx context.Context, family *goparent.Family) error {
	var familyKey *datastore.Key
	//if the family id is blank then we are creating a new family
	if family.ID == "" {
		u := uuid.New()
		familyKey = datastore.NewKey(ctx, FamilyKind, u.String(), 0, nil)
		family.CreatedAt = time.Now()
		family.ID = u.String()
	} else {
		familyKey = datastore.NewKey(ctx, FamilyKind, family.ID, 0, nil)
	}

	family.LastUpdated = time.Now()

	_, err := datastore.Put(ctx, familyKey, &family)
	if err != nil {
		return NewError("FamilyService.Save", err)
	}

	return nil
}

//Family -
func (s *FamilyService) Family(ctx context.Context, id string) (*goparent.Family, error) {
	var family goparent.Family
	familyKey := datastore.NewKey(ctx, FamilyKind, id, 0, nil)
	err := datastore.Get(ctx, familyKey, &family)
	if err != nil {
		return nil, NewError("FamilyService.Family", err)
	}

	return &family, nil
}

//Children -
func (s *FamilyService) Children(*goparent.Family) ([]*goparent.Child, error) {
	panic("not implemented")
}

//AddMember -
func (s *FamilyService) AddMember(*goparent.Family, *goparent.User) error {
	panic("not implemented")
}

//GetAdminFamily -
func (s *FamilyService) GetAdminFamily(ctx context.Context, user *goparent.User) (*goparent.Family, error) {
	var families []goparent.Family
	q := datastore.NewQuery(FamilyKind).Filter("Admin=", user.ID)
	_, err := q.GetAll(ctx, &families)
	if err != nil {
		return nil, NewError("FamilyService.GetAdminFamily", err)
	}

	if len(families) < 1 {
		return nil, NewError("FamilyService.GetAdminFamily", ErrNoFamilyFound)
	}

	return &families[0], nil
}
