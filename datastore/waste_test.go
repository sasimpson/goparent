package datastore_test

import (
	"testing"
	"time"

	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/datastore"
	"github.com/stretchr/testify/assert"
	"google.golang.org/appengine/aetest"
)

func TestDatastoreWaste(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		t.Error("error", err)
	}
	//setup:
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

	child := &goparent.Child{
		Name:     "Test User Jr",
		ParentID: user.ID,
		FamilyID: family.ID,
		Birthday: time.Date(2014, time.October, 1, 0, 0, 0, 0, time.UTC),
	}
	err = childService.Save(ctx, child)
	assert.Nil(t, err)
	assert.NotNil(t, child.ID)

	//actual waste tests:
	waste := &goparent.Waste{
		ChildID:   child.ID,
		FamilyID:  family.ID,
		UserID:    user.ID,
		Type:      1,
		TimeStamp: time.Date(2011, time.October, 02, 0, 0, 0, 0, time.UTC),
	}
	assert.Empty(t, waste.ID)
	assert.Empty(t, waste.CreatedAt)
	assert.Empty(t, waste.LastUpdated)

	wasteService := datastore.WasteService{}
	err = wasteService.Save(ctx, waste)
	assert.Nil(t, err)
	assert.NotEmpty(t, waste.ID)
	assert.NotEmpty(t, waste.CreatedAt)
	assert.NotEmpty(t, waste.LastUpdated)
	assert.Equal(t, waste.CreatedAt, waste.LastUpdated)

	err = wasteService.Save(ctx, waste)
	assert.Nil(t, err)
	assert.NotEqual(t, waste.CreatedAt, waste.LastUpdated)

	//test out of date range
	allWaste, err := wasteService.Waste(ctx, family, 7)
	assert.Nil(t, err)
	assert.Empty(t, allWaste)
	assert.Len(t, allWaste, 0)

	//add a new one in range
	wasteTS := time.Now()
	waste = &goparent.Waste{
		ChildID:   child.ID,
		FamilyID:  family.ID,
		UserID:    user.ID,
		Type:      1,
		TimeStamp: wasteTS,
	}
	assert.Empty(t, waste.ID)
	assert.Empty(t, waste.CreatedAt)
	assert.Empty(t, waste.LastUpdated)

	wasteService = datastore.WasteService{}
	err = wasteService.Save(ctx, waste)
	assert.Nil(t, err)
	assert.Equal(t, wasteTS, waste.TimeStamp)
	assert.NotEmpty(t, waste.ID)
	assert.NotEmpty(t, waste.CreatedAt)
	assert.NotEmpty(t, waste.LastUpdated)
	assert.Equal(t, waste.CreatedAt, waste.LastUpdated)

	//test in date range
	allWaste, err = wasteService.Waste(ctx, family, 7)
	assert.Nil(t, err)
	assert.NotEmpty(t, allWaste)
	assert.Len(t, allWaste, 1)
}
