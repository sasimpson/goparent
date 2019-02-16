package rethinkdb

import (
	"context"
	"errors"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sasimpson/goparent"
	gorethink "gopkg.in/gorethink/gorethink.v3"
)

//UserService -
type UserService struct {
	Env *goparent.Env
	DB  *DBEnv
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

//User - gets the user data based on the id string
func (us *UserService) User(ctx context.Context, id string) (*goparent.User, error) {
	err := us.DB.GetConnection()
	if err != nil {
		return nil, err
	}
	res, err := gorethink.Table("users").Get(id).Run(us.DB.Session)
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

//UserByLogin - gets a user by their username and password
func (us *UserService) UserByLogin(ctx context.Context, username string, password string) (*goparent.User, error) {
	//TODO: need to hash the password
	err := us.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	res, err := gorethink.Table("users").Filter(map[string]interface{}{
		"email":    username,
		"password": password}).Run(us.DB.Session)
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
func (us *UserService) Save(ctx context.Context, user *goparent.User) error {
	err := us.DB.GetConnection()
	if err != nil {
		return err
	}

	//check to see if a user with that email exists already
	res, err := gorethink.Table("users").Filter(map[string]interface{}{
		"email": user.Email,
	}).Run(us.DB.Session)
	if err != nil {
		return err
	}
	defer res.Close()

	var family *goparent.Family
	fs := &FamilyService{Env: us.Env, DB: us.DB}
	//if the user doesn't have a current family
	if user.CurrentFamily == "" && user.ID != "" {
		//get the family for which the user is the admin.
		family, err := fs.GetAdminFamily(ctx, user)
		if err != nil {
			//if no result is returned, create a new family
			if err == gorethink.ErrEmptyResult {
				family = &goparent.Family{Admin: user.ID, Members: []string{user.ID}}
				fs.Save(ctx, family)
			} else {
				return err
			}

		}
		user.CurrentFamily = family.ID
	}

	//if the user has no current family and no id, then we need to create a family.
	if user.CurrentFamily == "" && user.ID == "" {
		family = &goparent.Family{}
		fs.Save(ctx, family)
		user.CurrentFamily = family.ID
	}

	//if the user doesn't exist in the db OR the ID exists, save to db
	if res.IsNil() || user.ID != "" {
		res2, err := gorethink.Table("users").Insert(user, gorethink.InsertOpts{Conflict: "replace"}).RunWrite(us.DB.Session)
		if err != nil {
			return err
		}

		//if the user was inserted and not replaced, set the ID and family bits.s
		if res2.Inserted > 0 {
			user.ID = res2.GeneratedKeys[0]
			family.Admin = user.ID
			family.Members = []string{user.ID}
			fs.Save(ctx, family)
		}
		return nil
	}
	return errors.New("there needs to be an ID in the user if one with that email exists")
}

//GetToken - gets the user token
func (us *UserService) GetToken(user *goparent.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["Name"] = user.Name
	claims["ID"] = user.ID
	claims["Email"] = user.Email
	claims["Username"] = user.Username
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString(us.Env.Auth.SigningKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

//ValidateToken - validate token against signing method and populate user.
func (us *UserService) ValidateToken(ctx context.Context, tokenString string) (*goparent.User, bool, error) {
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
		user, err := us.User(ctx, claims.ID)
		if err != nil {
			return nil, false, err
		}
		return user, true, nil
	}
	return nil, false, errors.New("invalid token")
}

//GetFamily - return the family for a user. used for lookups
func (us *UserService) GetFamily(ctx context.Context, user *goparent.User) (*goparent.Family, error) {
	if user.CurrentFamily == "" {
		return nil, errors.New("user has no current family")
	}

	fs := FamilyService{Env: us.Env, DB: us.DB}
	family, err := fs.Family(ctx, user.CurrentFamily)
	if err != nil {
		return nil, err
	}
	return family, nil
}

//GetAllFamily - return the family for a user. used for lookups
func (us *UserService) GetAllFamily(ctx context.Context, user *goparent.User) ([]*goparent.Family, error) {
	err := us.DB.GetConnection()
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
		Run(us.DB.Session)
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

//RequestResetPassword - needs implementation
func (us *UserService) RequestResetPassword(context.Context, string, string) error {
	return nil
}

//ResetPassword - needs implementation
func (us *UserService) ResetPassword(context.Context, string, string) error {
	return nil
}
