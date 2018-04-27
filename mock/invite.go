package mock

import (
	"time"

	"github.com/sasimpson/goparent"
)

type MockUserInvitationService struct {
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

func (m *MockUserInvitationService) InviteParent(*goparent.User, string, time.Time) error {
	if m.InviteParentErr != nil {
		return m.InviteParentErr
	}
	return nil
}

func (m *MockUserInvitationService) SentInvites(*goparent.User) ([]*goparent.UserInvitation, error) {
	if m.SentInvitesErr != nil {
		return nil, m.SentInvitesErr
	}
	return nil, nil
}

func (m *MockUserInvitationService) Invite(string) (*goparent.UserInvitation, error) {
	if m.InviteErr != nil {
		return nil, m.InviteErr
	}
	if m.GetInvite != nil {
		return m.GetInvite, nil
	}
	return nil, nil
}

func (m *MockUserInvitationService) Invites(*goparent.User) ([]*goparent.UserInvitation, error) {
	if m.InvitesErr != nil {
		return nil, m.InvitesErr
	}

	if m.GetInvites != nil {
		return m.GetInvites, nil
	}

	return nil, nil
}

func (m *MockUserInvitationService) Accept(*goparent.User, string) error {
	if m.AcceptErr != nil {
		return m.AcceptErr
	}
	return nil
}

func (m *MockUserInvitationService) Delete(*goparent.UserInvitation) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	return nil
}
