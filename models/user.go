package models

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sasimpson/goparent/config"
	gorethink "gopkg.in/gorethink/gorethink.v3"
)

//User - structure for storing user data
type User struct {
	ID       string `json:"id" gorethink:"id,omitempty"`
	Name     string `json:"name" gorethink:"name"`
	Email    string `json:"email" gorethink:"email"`
	Username string `json:"username" gorethink:"username"`
	Password string `json:"-" gorethink:"password"`
}

type UserClaims struct {
	ID       string
	Name     string
	Email    string
	Username string
	Password string
	jwt.StandardClaims
}

//GetUser - gets the user data based on the id string
func (user *User) GetUser(env *config.Env, id string) error {
	session, err := env.DB.GetConnection()
	if err != nil {
		return err
	}
	res, err := gorethink.Table("users").Get(id).Run(session)
	if err != nil {
		return err
	}
	defer res.Close()
	if res.IsNil() {
		return errors.New("no result for that id")
	}
	res.One(&user)
	return nil
}

func (user *User) GetUserByLogin(env *config.Env, username string, password string) error {
	//TODO: need to hash the password
	session, err := env.DB.GetConnection()
	if err != nil {
		return err
	}
	res, err := gorethink.Table("users").Filter(map[string]interface{}{
		"email":    username,
		"password": password}).Run(session)
	if err != nil {
		return err
	}
	defer res.Close()
	if res.IsNil() {
		return errors.New("no result for that username password combo")
	}
	res.One(&user)
	return nil
}

//Save - saves the user. creates it if it doesn't exist.  upsert only works if there is an id and that email exists.
func (user *User) Save(env *config.Env) error {
	session, err := env.DB.GetConnection()
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
	if res.IsNil() || user.ID != "" {
		res2, err := gorethink.Table("users").Insert(user, gorethink.InsertOpts{Conflict: "replace"}).RunWrite(session)
		if err != nil {
			log.Println("error with insert from users upsert in user.Save()")
			return err
		}
		if res2.Inserted > 0 {
			user.ID = res2.GeneratedKeys[0]
		}
		return nil
	}
	return errors.New("there needs to be an ID in the user if one with that email exists")
}

func (user *User) GetToken(env *config.Env) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["Name"] = user.Name
	claims["ID"] = user.ID
	claims["Email"] = user.Email
	claims["Username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString(env.Auth.SigningKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (user *User) ValidateToken(env *config.Env, tokenString string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return env.Auth.SigningKey, nil
	})
	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		user.GetUser(env, claims.ID)
		return true, nil
	}
	return false, errors.New("invalid token")
}

//UserFromContext - helper to get the user from the request context
func UserFromContext(ctx context.Context) (User, error) {
	user, ok := ctx.Value("user").(User)
	if !ok {
		return User{}, errors.New("no user found in context")
	}
	return user, nil
}
