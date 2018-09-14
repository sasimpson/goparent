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

	userService := datastore.UserService{}
	testUser := &goparent.User{}
	err = userService.Save(ctx, testUser)
	assert.Nil(t, err)
	err = familyService.AddMember(ctx, newFamily, testUser)
	assert.Nil(t, err)
	err = familyService.AddMember(ctx, newFamily, testUser)
	assert.NotNil(t, err)

}
