package models

import (
	"errors"
	"fmt"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"gopkg.in/gorethink/gorethink.v3"
)

//User - structure for storing user data
type User struct {
	ID       string `json:"id" gorethink:"id,omitempty"`
	Name     string `json:"name" gorethink:"name"`
	Email    string `json:"email" gorethink:"email"`
	Username string `json:"username" gorethink:"username"`
	Password string `json:"password" gorethink:"password"`
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
func (user *User) GetUser(id string) error {
	session, err := GetConnection()
	if err != nil {
		return err
	}
	defer session.Close()
	resp, err := gorethink.Table("users").Get(id).Run(session)
	if err != nil {
		return err
	}
	if resp.IsNil() {
		return errors.New("no result for that id")
	}
	resp.One(&user)
	return nil
}

func (user *User) GetUserByLogin(username string, password string) error {
	//TODO: need to hash the password
	session, err := GetConnection()
	if err != nil {
		return err
	}
	defer session.Close()
	resp, err := gorethink.Table("users").Filter(map[string]interface{}{
		"email":    username,
		"password": password}).Run(session)
	if err != nil {
		return err
	}
	if resp.IsNil() {
		return errors.New("no result for that username password combo")
	}
	resp.One(&user)
	return nil
}

//Save - saves the user. creates it if it doesn't exist.  upsert only works if there is an id and that email exists.
func (user *User) Save() error {
	session, err := GetConnection()
	if err != nil {
		return err
	}
	defer session.Close()
	resp, err := gorethink.Table("users").Filter(map[string]interface{}{
		"email": user.Email,
	}).Run(session)
	if err != nil {
		return err
	}
	defer resp.Close()
	if resp.IsNil() || user.ID != "" {
		resp2, err := gorethink.Table("users").Insert(user, gorethink.InsertOpts{Conflict: "replace"}).RunWrite(session)
		if err != nil {
			log.Println("error with insert from users upsert in user.Save()")
			return err
		}
		if resp2.Inserted > 0 {
			user.ID = resp2.GeneratedKeys[0]
		}
		return nil
	}
	return errors.New("there needs to be an ID in the user if one with that email exists")
}

func (user *User) GetToken(mySigningKey []byte) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["Name"] = user.Name
	claims["ID"] = user.ID
	claims["Email"] = user.Email
	claims["Username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (user *User) ValidateToken(tokenString string, mySigningKey []byte) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return mySigningKey, nil
	})
	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		user.GetUser(claims.ID)
		return true, nil
	}
	return false, errors.New("invalid token")
}
