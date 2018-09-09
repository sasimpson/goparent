package datastore

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/sasimpson/goparent"
)

//UserService -
type UserService struct {
	Env *goparent.Env
}

//UserClaims -
type UserClaims struct {
	ID       string
	Name     string
	Email    string
	Username string
	Password string
	jwt.StandardClaims
}

//User -
func (s *UserService) User(string) (*goparent.User, error) {
	panic("not implemented")
}

//UserByLogin -
func (s *UserService) UserByLogin(string, string) (*goparent.User, error) {
	panic("not implemented")
}

//Save -
func (s *UserService) Save(*goparent.User) error {
	panic("not implemented")
}

//GetToken -
func (s *UserService) GetToken(*goparent.User) (string, error) {
	panic("not implemented")
}

//ValidateToken -
func (s *UserService) ValidateToken(string) (*goparent.User, bool, error) {
	panic("not implemented")
}

//GetFamily -
func (s *UserService) GetFamily(*goparent.User) (*goparent.Family, error) {
	panic("not implemented")
}

//GetAllFamily -
func (s *UserService) GetAllFamily(*goparent.User) ([]*goparent.Family, error) {
	panic("not implemented")
}
