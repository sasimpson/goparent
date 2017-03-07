package models

import (
	"errors"
	"log"

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

//GetUser - gets the user data based on the id string
func GetUser(id string) (User, error) {
	session, err := GetConnection()
	if err != nil {
		return User{}, err
	}
	defer session.Close()
	resp, err := gorethink.Table("users").Get(id).Run(session)
	if err != nil {
		return User{}, err
	}
	if resp.IsNil() {
		return User{}, errors.New("no result for that id")
	}
	var user User
	resp.One(&user)
	return user, nil
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
			log.Println("error with return from users get in user.Save() 1")
			return err
		}
		log.Printf("resp2: %v", resp2)
		if resp2.Inserted > 0 {
			user.ID = resp2.GeneratedKeys[0]
		}
		return nil
	}
	return errors.New("there needs to be an ID in the user if one with that email exists")
}

//wendy_burruel@yahoo.com
