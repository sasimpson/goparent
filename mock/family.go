package mock

import (
	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/config"
)

type MockFamilyService struct {
	Env    *config.Env
	Kids   []*goparent.Child
	GetErr error
}

func (mfs *MockFamilyService) Save(*goparent.Family) error {
	panic("not implemented")
}

func (mfs *MockFamilyService) Family(string) (*goparent.Family, error) {
	panic("not implemented")
}

func (mfs *MockFamilyService) Children(*goparent.Family) ([]*goparent.Child, error) {
	if mfs.GetErr != nil {
		return nil, mfs.GetErr
	}
	return mfs.Kids, nil
}

func (mfs *MockFamilyService) AddMember(*goparent.Family, *goparent.User) error {
	panic("not implemented")
}

func (mfs *MockFamilyService) GetAdminFamily(*goparent.User) (*goparent.Family, error) {
	panic("not implemented")
}
