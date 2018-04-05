package mock

import (
	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/config"
)

type MockUserService struct {
	Env    *config.Env
	Family *goparent.Family
}

func (mus *MockUserService) User(string) (*goparent.User, error) {
	panic("not implemented")
}

func (mus *MockUserService) UserByLogin(string, string) (*goparent.User, error) {
	panic("not implemented")
}

func (mus *MockUserService) Save(*goparent.User) error {
	panic("not implemented")
}

func (mus *MockUserService) GetToken(*goparent.User) (string, error) {
	panic("not implemented")
}

func (mus *MockUserService) ValidateToken(string) (*goparent.User, bool, error) {
	panic("not implemented")
}

func (mus *MockUserService) GetFamily(*goparent.User) (*goparent.Family, error) {
	return mus.Family, nil
}

func (mus *MockUserService) GetAllFamily(*goparent.User) ([]*goparent.Family, error) {
	panic("not implemented")
}
