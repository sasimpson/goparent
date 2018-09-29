package datastore_test

import (
	"log"
	"reflect"
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
}

func TestDatastoreFamilyChildren(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		t.Error("appengine context error", err)
	}

	//some setup
	familyService := datastore.FamilyService{}
	userService := datastore.UserService{}
	childService := datastore.ChildService{}
	user := &goparent.User{
		Name:     "Test User",
		Email:    "test@test.com",
		Username: "test@test.com",
		Password: "testing",
	}
	err = userService.Save(ctx, user)
	assert.Nil(t, err)
	assert.NotNil(t, user.CurrentFamily)

	family, err := familyService.Family(ctx, user.CurrentFamily)
	assert.Nil(t, err)
	assert.NotNil(t, family)

	// //test no children yet
	noChildren, err := familyService.Children(ctx, family)
	assert.Nil(t, err)
	assert.Empty(t, noChildren)

	child := &goparent.Child{
		Name:     "Test User Jr",
		FamilyID: user.CurrentFamily,
		ParentID: user.ID,
	}
	err = childService.Save(ctx, child)
	assert.Nil(t, err)
	assert.NotEmpty(t, child.ID)

	children, err := familyService.Children(ctx, family)

	log.Println("typeOf: ", reflect.TypeOf(children))
	assert.Nil(t, err)
	assert.NotEmpty(t, children)
	assert.Len(t, children, 1)

}
