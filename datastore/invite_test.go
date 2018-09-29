package datastore_test

import (
	"testing"
	"time"

	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/datastore"
	"github.com/stretchr/testify/assert"
	"google.golang.org/appengine/aetest"
)

func TestUserInvite(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	defer done()
	if err != nil {
		t.Error("appengine datastore context error", err)
	}
	inviteService := datastore.UserInviteService{}

	//some setup
	userService := datastore.UserService{}
	user := &goparent.User{
		Name:     "Test User",
		Email:    "test@test.com",
		Username: "test@test.com",
		Password: "testing",
	}
	err = userService.Save(ctx, user)
	assert.Nil(t, err)
	assert.NotNil(t, user.CurrentFamily)

	userInvites, err := inviteService.SentInvites(ctx, user)
	assert.Nil(t, err)
	assert.Len(t, userInvites, 0)

	err = inviteService.InviteParent(ctx, user, "testInvite@test.com", time.Now())
	assert.Nil(t, err)

	userInvites, err = inviteService.SentInvites(ctx, user)
	assert.Nil(t, err)
	assert.Len(t, userInvites, 1)

	lookupInvite, err := inviteService.Invite(ctx, userInvites[0].ID)
	assert.Nil(t, err)
	assert.EqualValues(t, userInvites[0], lookupInvite)
}
