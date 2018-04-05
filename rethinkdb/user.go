package rethinkdb

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/config"
	gorethink "gopkg.in/gorethink/gorethink.v3"
)

//ErrExistingInvitation - the user already has an invitation
const ErrExistingInvitation string = "existing invitation for that user"

type UserService struct {
	Env *config.Env
}

//UserClaims - structure for inserting claims into a jwt auth token
type UserClaims struct {
	ID       string
	Name     string
	Email    string
	Username string
	Password string
	jwt.StandardClaims
}

//GetUser - gets the user data based on the id string
func (us *UserService) User(id string) (*goparent.User, error) {
	session, err := us.Env.DB.GetConnection()
	if err != nil {
		return nil, err
	}
	res, err := gorethink.Table("users").Get(id).Run(session)
	if err != nil {
		return nil, err
	}

	defer res.Close()
	if res.IsNil() {
		return nil, errors.New("no result for that id")
	}

	var user goparent.User
	res.One(&user)
	return &user, nil
}

//GetUserByLogin - gets a user by their username and password
func (us *UserService) UserByLogin(username string, password string) (*goparent.User, error) {
	//TODO: need to hash the password
	session, err := us.Env.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	res, err := gorethink.Table("users").Filter(map[string]interface{}{
		"email":    username,
		"password": password}).Run(session)
	if err != nil {
		return nil, err
	}

	defer res.Close()
	if res.IsNil() {
		return nil, errors.New("no result for that username password combo")
	}

	var user goparent.User
	err = res.One(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

//Save - saves the user. creates it if it doesn't exist.  upsert only works if there is an id and that email exists.
func (us *UserService) Save(user *goparent.User) error {
	session, err := us.Env.DB.GetConnection()
	if err != nil {
		return err
	}

	res, err := gorethink.Table("users").Filter(map[string]interface{}{
		"email": user.Email,
	}).Run(session)
	if err != nil {
		return err
	}
	defer res.Close()

	var family *goparent.Family
	fs := &FamilyService{Env: us.Env}
	if user.CurrentFamily == "" && user.ID != "" {
		family, err := fs.GetAdminFamily(user)
		if err != nil {
			if err == gorethink.ErrEmptyResult {
				family = &goparent.Family{Admin: user.ID, Members: []string{user.ID}}
				fs.Save(family)
			} else {
				return err
			}

		}
		user.CurrentFamily = family.ID
	}

	if user.CurrentFamily == "" && user.ID == "" {
		family = &goparent.Family{Admin: user.ID, Members: []string{user.ID}}
		fs.Save(family)
		user.CurrentFamily = family.ID
	}

	if res.IsNil() || user.ID != "" {
		res2, err := gorethink.Table("users").Insert(user, gorethink.InsertOpts{Conflict: "replace"}).RunWrite(session)
		if err != nil {
			return err
		}

		if res2.Inserted > 0 {
			user.ID = res2.GeneratedKeys[0]
			family.Admin = user.ID
			family.Members = []string{user.ID}
			fs.Save(family)
		}
		return nil
	}
	return errors.New("there needs to be an ID in the user if one with that email exists")
}

//GetToken - gets the user token
func (us *UserService) GetToken(user *goparent.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["Name"] = user.Name
	claims["ID"] = user.ID
	claims["Email"] = user.Email
	claims["Username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString(us.Env.Auth.SigningKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

//ValidateToken - validate token against signing method and populate user.
func (us *UserService) ValidateToken(tokenString string) (*goparent.User, bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return us.Env.Auth.SigningKey, nil
	})
	if err != nil {
		return nil, false, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		user, err := us.User(claims.ID)
		if err != nil {
			return nil, false, err
		}
		return user, true, nil
	}
	return nil, false, errors.New("invalid token")
}

//GetFamily - return the family for a user. used for lookups
func (us *UserService) GetFamily(user *goparent.User) (*goparent.Family, error) {
	if user.CurrentFamily == "" {
		return nil, errors.New("user has no current family")
	}

	fs := FamilyService{Env: us.Env}
	family, err := fs.Family(user.CurrentFamily)
	if err != nil {
		return nil, err
	}
	return family, nil
}

//GetAllFamily - return the family for a user. used for lookups
func (us *UserService) GetAllFamily(user *goparent.User) ([]*goparent.Family, error) {
	session, err := us.Env.DB.GetConnection()
	if err != nil {
		return nil, err
	}
	/*  rethinkdb query:
	r.db('goparent').table('family').filter(function(famRow) {
		return famRow('members').contains("userid")
	})
	*/
	res, err := gorethink.Table("family").
		Filter(
			func(row gorethink.Term) gorethink.Term {
				return row.Field("members").Contains(user.ID)
			},
		).
		Run(session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var family []*goparent.Family
	err = res.All(&family)
	if err != nil {
		return nil, err
	}
	return family, nil
}
