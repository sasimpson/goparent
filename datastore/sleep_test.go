package datastore_test

import (
	"testing"
	"time"

	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/datastore"
	"github.com/stretchr/testify/assert"
	"google.golang.org/appengine/aetest"
)

func TestDatastoreSleep(t *testing.T) {
	// t.Skip()
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

	//tests!
	timestamp := time.Now()

	sleep := &goparent.Sleep{
		ChildID:  child.ID,
		FamilyID: family.ID,
		UserID:   user.ID,
		Start:    timestamp.Add(time.Minute * -5),
		End:      timestamp,
	}
	assert.Empty(t, sleep.ID)
	assert.Empty(t, sleep.CreatedAt)
	assert.Empty(t, sleep.LastUpdated)

	sleepService := datastore.SleepService{}
	err = sleepService.Save(ctx, sleep)
	assert.Nil(t, err)
	assert.NotEmpty(t, sleep.ID)
	assert.NotEmpty(t, sleep.CreatedAt)
	assert.NotEmpty(t, sleep.LastUpdated)

	err = sleepService.Save(ctx, sleep)
	assert.Nil(t, err)
	assert.NotEqual(t, sleep.CreatedAt, sleep.LastUpdated)

	// summary, err := sleepService.Stats(ctx, child)
	// log.Printf("%#v", summary)
	// assert.Nil(err)
	// assert.Len(summary.Data, 1)
	// assert.Equal(1, summary.Range)
	// assert.Equal(300, summary.Total)

	// //status should be false right here because we haven't started a sleep
	// sleep, status, err := sleepService.Status(ctx, family, child)
	// assert.Nil(t, err)
	// assert.Nil(t, sleep)
	// assert.False(t, status)

	// //started sleep should now exist
	// err = sleepService.Start(ctx, family, child)
	// assert.Nil(t, err)

	// //verify with status check
	// sleep, status, err = sleepService.Status(ctx, family, child)
	// assert.Nil(t, err)
	// assert.NotNil(t, sleep)
	// assert.True(t, status)

	// //starting a new sleep should return an error since there is already one started
	// err = sleepService.Start(ctx, family, child)
	// assert.NotNil(t, err)
	// assert.EqualError(t, err, goparent.ErrExistingStart.Error())

	// //end should end the last sleep
	// err = sleepService.End(ctx, family, child)
	// assert.Nil(t, err)

	// //status should be back to false
	// sleep, status, err = sleepService.Status(ctx, family, child)
	// assert.Nil(t, sleep)
	// assert.Nil(t, err)
	// assert.False(t, status)

	// //now the error should say we cannot end something that doesn't exist.
	// err = sleepService.End(ctx, family, child)
	// assert.EqualError(t, err, goparent.ErrNoExistingSession.Error())

}
