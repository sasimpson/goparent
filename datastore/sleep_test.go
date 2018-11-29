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

	//tests!
	timestamp := time.Now()

	sleep := &goparent.Sleep{
		ChildID:  child.ID,
		FamilyID: family.ID,
		UserID:   user.ID,
		Start:    timestamp.Add(time.Minute * -5),
		End:      timestamp,
	}
	assert.Empty(sleep.ID)
	assert.Empty(sleep.CreatedAt)
	assert.Empty(sleep.LastUpdated)

	sleepService := datastore.SleepService{}
	err = sleepService.Save(ctx, sleep)
	assert.Nil(err)
	assert.NotEmpty(sleep.ID)
	assert.NotEmpty(sleep.CreatedAt)
	assert.NotEmpty(sleep.LastUpdated)

	err = sleepService.Save(ctx, sleep)
	assert.Nil(err)
	assert.NotEqual(sleep.CreatedAt, sleep.LastUpdated)

	// summary, err := sleepService.Stats(ctx, child)
	// log.Printf("%#v", summary)
	// assert.Nil(err)
	// assert.Len(summary.Data, 1)
	// assert.Equal(1, summary.Range)
	// assert.Equal(300, summary.Total)

	//status should be false right here because we haven't started a sleep
	sleep, status, err := sleepService.Status(ctx, family, child)
	assert.Nil(err)
	assert.Nil(sleep)
	assert.False(status)

	//started sleep should now exist
	err = sleepService.Start(ctx, family, child)
	assert.Nil(err)

	//verify with status check
	sleep, status, err = sleepService.Status(ctx, family, child)
	assert.Nil(err)
	assert.NotNil(sleep)
	assert.True(status)

	//starting a new sleep should return an error since there is already one started
	err = sleepService.Start(ctx, family, child)
	assert.NotNil(err)
	assert.EqualError(err, goparent.ErrExistingStart.Error())

	//end should end the last sleep
	err = sleepService.End(ctx, family, child)
	assert.Nil(err)

	//status should be back to false
	sleep, status, err = sleepService.Status(ctx, family, child)
	assert.Nil(sleep)
	assert.Nil(err)
	assert.False(status)

	//now the error should say we cannot end something that doesn't exist.
	err = sleepService.End(ctx, family, child)
	assert.EqualError(err, goparent.ErrNoExistingSession.Error())

}
