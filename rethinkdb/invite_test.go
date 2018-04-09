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

func TestSentInvites(t *testing.T) {
	testCases := []struct {
		desc         string
		env          *config.Env
		user         *goparent.User
		query        *r.MockQuery
		returnLength int
		returnError  error
	}{
		{
			desc: "nothing returned",
			env:  &config.Env{},
			user: &goparent.User{ID: "1"},
			query: (&r.Mock{}).On(
				r.Table("invites").Filter(map[string]interface{}{
					"userID": "1",
				}).OrderBy(r.Desc("timestamp")),
			).Return(nil, nil),
			returnLength: 1,
		},
		{
			desc: "nothing returned",
			env:  &config.Env{},
			user: &goparent.User{ID: "1"},
			query: (&r.Mock{}).On(
				r.Table("invites").Filter(map[string]interface{}{
					"userID": "1",
				}).OrderBy(r.Desc("timestamp")),
			).Return(nil, errors.New("test error")),
			returnLength: 0,
			returnError:  errors.New("test error"),
		},
		{
			desc: "two invites returned",
			env:  &config.Env{},
			user: &goparent.User{ID: "1"},
			query: (&r.Mock{}).On(
				r.Table("invites").Filter(map[string]interface{}{
					"userID": "1",
				}).OrderBy(r.Desc("timestamp")),
			).Return([]map[string]interface{}{
				map[string]interface{}{
					"id":          "1",
					"inviteEmail": "invitedUser@test.com",
					"userID":      "1",
					"timestamp":   time.Now(),
				},
				map[string]interface{}{
					"id":          "2",
					"inviteEmail": "invitedUser2@test.com",
					"userID":      "1",
					"timestamp":   time.Now(),
				},
			}, nil),
			returnLength: 2,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)
			tC.env.DB = config.DBEnv{Session: mock}
			uis := UserInviteService{Env: tC.env}
			invites, err := uis.SentInvites(tC.user)
			if tC.returnError != nil {
				assert.Error(t, tC.returnError, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tC.returnLength, len(invites))
		})
	}
}

func TestInvite(t *testing.T) {
	timestamp := time.Now()
	testCases := []struct {
		desc        string
		env         *config.Env
		query       *r.MockQuery
		id          string
		invite      *goparent.UserInvitation
		returnError error
	}{
		{
			desc: "nothing returned",
			env:  &config.Env{},
			query: (&r.Mock{}).On(
				r.Table("invites").Get("1"),
			).Return(nil, nil),
			id:          "1",
			returnError: r.ErrEmptyResult,
		},
		{
			desc: "invite returned",
			env:  &config.Env{},
			query: (&r.Mock{}).On(
				r.Table("invites").Get("1"),
			).Return(map[string]interface{}{
				"id":          "1",
				"inviteEmail": "invitedUser@test.com",
				"userID":      "1",
				"timestamp":   timestamp,
			}, nil),
			id: "1",
			invite: &goparent.UserInvitation{
				ID:          "1",
				InviteEmail: "invitedUser@test.com",
				UserID:      "1",
				Timestamp:   timestamp,
			},
		},
		{
			desc: "error returned",
			env:  &config.Env{},
			query: (&r.Mock{}).On(
				r.Table("invites").Get("1"),
			).Return(nil, errors.New("test error")),
			id:          "1",
			returnError: errors.New("test error"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)
			tC.env.DB = config.DBEnv{Session: mock}
			uis := UserInviteService{Env: tC.env}
			invite, err := uis.Invite(tC.id)
			if tC.returnError != nil {
				assert.Error(t, tC.returnError, err)
			} else {
				assert.Nil(t, err)
				t.Logf("%#v", invite)
				assert.Equal(t, tC.id, invite.ID)
			}

		})
	}
}

func TestInvites(t *testing.T) {
	testCases := []struct {
		desc         string
		env          *config.Env
		user         *goparent.User
		query        *r.MockQuery
		returnLength int
		returnError  error
	}{
		{
			desc: "nothing returned",
			env:  &config.Env{},
			user: &goparent.User{ID: "1", Email: "testUser@test.com"},
			query: (&r.Mock{}).On(
				r.Table("invites").Filter(map[string]interface{}{
					"inviteEmail": "testUser@test.com",
				}).OrderBy(r.Desc("timestamp")),
			).Return(nil, nil),
			returnLength: 1,
		},
		{
			desc: "nothing returned",
			env:  &config.Env{},
			user: &goparent.User{ID: "1", Email: "testUser@test.com"},
			query: (&r.Mock{}).On(
				r.Table("invites").Filter(map[string]interface{}{
					"inviteEmail": "testUser@test.com",
				}).OrderBy(r.Desc("timestamp")),
			).Return(nil, errors.New("test error")),
			returnLength: 0,
			returnError:  errors.New("test error"),
		},
		{
			desc: "two invites returned",
			env:  &config.Env{},
			user: &goparent.User{ID: "1", Email: "testUser@test.com"},
			query: (&r.Mock{}).On(
				r.Table("invites").Filter(map[string]interface{}{
					"inviteEmail": "testUser@test.com",
				}).OrderBy(r.Desc("timestamp")),
			).Return([]map[string]interface{}{
				map[string]interface{}{
					"id":          "1",
					"inviteEmail": "testUser@test.com",
					"userID":      "1",
					"timestamp":   time.Now(),
				},
				map[string]interface{}{
					"id":          "2",
					"inviteEmail": "testUser@test.com",
					"userID":      "2",
					"timestamp":   time.Now(),
				},
			}, nil),
			returnLength: 2,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)
			tC.env.DB = config.DBEnv{Session: mock}
			uis := UserInviteService{Env: tC.env}
			invites, err := uis.Invites(tC.user)
			t.Logf("invites: %#v len: %#v err: %#v", invites, len(invites), err)
			if tC.returnError != nil {
				assert.Error(t, tC.returnError, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tC.returnLength, len(invites))
		})
	}
}
