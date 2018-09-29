package rethinkdb

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/mock"
	"github.com/stretchr/testify/assert"
	r "gopkg.in/gorethink/gorethink.v3"
)

func TestInviteParent(t *testing.T) {
	timestamp := time.Now()
	testCases := []struct {
		desc        string
		env         *goparent.Env
		query1      *r.MockQuery
		query2      *r.MockQuery
		inviteEmail string
		user        *goparent.User
		returnError error
	}{
		{
			desc: "invite parent",
			env:  &goparent.Env{DB: &mock.DBEnv{}},
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
			env:  &goparent.Env{DB: &mock.DBEnv{}},
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
			env:  &goparent.Env{DB: &mock.DBEnv{}},
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
			env:  &goparent.Env{DB: &mock.DBEnv{}},
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
			ctx := context.Background()
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query1)
			if tC.query2 != nil {
				mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query2)
			}
			uis := UserInviteService{Env: tC.env, DB: &DBEnv{Session: mock}}
			err := uis.InviteParent(ctx, tC.user, tC.inviteEmail, timestamp)
			if tC.returnError != nil {
				assert.EqualError(t, tC.returnError, err.Error())
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
		env          *goparent.Env
		user         *goparent.User
		query        *r.MockQuery
		returnLength int
		returnError  error
	}{
		{
			desc: "nothing returned",
			env:  &goparent.Env{DB: &mock.DBEnv{}},
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
			env:  &goparent.Env{DB: &mock.DBEnv{}},
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
			env:  &goparent.Env{DB: &mock.DBEnv{}},
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
			ctx := context.Background()
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)

			uis := UserInviteService{Env: tC.env, DB: &DBEnv{Session: mock}}
			invites, err := uis.SentInvites(ctx, tC.user)
			if tC.returnError != nil {
				assert.EqualError(t, tC.returnError, err.Error())
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
		env         *goparent.Env
		query       *r.MockQuery
		id          string
		invite      *goparent.UserInvitation
		returnError error
	}{
		{
			desc: "nothing returned",
			env:  &goparent.Env{DB: &mock.DBEnv{}},
			query: (&r.Mock{}).On(
				r.Table("invites").Get("1"),
			).Return(nil, nil),
			id:          "1",
			returnError: r.ErrEmptyResult,
		},
		{
			desc: "invite returned",
			env:  &goparent.Env{DB: &mock.DBEnv{}},
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
			env:  &goparent.Env{DB: &mock.DBEnv{}},
			query: (&r.Mock{}).On(
				r.Table("invites").Get("1"),
			).Return(nil, errors.New("test error")),
			id:          "1",
			returnError: errors.New("test error"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctx := context.Background()
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)

			uis := UserInviteService{Env: tC.env, DB: &DBEnv{Session: mock}}
			invite, err := uis.Invite(ctx, tC.id)
			if tC.returnError != nil {
				assert.EqualError(t, tC.returnError, err.Error())
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
		env          *goparent.Env
		user         *goparent.User
		query        *r.MockQuery
		returnLength int
		returnError  error
	}{
		{
			desc: "nothing returned",
			env:  &goparent.Env{DB: &mock.DBEnv{}},
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
			env:  &goparent.Env{DB: &mock.DBEnv{}},
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
			env:  &goparent.Env{DB: &mock.DBEnv{}},
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
			ctx := context.Background()
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)

			uis := UserInviteService{Env: tC.env, DB: &DBEnv{Session: mock}}
			invites, err := uis.Invites(ctx, tC.user)
			t.Logf("invites: %#v len: %#v err: %#v", invites, len(invites), err)
			if tC.returnError != nil {
				assert.EqualError(t, tC.returnError, err.Error())
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tC.returnLength, len(invites))
		})
	}
}

func TestAcceptInvite(t *testing.T) {
	//TODO: add error conditions
	timestamp := time.Now()
	testCases := []struct {
		desc           string
		env            *goparent.Env
		id             string
		invitedUser    *goparent.User
		inviteQuery    *r.MockQuery
		returnError    error
		userQuery      *r.MockQuery
		familyQuery    *r.MockQuery
		addMemberQuery *r.MockQuery
		deleteQuery    *r.MockQuery
	}{
		{
			desc: "invite accepted, no errors",
			env:  &goparent.Env{DB: &mock.DBEnv{}},
			id:   "invite-1",
			invitedUser: &goparent.User{
				ID: "user-2", Email: "invitedUser@test.com",
			},
			inviteQuery: (&r.Mock{}).On(
				r.Table("invites").Filter(
					map[string]interface{}{
						"id":          "invite-1",
						"inviteEmail": "invitedUser@test.com",
					},
				).OrderBy(r.Desc("timestamp")),
			).Return(map[string]interface{}{
				"id":          "invite-1",
				"inviteEmail": "invitedUser@test.com",
				"userID":      "user-1",
				"timestamp":   time.Now(),
			}, nil),
			userQuery: (&r.Mock{}).On(r.Table("users").MockAnything()).Once().Return(
				map[string]interface{}{
					"id":            "user-1",
					"currentFamily": "family-1",
				}, nil),
			familyQuery: (&r.Mock{}).On(r.Table("family").MockAnything()).Once().Return(
				map[string]interface{}{
					"id":           "family-1",
					"created_at":   timestamp,
					"last_updated": timestamp,
				}, nil),
			addMemberQuery: (&r.Mock{}).On(r.Table("family").MockAnything()).Once().Return(
				r.WriteResponse{
					Updated: 1,
					Errors:  0,
				}, nil),
			deleteQuery: (&r.Mock{}).On(r.Table("invites").Filter(map[string]interface{}{
				"id": "invite-1",
			}).Delete()).Return(
				r.WriteResponse{
					Deleted: 1,
					Errors:  0,
				}, nil),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctx := context.Background()
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.inviteQuery)
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.userQuery)
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.familyQuery)
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.addMemberQuery)
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.deleteQuery)

			uis := UserInviteService{Env: tC.env, DB: &DBEnv{Session: mock}}
			err := uis.Accept(ctx, tC.invitedUser, tC.id)
			if tC.returnError != nil {
				assert.EqualError(t, tC.returnError, err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		desc        string
		env         *goparent.Env
		invite      *goparent.UserInvitation
		deleteQuery *r.MockQuery
		returnError error
	}{
		{
			desc: "valid delete",
			env:  &goparent.Env{DB: &mock.DBEnv{}},
			invite: &goparent.UserInvitation{
				ID: "invite-1",
			},
			deleteQuery: (&r.Mock{}).On(
				r.Table("invites").Filter(map[string]interface{}{
					"id": "invite-1",
				}).Delete(),
			).Return(r.WriteResponse{
				Errors:  0,
				Deleted: 1,
			}, nil),
			returnError: nil,
		},
		{
			desc: "delete error",
			env:  &goparent.Env{DB: &mock.DBEnv{}},
			invite: &goparent.UserInvitation{
				ID: "invite-1",
			},
			deleteQuery: (&r.Mock{}).On(
				r.Table("invites").Filter(map[string]interface{}{
					"id": "invite-1",
				}).Delete(),
			).Return(r.WriteResponse{
				Errors:  1,
				Deleted: 0,
			}, errors.New("test error")),
			returnError: errors.New("test error"),
		},
		{
			desc: "nothing deleted",
			env:  &goparent.Env{DB: &mock.DBEnv{}},
			invite: &goparent.UserInvitation{
				ID: "invite-1",
			},
			deleteQuery: (&r.Mock{}).On(
				r.Table("invites").Filter(map[string]interface{}{
					"id": "invite-1",
				}).Delete(),
			).Return(r.WriteResponse{
				Errors:  0,
				Deleted: 0,
			}, nil),
			returnError: errors.New("no record to delete"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctx := context.Background()
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.deleteQuery)

			uis := UserInviteService{Env: tC.env, DB: &DBEnv{Session: mock}}
			err := uis.Delete(ctx, tC.invite)
			if tC.returnError != nil {
				assert.EqualError(t, tC.returnError, err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
