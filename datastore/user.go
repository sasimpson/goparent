package datastore

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sasimpson/goparent"
	"google.golang.org/appengine/datastore"
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

//UserKind - constant string for all user entities in datastore
const UserKind = "User"

var (
	//ErrNoUserFound is when there is no user returned by the datastore
	ErrNoUserFound = errors.New("no result for that id")
	//ErrInvalidLogin is when the password/user combo do not match
	ErrInvalidLogin = errors.New("no result for that username password combo")
)

//User - get a user by the key/id
func (s *UserService) User(ctx context.Context, key string) (*goparent.User, error) {
	var user goparent.User
	userKey := datastore.NewKey(ctx, UserKind, key, 0, nil)
	err := datastore.Get(ctx, userKey, &user)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, NewError("datastore.UserService.User", ErrNoUserFound)
		}
		return nil, err
	}
	return &user, nil
}

//UserByLogin - get a user by their login, email and password
func (s *UserService) UserByLogin(ctx context.Context, email string, password string) (*goparent.User, error) {
	var user goparent.User
	userKey := datastore.NewKey(ctx, UserKind, md5Email(email), 0, nil)
	err := datastore.Get(ctx, userKey, &user)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, ErrInvalidLogin
		}
		return nil, NewError("datastore.UserService.User", err)
	}

	//verify the email and password match, then return it
	if user.Password == password && user.Email == email {
		return &user, nil
	}
	return nil, ErrInvalidLogin
}

//Save - save a user's current values
func (s *UserService) Save(ctx context.Context, user *goparent.User) error {
	userKey := datastore.NewKey(ctx, UserKind, md5Email(user.Email), 0, nil)
	//we use id as a way to lookup stuff, but don't actually _use_ id since we are using key...

	var family *goparent.Family
	fs := &FamilyService{Env: s.Env}
	//user doesn't have a current family
	if user.CurrentFamily == "" && user.ID != "" {
		log.Println("user with no current family but with an ID")
		//get the family for which the user is the admin
		family, err := fs.GetAdminFamily(ctx, user)
		if err == ErrNoFamilyFound {
			//didn't find a family, creating one.
			family = &goparent.Family{Admin: userKey.StringID(), Members: []string{userKey.StringID()}}
			err = fs.Save(ctx, family)
			if err != nil {
				return NewError("datastore.UserService.Save.1a", err)
			}
			log.Printf("%#v", family.ID)
		}
		if err != nil {
			return NewError("datastore.UserService.Save.1b", err)
		}
		log.Println("setting current family")
		user.CurrentFamily = family.ID
	}

	//if the user has no current family and no id, then we need to create a family.
	if user.CurrentFamily == "" && user.ID == "" {
		log.Println("user is new,  has no id or current family, so we're going to create one")
		family = &goparent.Family{Admin: userKey.StringID(), Members: []string{userKey.StringID()}}
		err := fs.Save(ctx, family)
		if err != nil {
			return NewError("datastore.UserService.Save.2", err)
		}
		user.CurrentFamily = family.ID
		log.Printf("%#v", family)
		log.Println("setting current family", family.ID)
	}

	user.ID = userKey.StringID()
	log.Printf("%#v", user)
	//this will save if there is or isn't a record for this user.
	_, err := datastore.Put(ctx, userKey, user)
	if err != nil {
		return NewError("datastore.UserService.Save.3", err)
	}
	return nil
}

//GetToken - gets the user token
func (s *UserService) GetToken(user *goparent.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["Name"] = user.Name
	claims["ID"] = user.ID
	claims["Email"] = user.Email
	claims["Username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	log.Printf("%s", s.Env.Auth.SigningKey)
	tokenString, err := token.SignedString(s.Env.Auth.SigningKey)
	if err != nil {
		return "", NewError("datastore.UserService.GetToken", err)
	}
	return tokenString, nil
}

//ValidateToken - validate token against signing method and populate user.
func (s *UserService) ValidateToken(ctx context.Context, tokenString string) (*goparent.User, bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return s.Env.Auth.SigningKey, nil
	})
	if err != nil {
		return nil, false, NewError("datastore.UserService.ValidateToken", err)
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		user, err := s.User(ctx, claims.ID)
		if err != nil {
			return nil, false, NewError("datastore.UserService.ValidateToken", err)
		}
		return user, true, nil
	}
	return nil, false, errors.New("invalid token")
}

//GetFamily -
func (s *UserService) GetFamily(ctx context.Context, user *goparent.User) (*goparent.Family, error) {
	if user.CurrentFamily == "" {
		return nil, errors.New("user has no current family")
	}

	fs := FamilyService{Env: s.Env}
	family, err := fs.Family(ctx, user.CurrentFamily)
	if err != nil {
		return nil, NewError("datastore.UserService.GetFamily", err)
	}
	return family, nil
}

//GetAllFamily -
func (s *UserService) GetAllFamily(ctx context.Context, user *goparent.User) ([]*goparent.Family, error) {
	var families []*goparent.Family
	q := datastore.NewQuery(FamilyKind).Filter("Members =", user.ID)
	t := q.Run(ctx)
	for {
		var family goparent.Family
		_, err := t.Next(&family)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		families = append(families, &family)
	}
	return families, nil
}

func md5Email(email string) string {
	h := md5.New()
	h.Write([]byte(email))
	return hex.EncodeToString(h.Sum(nil))
}
