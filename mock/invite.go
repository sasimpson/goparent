package mock

import (
	"time"

	"github.com/sasimpson/goparent"
)

type MockUserInvitationService struct {
	InviteParentErr error
}

func (m *MockUserInvitationService) InviteParent(*goparent.User, string, time.Time) error {
	if m.InviteParentErr != nil {
		return m.InviteParentErr
	}
	return nil
}

func (m *MockUserInvitationService) SentInvites(*goparent.User) ([]*goparent.UserInvitation, error) {
	panic("not implemented")
}

func (m *MockUserInvitationService) Invite(string) (*goparent.UserInvitation, error) {
	panic("not implemented")
}

func (m *MockUserInvitationService) Invites(*goparent.User) ([]*goparent.UserInvitation, error) {
	panic("not implemented")
}

func (m *MockUserInvitationService) Accept(*goparent.User, string) error {
	panic("not implemented")
}

func (m *MockUserInvitationService) Delete(*goparent.UserInvitation) error {
	panic("not implemented")
}
