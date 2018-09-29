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

	//test no invites sent
	userInvites, err := inviteService.SentInvites(ctx, user)
	assert.Nil(t, err)
	assert.Len(t, userInvites, 0)

	//invite someone
	err = inviteService.InviteParent(ctx, user, "mrstest@test.com", time.Now())
	assert.Nil(t, err)

	//test sent invite
	userInvites, err = inviteService.SentInvites(ctx, user)
	assert.Nil(t, err)
	assert.Len(t, userInvites, 1)

	//find that invite by id
	lookupInvite, err := inviteService.Invite(ctx, userInvites[0].ID)
	assert.Nil(t, err)
	assert.EqualValues(t, userInvites[0], lookupInvite)

	//create a new user that matches our invite from before
	mrsUser := &goparent.User{
		Name:     "Mrs Test User",
		Email:    "mrstest@test.com",
		Username: "mrstest@test.com",
		Password: "testing",
	}
	err = userService.Save(ctx, mrsUser)
	assert.Nil(t, err)

	//test getting existing invites for a user
	invites, err := inviteService.Invites(ctx, mrsUser)
	assert.Nil(t, err)
	assert.Len(t, invites, 1)

	//test accepting that invite
	err = inviteService.Accept(ctx, mrsUser, invites[0].ID)

	//test invite was deleted
	// deletedInvite, err := inviteService.Invite(ctx, invites[0].ID)
	// assert.Nil(t, err)
	// assert.Nil(t, deletedInvite)
}
