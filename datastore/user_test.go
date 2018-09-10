package datastore_test

import (
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
	us := datastore.UserService{}

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
	t.Logf("%#v", user)
	assert.Nil(t, err)

	//get user
	nextUser, err := us.User(ctx, user.ID)
	assert.Nil(t, err)
	assert.NotNil(t, nextUser)

	//save user with no current family
	nextUser.CurrentFamily = ""
	err = us.Save(ctx, nextUser)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "datastore.UserService.Save.1b: FamilyService.GetAdminFamily: no family found with user as admin")

	//test login succeed
	loggedInUser, err := us.UserByLogin(ctx, "test@test.com", "testing")
	assert.NotNil(t, loggedInUser)
	assert.Nil(t, err)

	//test bad login
	notLoggedInUser, err := us.UserByLogin(ctx, "test@test.com", "test")
	assert.Nil(t, notLoggedInUser)
	assert.NotNil(t, err)
	assert.EqualError(t, err, "no result for that username password combo")

}
