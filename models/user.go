package models

import (
	"errors"
	"fmt"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sasimpson/goparent/config"
	gorethink "gopkg.in/gorethink/gorethink.v3"
)

const ErrExistingInvitation string = "existing invitation for that user"

//User - structure for storing user data
type User struct {
	ID       string `json:"id" gorethink:"id,omitempty"`
	Name     string `json:"name" gorethink:"name"`
	Email    string `json:"email" gorethink:"email"`
	Username string `json:"username" gorethink:"username"`
	Password string `json:"-" gorethink:"password"`
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

//GetUserByLogin - gets a user by their username and password
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

//GetToken - gets the user token
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

//ValidateToken - validate token against signing method and populate user.
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

//UserInvitation - structure for storing invitations
type UserInvitation struct {
	ID          string    `json:"id" gorethink:"id,omitempty"`
	UserID      string    `json:"userID" gorethink:"userID"`
	InviteEmail string    `json:"inviteEmail" gorethink:"inviteEmail"`
	Timestamp   time.Time `json:"timestamp" gorethink:"timestamp"`
}

//InviteParent - add an invitation for another parent to join in on user's data.
func (user *User) InviteParent(env *config.Env, inviteEmail string) error {
	session, err := env.DB.GetConnection()
	if err != nil {
		return err
	}

	res, err := gorethink.Table("invites").Filter(map[string]interface{}{
		"inviteEmail": inviteEmail,
	}).Run(session)
	if err != nil {
		return err
	}
	defer res.Close()

	//if there is already an invite for that user, return error.
	if !res.IsNil() {
		return errors.New(ErrExistingInvitation)
	}

	inviteUser := UserInvitation{
		UserID:      user.ID,
		InviteEmail: inviteEmail,
		Timestamp:   time.Now(),
	}
	_, err = gorethink.Table("invites").Insert(inviteUser).Run(session)
	if err != nil {
		return err
	}

	return nil
}

//GetInvites - return the current invites a user has sent out.
func (user *User) GetInvites(env *config.Env) ([]UserInvitation, error) {
	session, err := env.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	res, err := gorethink.Table("invites").
		Filter(map[string]interface{}{
			"userID": user.ID,
		}).
		OrderBy(gorethink.Desc("timestamp")).
		Run(session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var rows []UserInvitation
	err = res.All(&rows)
	if err != nil {
		return rows, err
	}

	return rows, nil
}

func (user *User) DeleteInvite(env *config.Env, id string) error {
	session, err := env.DB.GetConnection()
	if err != nil {
		return err
	}

	res, err := gorethink.Table("invites").
		Filter(map[string]interface{}{
			"userID": user.ID,
			"id":     id,
		}).
		Delete().
		Run(session)
	if err != nil {
		return err
	}
	defer res.Close()

	var answer gorethink.WriteResponse
	err = res.One(&answer)
	if err != nil {
		return err
	}

	if answer.Deleted > 0 {
		return nil
	}

	return errors.New("no record to delete")
}

//GetFamily - return the family for a user. used for lookups
func (user *User) GetFamily(env *config.Env) (Family, error) {
	session, err := env.DB.GetConnection()
	if err != nil {
		return Family{}, err
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
		return Family{}, err
	}
	defer res.Close()

	var family Family
	err = res.One(&family)
	if err != nil {
		return Family{}, err
	}
	return family, nil
}
