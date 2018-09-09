package datastore

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
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

//UserEntity - constant string for all user entities in datastore
const UserEntity = "User"

//User - get a user by the key/id
func (s *UserService) User(ctx context.Context, key string) (*goparent.User, error) {
	var user goparent.User
	userKey := datastore.NewKey(ctx, UserEntity, key, 0, nil)
	err := datastore.Get(ctx, userKey, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

//UserByLogin - get a user by their login, email and password
func (s *UserService) UserByLogin(ctx context.Context, email string, password string) (*goparent.User, error) {
	var user goparent.User
	userKey := datastore.NewKey(ctx, UserEntity, md5Email(email), 0, nil)
	err := datastore.Get(ctx, userKey, &user)
	if err != nil {
		return nil, err
	}

	//verify the email and password match, then return it
	if user.Password == password && user.Email == email {
		return &user, nil
	}
	return nil, errors.New("no result for that username password combo")
}

//Save - save a user's current values
func (s *UserService) Save(user *goparent.User) error {
	ctx := context.Background()
	userKey := datastore.NewKey(ctx, UserEntity, md5Email(user.Email), 0, nil)
	err := datastore.Get(ctx, userKey, &user)
	//check for error or no such entity, we want to continue if it exists or it doesn't
	if err != nil && err != datastore.ErrNoSuchEntity {
		return err
	}
	user.ID = userKey.String()
	_, err = datastore.Put(ctx, userKey, user)
	if err != nil {
		return err
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

	tokenString, err := token.SignedString(s.Env.Auth.SigningKey)
	if err != nil {
		return "", err
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
		return nil, false, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		user, err := s.User(ctx, claims.ID)
		if err != nil {
			return nil, false, err
		}
		return user, true, nil
	}
	return nil, false, errors.New("invalid token")
}

//GetFamily -
func (s *UserService) GetFamily(user *goparent.User) (*goparent.Family, error) {
	if user.CurrentFamily == "" {
		return nil, errors.New("user has no current family")
	}

	panic("not implemented")
}

//GetAllFamily -
func (s *UserService) GetAllFamily(*goparent.User) ([]*goparent.Family, error) {
	panic("not implemented")
}

func md5Email(email string) string {
	h := md5.New()
	h.Write([]byte(email))
	return hex.EncodeToString(h.Sum(nil))
}
