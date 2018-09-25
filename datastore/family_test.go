package datastore_test

import (
	"testing"

	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/datastore"
	"github.com/stretchr/testify/assert"
	"google.golang.org/appengine/aetest"
)

func TestDatastoreFamily(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		t.Error("error", err)
	}

	familyService := datastore.FamilyService{}

	nilFamily, err := familyService.Family(ctx, "123")
	assert.Nil(t, nilFamily)
	assert.NotNil(t, err)

	//test new family creation
	newFamily := &goparent.Family{}
	err = familyService.Save(ctx, newFamily)
	assert.Nil(t, err)
	//test family save w/id
	err = familyService.Save(ctx, newFamily)
	assert.Nil(t, err)

	//create a user to test adding to a family
	userService := datastore.UserService{}
	testUser := &goparent.User{
		Name:     "Test User",
		Email:    "test@test.com",
		Username: "test@test.com",
		Password: "testing",
	}
	//save it
	err = userService.Save(ctx, testUser)
	assert.Nil(t, err)
	assert.NotNil(t, testUser.ID)

	//add the member to the family
	err = familyService.AddMember(ctx, newFamily, testUser)
	assert.Nil(t, err)

	//get the admin family for the user, which should have been
	//created on save. should be equal to the returned admin family id
	// adminFamily, err := familyService.GetAdminFamily(ctx, testUser)
	// t.Logf("%#v", adminFamily)
	// t.Logf("%#v", err)
	// assert.Nil(t, err)
	// assert.Equal(t, testUser.CurrentFamily, adminFamily.ID)
}
