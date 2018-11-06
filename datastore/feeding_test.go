package datastore_test

import (
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/datastore"
	"github.com/stretchr/testify/assert"
	"google.golang.org/appengine/aetest"
)

func TestDatastoreFeeding(t *testing.T) {
	assert := assert.New(t)
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
		Name:     "Test User Feeding",
		Email:    "testf@test.com",
		Username: "testf@test.com",
		Password: "testing",
	}
	err = userService.Save(ctx, user)
	if err != nil {
		t.Fatalf("error: %#v", err)
	}
	log.Printf("%s", reflect.TypeOf(err))
	assert.Nil(err)
	assert.NotNil(user.CurrentFamily)

	family, err := familyService.Family(ctx, user.CurrentFamily)
	assert.Nil(err)
	assert.NotNil(family)

	child := &goparent.Child{
		Name:     "Test User Jr",
		ParentID: user.ID,
		FamilyID: family.ID,
		Birthday: time.Date(2014, time.October, 1, 0, 0, 0, 0, time.UTC),
	}
	err = childService.Save(ctx, child)
	assert.Nil(err)
	assert.NotNil(child.ID)

	feeding := &goparent.Feeding{
		ChildID:   child.ID,
		FamilyID:  family.ID,
		UserID:    user.ID,
		Type:      "bottle",
		Amount:    10.0,
		TimeStamp: time.Date(2011, time.October, 02, 0, 0, 0, 0, time.UTC),
	}
	assert.Empty(feeding.ID)
	assert.Empty(feeding.CreatedAt)
	assert.Empty(feeding.LastUpdated)

	feedingService := datastore.FeedingService{}
	err = feedingService.Save(ctx, feeding)
	assert.Nil(err)
	assert.NotEmpty(feeding.ID)
	assert.NotEmpty(feeding.CreatedAt)
	assert.NotEmpty(feeding.LastUpdated)

	err = feedingService.Save(ctx, feeding)
	assert.Nil(err)
	assert.NotEqual(feeding.CreatedAt, feeding.LastUpdated)
}
