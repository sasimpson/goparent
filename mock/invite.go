package mock

import (
	"context"
	"time"

	"github.com/sasimpson/goparent"
)

//UserInvitationService -
type UserInvitationService struct {
	GetSentInvites  []*goparent.UserInvitation
	GetInvites      []*goparent.UserInvitation
	GetInvite       *goparent.UserInvitation
	InviteParentErr error
	SentInvitesErr  error
	InvitesErr      error
	InviteErr       error
	AcceptErr       error
	DeleteErr       error
}

//InviteParent -
func (m *UserInvitationService) InviteParent(context.Context, *goparent.User, string, time.Time) error {
	if m.InviteParentErr != nil {
		return m.InviteParentErr
	}
	return nil
}

//SentInvites -
func (m *UserInvitationService) SentInvites(context.Context, *goparent.User) ([]*goparent.UserInvitation, error) {
	if m.SentInvitesErr != nil {
		return nil, m.SentInvitesErr
	}
	return nil, nil
}

//Invite -
func (m *UserInvitationService) Invite(context.Context, string) (*goparent.UserInvitation, error) {
	if m.InviteErr != nil {
		return nil, m.InviteErr
	}
	if m.GetInvite != nil {
		return m.GetInvite, nil
	}
	return nil, nil
}

//Invites -
func (m *UserInvitationService) Invites(*goparent.User) ([]*goparent.UserInvitation, error) {
	if m.InvitesErr != nil {
		return nil, m.InvitesErr
	}

	if m.GetInvites != nil {
		return m.GetInvites, nil
	}

	return nil, nil
}

//Accept -
func (m *UserInvitationService) Accept(*goparent.User, string) error {
	if m.AcceptErr != nil {
		return m.AcceptErr
	}
	return nil
}

//Delete -
func (m *UserInvitationService) Delete(*goparent.UserInvitation) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	return nil
}
