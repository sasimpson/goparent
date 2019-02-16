package mock

import (
	"context"
	"time"

	"github.com/sasimpson/goparent"
)

//UserService -
type UserService struct {
	Env          *goparent.Env
	Family       *goparent.Family
	ReturnedUser *goparent.User
	Token        string
	UserID       string
	AuthErr      error
	TokenErr     error
	FamilyErr    error
	SaveErr      error
}

//User -
func (m *UserService) User(context.Context, string) (*goparent.User, error) {
	panic("not implemented")
}

//UserByLogin -
func (m *UserService) UserByLogin(context.Context, string, string) (*goparent.User, error) {
	if m.AuthErr != nil {
		return nil, m.AuthErr
	}
	if m.ReturnedUser != nil {
		return m.ReturnedUser, nil
	}
	return nil, nil
}

//Save -
func (m *UserService) Save(ctx context.Context, user *goparent.User) error {
	if m.SaveErr != nil {
		return m.SaveErr
	}
	if m.UserID != "" {
		user.ID = m.UserID
	}
	return nil
}

//GetToken -
func (m *UserService) GetToken(*goparent.User, time.Duration) (string, error) {
	if m.TokenErr != nil {
		return "", m.TokenErr
	}
	if m.Token != "" {
		return m.Token, nil
	}
	return "", nil
}

//ValidateToken -
func (m *UserService) ValidateToken(context.Context, string) (*goparent.User, bool, error) {
	panic("not implemented")
}

//GetFamily -
func (m *UserService) GetFamily(context.Context, *goparent.User) (*goparent.Family, error) {
	if m.FamilyErr != nil {
		return nil, m.FamilyErr
	}
	return m.Family, nil
}

//GetAllFamily -
func (m *UserService) GetAllFamily(context.Context, *goparent.User) ([]*goparent.Family, error) {
	panic("not implemented")
}

//RequestResetPassword -
func (m *UserService) RequestResetPassword(context.Context, string, string) error {
	panic("not implemented")
}

//ResetPassword -
func (m *UserService) ResetPassword(context.Context, string, string) error {
	panic("not implemented")
}
