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

	childService := datastore.ChildService{}
	nilChild, err := childService.Child(ctx, "123")
	assert.Nil(t, nilChild)
	assert.NotNil(t, err)

	child := goparent.Child{
		Name:     "Test User Jr.",
		ParentID: "1",
		FamilyID: "1",
		Birthday: time.Now().AddDate(-1, 0, 0),
	}

	assert.Empty(t, child.ID)
	err = childService.Save(ctx, &child)
	assert.Nil(t, err)
	assert.NotEmpty(t, child.ID)
	err = childService.Save(ctx, &child)
	assert.Nil(t, err)
	assert.NotEqual(t, child.CreatedAt, child.LastUpdated)

	lookupChild, err := childService.Child(ctx, child.ID)
	assert.Nil(t, err)
	assert.ObjectsAreEqualValues(child, lookupChild)

	deletedCount, err := childService.Delete(ctx, &child)
	assert.Nil(t, err)
	assert.Equal(t, 1, deletedCount)

}
