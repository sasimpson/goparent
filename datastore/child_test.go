package datastore_test

import (
	"testing"
	"time"

	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/datastore"
	"github.com/stretchr/testify/assert"
	"google.golang.org/appengine/aetest"
)

func TestDatastoreChild(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		t.Error("error", err)
	}
	//setup
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

	//lookup child that doesn't exist
	nilChild, err := childService.Child(ctx, "123")
	assert.Nil(t, nilChild)
	assert.NotNil(t, err)
	//create new child
	child := goparent.Child{
		Name:     "Test User Jr.",
		ParentID: user.ID,
		FamilyID: family.ID,
		Birthday: time.Now().AddDate(-1, 0, 0),
	}
	//no id before save
	assert.Empty(t, child.ID)
	//save should get an id
	err = childService.Save(ctx, &child)
	assert.Nil(t, err)
	assert.NotEmpty(t, child.ID)

	//second save should catch the bit of code that saves not creates
	err = childService.Save(ctx, &child)
	assert.Nil(t, err)
	assert.NotEqual(t, child.CreatedAt, child.LastUpdated)

	//get child by the id, should be the same as the one we created
	lookupChild, err := childService.Child(ctx, child.ID)
	assert.Nil(t, err)
	assert.ObjectsAreEqualValues(child, lookupChild)

	//delete child.  count is silly because deletes are true even if already gone
	deletedCount, err := childService.Delete(ctx, &child)
	assert.Nil(t, err)
	assert.Equal(t, 1, deletedCount)

}
