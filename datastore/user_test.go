package datastore_test

import (
	"log"
	"testing"

	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/datastore"
	"github.com/stretchr/testify/assert"
	"google.golang.org/appengine/aetest"
)

func TestDatastoreUser(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		t.Error("error", err)
	}
	us := datastore.UserService{
		Env: &goparent.Env{
			Auth: goparent.Authentication{
				SigningKey: []byte("test"),
			},
		},
	}

	//test invalid user
	nilUser, err := us.User(ctx, "123")
	assert.Nil(t, nilUser)
	assert.NotNil(t, err)

	//test no user for login
	noUserForLogin, err := us.UserByLogin(ctx, "test@test.com", "testing")
	assert.Nil(t, noUserForLogin)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "no result for that username password combo")

	user := goparent.User{
		Name:     "test",
		Email:    "test@test.com",
		Password: "testing",
		Username: "test@test.com",
	}

	//save user
	err = us.Save(ctx, &user)
	assert.Nil(t, err)

	//check that the current family is the loaded family
	cFamily, err := us.GetFamily(ctx, &user)
	log.Println(err)
	log.Println(cFamily)
	assert.Nil(t, err)
	assert.NotNil(t, cFamily)
	assert.Equal(t, user.CurrentFamily, cFamily.ID)

	//get user
	nextUser, err := us.User(ctx, user.ID)
	assert.Nil(t, err)
	assert.NotNil(t, nextUser)

	//save user with no current family
	nextUser.CurrentFamily = ""
	err = us.Save(ctx, nextUser)
	assert.Nil(t, err)
	assert.NotEmpty(t, nextUser.CurrentFamily)

	families, err := us.GetAllFamily(ctx, nextUser)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(families))

	//test login succeed
	loggedInUser, err := us.UserByLogin(ctx, "test@test.com", "testing")
	assert.NotNil(t, loggedInUser)
	assert.Nil(t, err)

	//get token
	token, err := us.GetToken(loggedInUser)
	assert.NotNil(t, token)
	assert.Nil(t, err)

	//validate token
	validatedUser, ok, err := us.ValidateToken(ctx, token)
	assert.Nil(t, err)
	assert.True(t, ok)
	assert.EqualValues(t, loggedInUser, validatedUser)

	//test bad login
	notLoggedInUser, err := us.UserByLogin(ctx, "test@test.com", "test")
	assert.Nil(t, notLoggedInUser)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "no result for that username password combo")

}
