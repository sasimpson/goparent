package mock

import (
	"context"

	"github.com/sasimpson/goparent"
)

//FamilyService -
type FamilyService struct {
	Env    *goparent.Env
	Kids   []*goparent.Child
	GetErr error
}

//Save -
func (mfs *FamilyService) Save(context.Context, *goparent.Family) error {
	panic("not implemented")
}

//Family -
func (mfs *FamilyService) Family(context.Context, string) (*goparent.Family, error) {
	panic("not implemented")
}

//Children -
func (mfs *FamilyService) Children(context.Context, *goparent.Family) ([]*goparent.Child, error) {
	if mfs.GetErr != nil {
		return nil, mfs.GetErr
	}
	return mfs.Kids, nil
}

//AddMember -
func (mfs *FamilyService) AddMember(context.Context, *goparent.Family, *goparent.User) error {
	panic("not implemented")
}

//GetAdminFamily -
func (mfs *FamilyService) GetAdminFamily(context.Context, *goparent.User) (*goparent.Family, error) {
	panic("not implemented")
}
