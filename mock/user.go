package mock

import (
	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/config"
)

type MockUserService struct {
	Env          *config.Env
	Family       *goparent.Family
	ReturnedUser *goparent.User
	Token        string
	UserID       string
	AuthErr      error
	TokenErr     error
	FamilyErr    error
	SaveErr      error
}

func (m *MockUserService) User(string) (*goparent.User, error) {
	panic("not implemented")
}

func (m *MockUserService) UserByLogin(string, string) (*goparent.User, error) {
	if m.AuthErr != nil {
		return nil, m.AuthErr
	}
	if m.ReturnedUser != nil {
		return m.ReturnedUser, nil
	}
	return nil, nil
}

func (m *MockUserService) Save(user *goparent.User) error {
	if m.SaveErr != nil {
		return m.SaveErr
	}
	if m.UserID != "" {
		user.ID = m.UserID
	}
	return nil
}

func (m *MockUserService) GetToken(*goparent.User) (string, error) {
	if m.TokenErr != nil {
		return "", m.TokenErr
	}
	if m.Token != "" {
		return m.Token, nil
	}
	return "", nil
}

func (m *MockUserService) ValidateToken(string) (*goparent.User, bool, error) {
	panic("not implemented")
}

func (m *MockUserService) GetFamily(*goparent.User) (*goparent.Family, error) {
	if m.FamilyErr != nil {
		return nil, m.FamilyErr
	}
	return m.Family, nil
}

func (m *MockUserService) GetAllFamily(*goparent.User) ([]*goparent.Family, error) {
	panic("not implemented")
}
