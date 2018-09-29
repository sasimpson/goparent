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

//SentInvites gets the invites that a user has sent out to other people
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

//Invite gets an invite by its ID.
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

//Invites returns the invites that have been issued to a user based on their email
func (s *UserInviteService) Invites(ctx context.Context, user *goparent.User) ([]*goparent.UserInvitation, error) {
	var invites []*goparent.UserInvitation
	q := datastore.NewQuery(InviteKind).Filter("InviteEmail = ", user.Email)
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

		invites = append(invites, &invite)
	}
	return invites, nil
}

//Accept an invite from an existing user, add them to the family and set their current family.
func (s *UserInviteService) Accept(ctx context.Context, user *goparent.User, id string) error {
	//get invite by id and invited user
	invite, err := s.Invite(ctx, id)
	if err != nil {
		return err
	}

	//get the user and family that is doing the inviting
	us := UserService{Env: s.Env}
	invitingUser, err := us.User(ctx, invite.UserID)
	if err != nil {
		return err
	}

	family, err := us.GetFamily(ctx, invitingUser)
	if err != nil {
		return err
	}
	//add the user to the family of the inviting user
	fs := FamilyService{Env: s.Env}
	err = fs.AddMember(ctx, family, user)
	if err != nil {
		return err
	}

	//remove the invite from the system
	err = s.Delete(ctx, invite)
	if err != nil {
		return err
	}

	return nil
}

//Delete removes an invite
func (s *UserInviteService) Delete(ctx context.Context, invite *goparent.UserInvitation) error {
	userKey := datastore.NewKey(ctx, UserKind, invite.UserID, 0, nil)
	inviteKey := datastore.NewKey(ctx, InviteKind, invite.ID, 0, userKey)

	err := datastore.Delete(ctx, inviteKey)

	return err
}
