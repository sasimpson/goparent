package rethinkdb

import (
	"context"
	"errors"
	"time"

	"github.com/sasimpson/goparent"
	gorethink "gopkg.in/gorethink/gorethink.v3"
)

//ErrExistingInvitation - the user already has an invitation
const ErrExistingInvitation string = "existing invitation for that user"

//UserInviteService - struct for implementing the interface
type UserInviteService struct {
	Env *goparent.Env
	DB  *DBEnv
}

//InviteParent - add an invitation for another parent to join in on user's data.
func (uis *UserInviteService) InviteParent(ctx context.Context, user *goparent.User, inviteEmail string, timestamp time.Time) error {
	err := uis.DB.GetConnection()
	if err != nil {
		return err
	}

	res, err := gorethink.Table("invites").Filter(map[string]interface{}{
		"inviteEmail": inviteEmail,
	}).Run(uis.DB.Session)
	if err != nil {
		return err
	}
	defer res.Close()

	//if there is already an invite for that user, return error.
	if !res.IsNil() {
		return errors.New(ErrExistingInvitation)
	}

	inviteUser := goparent.UserInvitation{
		UserID:      user.ID,
		InviteEmail: inviteEmail,
		Timestamp:   timestamp,
	}
	_, err = gorethink.Table("invites").Insert(inviteUser).Run(uis.DB.Session)
	if err != nil {
		return err
	}

	return nil
}

//SentInvites - return the current invites a user has sent out.
func (uis *UserInviteService) SentInvites(ctx context.Context, user *goparent.User) ([]*goparent.UserInvitation, error) {
	err := uis.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	res, err := gorethink.Table("invites").
		Filter(map[string]interface{}{
			"userID": user.ID,
		}).
		OrderBy(gorethink.Desc("timestamp")).
		Run(uis.DB.Session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var rows []*goparent.UserInvitation
	err = res.All(&rows)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

//Invite - return the invite by the id
func (uis *UserInviteService) Invite(ctx context.Context, id string) (*goparent.UserInvitation, error) {
	err := uis.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	res, err := gorethink.Table("invites").Get(id).Run(uis.DB.Session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var invite goparent.UserInvitation
	err = res.One(&invite)
	if err != nil {
		return nil, err
	}

	return &invite, nil
}

//Invites - return the invites that have been issued to a user based on the email.
func (uis *UserInviteService) Invites(ctx context.Context, user *goparent.User) ([]*goparent.UserInvitation, error) {
	err := uis.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	res, err := gorethink.Table("invites").
		Filter(map[string]interface{}{
			"inviteEmail": user.Email,
		}).
		OrderBy(gorethink.Desc("timestamp")).
		Run(uis.DB.Session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var rows []*goparent.UserInvitation
	err = res.All(&rows)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

//Accept - user can accept an invite, this will set their
// CurrentFamily and add them as a member to that family.
func (uis *UserInviteService) Accept(ctx context.Context, user *goparent.User, id string) error {
	err := uis.DB.GetConnection()
	if err != nil {
		return err
	}

	//get invite by id and invited user (current):
	res, err := gorethink.Table("invites").
		Filter(map[string]interface{}{
			"id":          id,
			"inviteEmail": user.Email,
		}).
		OrderBy(gorethink.Desc("timestamp")).
		Run(uis.DB.Session)
	if err != nil {
		return err
	}
	defer res.Close()

	var invite goparent.UserInvitation
	err = res.One(&invite)
	if err != nil {
		return err
	}

	//get the user and family that is doing the inviting
	us := UserService{Env: uis.Env, DB: uis.DB}
	invitingUser, err := us.User(ctx, invite.UserID)
	if err != nil {
		return err
	}

	//NOTE: this would need to be set in the invite if we allowed family switching
	family, err := us.GetFamily(ctx, invitingUser)
	if err != nil {
		return err
	}
	fs := FamilyService{Env: uis.Env, DB: uis.DB}
	//add the user to the family of the inviting user
	err = fs.AddMember(ctx, family, user)
	if err != nil {
		return err
	}

	//remove invite from system
	err = uis.Delete(ctx, &invite)
	if err != nil {
		return err
	}

	return nil
}

//Delete - a user can delete invites they have sent.
func (uis *UserInviteService) Delete(ctx context.Context, invite *goparent.UserInvitation) error {
	err := uis.DB.GetConnection()
	if err != nil {
		return err
	}

	res, err := gorethink.Table("invites").
		Filter(map[string]interface{}{
			"id": invite.ID,
		}).
		Delete().
		Run(uis.DB.Session)
	if err != nil {
		return err
	}
	defer res.Close()

	var answer gorethink.WriteResponse
	err = res.One(&answer)
	if err != nil {
		return err
	}

	if answer.Deleted > 0 {
		return nil
	}

	return errors.New("no record to delete")
}
