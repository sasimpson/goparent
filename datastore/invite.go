package datastore

import (
	"context"
	"time"

	"github.com/sasimpson/goparent"
	"google.golang.org/appengine/datastore"
)

//UserInviteService -
type UserInviteService struct {
	Env *goparent.Env
}

//InviteKind is the datastore kind representation
const InviteKind = "Invite"

//InviteParent adds an invitation for another parent to join in on a family's data
func (s *UserInviteService) InviteParent(ctx context.Context, user *goparent.User, inviteEmail string, timestamp time.Time) error {
	//get existing invites by the invitee's email and make sure they don't already have one.
	q := datastore.NewQuery(InviteKind).Filter("InviteEmail = ", inviteEmail).KeysOnly()
	keys, err := q.GetAll(ctx, nil)
	if len(keys) > 0 {
		return NewError("UserInviteService.InviteParent", err)
	}

	//if not, add an invite for the invitee
	userKey := datastore.NewKey(ctx, UserKind, user.ID, 0, nil)
	//set the user as the parent so we can lookup all invites sent by a user by ancestry.
	inviteKey := datastore.NewIncompleteKey(ctx, InviteKind, userKey)
	inviteUser := goparent.UserInvitation{
		UserID:      user.ID,
		InviteEmail: inviteEmail,
		Timestamp:   timestamp,
	}

	_, err = datastore.Put(ctx, inviteKey, &inviteUser)
	if err != nil {
		return err
	}

	return nil
}

//SentInvites -
func (s *UserInviteService) SentInvites(ctx context.Context, user *goparent.User) ([]*goparent.UserInvitation, error) {
	var userInvites []*goparent.UserInvitation

	userKey := datastore.NewKey(ctx, UserKind, user.ID, 0, nil)
	q := datastore.NewQuery(InviteKind).Ancestor(userKey)
	itx := q.Run(ctx)

	for {
		var invite goparent.UserInvitation
		_, err := itx.Next(&invite)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		userInvites = append(userInvites, &invite)
	}

	return userInvites, nil
}

//Invite -
func (s *UserInviteService) Invite(ctx context.Context, id string) (*goparent.UserInvitation, error) {
	q := datastore.NewQuery(InviteKind).Filter("ID =", id)
	itx := q.Run(ctx)

	var invite goparent.UserInvitation
	_, err := itx.Next(&invite)
	if err == datastore.Done {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &invite, nil

}

//Invites -
func (s *UserInviteService) Invites(*goparent.User) ([]*goparent.UserInvitation, error) {
	panic("not implemented")
}

//Accept -
func (s *UserInviteService) Accept(*goparent.User, string) error {
	panic("not implemented")
}

//Delete -
func (s *UserInviteService) Delete(*goparent.UserInvitation) error {
	panic("not implemented")
}
