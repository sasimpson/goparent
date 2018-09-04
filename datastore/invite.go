package datastore

import (
	"time"

	"github.com/sasimpson/goparent"
)

//UserInviteService -
type UserInviteService struct {
}

//InviteParent -
func (s *UserInviteService) InviteParent(*goparent.User, string, time.Time) error {
	panic("not implemented")
}

//SentInvites -
func (s *UserInviteService) SentInvites(*goparent.User) ([]*goparent.UserInvitation, error) {
	panic("not implemented")
}

//Invite -
func (s *UserInviteService) Invite(string) (*goparent.UserInvitation, error) {
	panic("not implemented")
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
