package rethinkdb

import (
	"errors"
	"testing"
	"time"

	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/config"
	"github.com/stretchr/testify/assert"
	r "gopkg.in/gorethink/gorethink.v3"
)

func TestInviteParent(t *testing.T) {
	timestamp := time.Now()
	testCases := []struct {
		desc        string
		env         *config.Env
		query1      *r.MockQuery
		query2      *r.MockQuery
		inviteEmail string
		user        *goparent.User
		returnError error
	}{
		{
			desc: "invite parent",
			env:  &config.Env{},
			query1: (&r.Mock{}).On(
				r.Table("invites").Filter(
					map[string]interface{}{
						"inviteEmail": "invitedUser@test.com",
					}),
			).Return(nil, nil),
			query2: (&r.Mock{}).
				On(
					r.Table("invites").Insert(
						map[string]interface{}{
							"userID":      "1",
							"inviteEmail": "invitedUser@test.com",
							"timestamp":   timestamp,
						}),
				).Return(nil, nil),
			inviteEmail: "invitedUser@test.com",
			user:        &goparent.User{ID: "1"},
		},
		{
			desc: "invite parent existing invite",
			env:  &config.Env{},
			query1: (&r.Mock{}).On(
				r.Table("invites").Filter(
					map[string]interface{}{
						"inviteEmail": "invitedUser@test.com",
					}),
			).Return(map[string]interface{}{
				"id":          "1",
				"userID":      "1",
				"inviteEmail": "invitedUser@test.com",
				"timestamp":   timestamp,
			}, nil),
			query2:      nil,
			inviteEmail: "invitedUser@test.com",
			user:        &goparent.User{ID: "1"},
			returnError: errors.New(ErrExistingInvitation),
		},
		{
			desc: "invite parent check error",
			env:  &config.Env{},
			query1: (&r.Mock{}).On(
				r.Table("invites").Filter(
					map[string]interface{}{
						"inviteEmail": "invitedUser@test.com",
					}),
			).Return(nil, errors.New("test error")),
			query2:      nil,
			inviteEmail: "invitedUser@test.com",
			user:        &goparent.User{ID: "1"},
			returnError: errors.New("test error"),
		},
		{
			desc: "invite parent insert error",
			env:  &config.Env{},
			query1: (&r.Mock{}).On(
				r.Table("invites").Filter(
					map[string]interface{}{
						"inviteEmail": "invitedUser@test.com",
					}),
			).Return(nil, nil),
			query2: (&r.Mock{}).
				On(
					r.Table("invites").Insert(
						map[string]interface{}{
							"userID":      "1",
							"inviteEmail": "invitedUser@test.com",
							"timestamp":   timestamp,
						}),
				).Return(nil, errors.New("test error")),
			inviteEmail: "invitedUser@test.com",
			user:        &goparent.User{ID: "1"},
			returnError: errors.New("test error"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query1)
			if tC.query2 != nil {
				mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query2)
			}
			tC.env.DB = config.DBEnv{Session: mock}
			uis := UserInviteService{Env: tC.env}
			err := uis.InviteParent(tC.user, tC.inviteEmail, timestamp)
			if tC.returnError != nil {
				assert.Error(t, tC.returnError, err)
			} else {
				assert.Nil(t, err)
			}
			mock.AssertExpectations(t)
		})
	}
}
