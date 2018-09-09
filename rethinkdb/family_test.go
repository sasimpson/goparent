package rethinkdb

import (
	"errors"
	"testing"
	"time"

	"github.com/sasimpson/goparent"
	"github.com/stretchr/testify/assert"

	r "gopkg.in/gorethink/gorethink.v3"
)

func TestFamilySave(t *testing.T) {
	testCases := []struct {
		desc        string
		family      *goparent.Family
		query       *r.MockQuery
		env         *goparent.Env
		returnError error
	}{
		{
			desc: "Save 1",
			family: &goparent.Family{
				Admin:       "1",
				Members:     []string{"1"},
				CreatedAt:   time.Now(),
				LastUpdated: time.Now(),
			},
			//one issue with the rethinkdb mocking is that you cannot mock out
			query: (&r.Mock{}).On(r.Table("family").MockAnything()).Once().Return(
				r.WriteResponse{
					Inserted:      1,
					Errors:        0,
					GeneratedKeys: []string{"1"},
				}, nil),
			env:         &goparent.Env{},
			returnError: nil,
		},
		{
			desc: "Save 2",
			family: &goparent.Family{
				ID:          "1",
				Admin:       "1",
				Members:     []string{"1"},
				CreatedAt:   time.Now(),
				LastUpdated: time.Now(),
			},
			query: (&r.Mock{}).On(r.Table("family").MockAnything()).Once().Return(
				r.WriteResponse{
					Replaced:      1,
					Updated:       0,
					Inserted:      0,
					Errors:        0,
					GeneratedKeys: []string{"1"}}, nil),
			env:         &goparent.Env{},
			returnError: nil,
		},
		{
			desc: "Save error",
			family: &goparent.Family{
				ID:          "1",
				Admin:       "1",
				Members:     []string{"1"},
				CreatedAt:   time.Now(),
				LastUpdated: time.Now(),
			},
			query: (&r.Mock{}).On(r.Table("family").MockAnything()).Once().Return(
				r.WriteResponse{
					Replaced:      0,
					Updated:       0,
					Inserted:      0,
					Errors:        1,
					GeneratedKeys: []string{"1"}}, errors.New("test error")),
			env:         &goparent.Env{},
			returnError: errors.New("test error"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)
			fs := FamilyService{Env: tC.env, DB: &DBEnv{Session: mock}}
			err := fs.Save(tC.family)
			if tC.returnError != nil {
				assert.Error(t, tC.returnError, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestFamily(t *testing.T) {
	timestamp := time.Now()
	testCases := []struct {
		desc        string
		env         *goparent.Env
		id          string
		family      *goparent.Family
		query       *r.MockQuery
		returnError error
	}{
		{
			desc: "return family",
			env:  &goparent.Env{},
			id:   "family-1",
			family: &goparent.Family{
				ID:          "family-1",
				Admin:       "1",
				Members:     []string{"1"},
				CreatedAt:   timestamp,
				LastUpdated: timestamp,
			},
			query: (&r.Mock{}).On(
				r.Table("family").Get("family-1"),
			).Return(
				map[string]interface{}{
					"id":           "family-1",
					"admin":        "1",
					"members":      []string{"1"},
					"created_at":   timestamp,
					"last_updated": timestamp,
				}, nil),
		},
		{
			desc: "return error",
			env:  &goparent.Env{},
			id:   "family-1",
			family: &goparent.Family{
				ID:          "family-1",
				Admin:       "1",
				Members:     []string{"1"},
				CreatedAt:   timestamp,
				LastUpdated: timestamp,
			},
			query: (&r.Mock{}).On(
				r.Table("family").Get("family-1"),
			).Return(
				nil, r.ErrEmptyResult),
			returnError: r.ErrEmptyResult,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)
			fs := FamilyService{Env: tC.env, DB: &DBEnv{Session: mock}}
			family, err := fs.Family(tC.id)
			if tC.returnError != nil {
				assert.Error(t, tC.returnError, err)
			} else {
				assert.Nil(t, err)
				assert.EqualValues(t, tC.family, family)
			}
		})
	}
}

func TestChildren(t *testing.T) {
	timestamp := time.Now()
	testCases := []struct {
		desc         string
		env          *goparent.Env
		family       *goparent.Family
		query        *r.MockQuery
		resultLength int
		returnError  error
	}{
		{
			desc: "get some children",
			env:  &goparent.Env{},
			family: &goparent.Family{
				ID:          "family-1",
				Admin:       "1",
				Members:     []string{"1"},
				CreatedAt:   timestamp,
				LastUpdated: timestamp,
			},
			query: (&r.Mock{}).On(
				r.Table("children").Filter(
					map[string]interface{}{
						"familyID": "family-1",
					}).OrderBy(r.Desc("birthday")),
			).Return([]map[string]interface{}{
				{
					"id":       "child-1",
					"name":     "test-child-1",
					"userID":   "user-1",
					"familyID": "family-1",
					"birthday": timestamp.AddDate(0, -30, 0),
				},
				{
					"id":       "child-2",
					"name":     "test-child-2",
					"userID":   "user-1",
					"familyID": "family-1",
					"birthday": timestamp.AddDate(-2, -30, 0),
				},
			}, nil),
			resultLength: 2,
		},
		{
			desc: "get error",
			env:  &goparent.Env{},
			family: &goparent.Family{
				ID:          "family-1",
				Admin:       "1",
				Members:     []string{"1"},
				CreatedAt:   timestamp,
				LastUpdated: timestamp,
			},
			query: (&r.Mock{}).On(
				r.Table("children").Filter(
					map[string]interface{}{
						"familyID": "family-1",
					}).OrderBy(r.Desc("birthday")),
			).Return(nil, errors.New("test error")),
			returnError: errors.New("test error"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)
			fs := FamilyService{Env: tC.env, DB: &DBEnv{Session: mock}}
			children, err := fs.Children(tC.family)
			if tC.returnError != nil {
				assert.Error(t, tC.returnError, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tC.resultLength, len(children))
			}
		})
	}
}

func TestAddMember(t *testing.T) {
	timestamp := time.Now()
	testCases := []struct {
		desc        string
		env         *goparent.Env
		family      *goparent.Family
		user        *goparent.User
		query       *r.MockQuery
		returnError error
	}{
		{
			desc: "add member",
			env:  &goparent.Env{},
			family: &goparent.Family{
				ID:          "family-1",
				Admin:       "1",
				Members:     []string{"1"},
				CreatedAt:   timestamp,
				LastUpdated: timestamp,
			},
			user: &goparent.User{
				ID: "user-1",
			},
			query: (&r.Mock{}).On(r.Table("family").MockAnything()).Once().Return(
				r.WriteResponse{
					Updated: 1,
					Errors:  0,
				}, nil),
		},
		{
			desc: "add member fail",
			env:  &goparent.Env{},
			family: &goparent.Family{
				ID:          "family-1",
				Admin:       "1",
				Members:     []string{"user-1"},
				CreatedAt:   timestamp,
				LastUpdated: timestamp,
			},
			user: &goparent.User{
				ID: "user-1",
			},
			query: (&r.Mock{}).On(r.Table("family").MockAnything()).Once().Return(
				r.WriteResponse{
					Updated: 1,
					Errors:  0,
				}, nil),
			returnError: errors.New("user already in that family"),
		},
		{
			desc: "add member db error",
			env:  &goparent.Env{},
			family: &goparent.Family{
				ID:          "family-1",
				Admin:       "1",
				Members:     []string{"user-1"},
				CreatedAt:   timestamp,
				LastUpdated: timestamp,
			},
			user: &goparent.User{
				ID: "user-1",
			},
			query: (&r.Mock{}).On(r.Table("family").MockAnything()).Once().Return(
				r.WriteResponse{
					Updated: 0,
					Errors:  1,
				}, errors.New("test error")),
			returnError: errors.New("test error"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)
			fs := FamilyService{Env: tC.env, DB: &DBEnv{Session: mock}}
			err := fs.AddMember(tC.family, tC.user)
			if tC.returnError != nil {
				assert.Error(t, tC.returnError, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestGetAdminFamily(t *testing.T) {
	timestamp := time.Now()
	testCases := []struct {
		desc        string
		env         *goparent.Env
		user        *goparent.User
		family      *goparent.Family
		query       *r.MockQuery
		returnError error
	}{
		{
			desc: "valid call",
			env:  &goparent.Env{},
			user: &goparent.User{
				ID: "user-1",
			},
			family: &goparent.Family{
				ID:          "family-1",
				Admin:       "user-1",
				Members:     []string{"user-1"},
				CreatedAt:   timestamp,
				LastUpdated: timestamp,
			},
			query: (&r.Mock{}).On(r.Table("family").Filter(map[string]interface{}{
				"admin": "user-1",
			})).Return(map[string]interface{}{
				"id":           "family-1",
				"admin":        "user-1",
				"members":      []string{"user-1"},
				"created_at":   timestamp,
				"last_updated": timestamp,
			}, nil),
		},
		{
			desc: "nothing returned",
			env:  &goparent.Env{},
			user: &goparent.User{
				ID: "user-1",
			},
			family: &goparent.Family{
				ID:          "family-1",
				Admin:       "user-1",
				Members:     []string{"user-1"},
				CreatedAt:   timestamp,
				LastUpdated: timestamp,
			},
			query: (&r.Mock{}).On(r.Table("family").Filter(map[string]interface{}{
				"admin": "user-1",
			})).Return(nil, r.ErrEmptyResult),
			returnError: r.ErrEmptyResult,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)
			fs := FamilyService{Env: tC.env, DB: &DBEnv{Session: mock}}
			family, err := fs.GetAdminFamily(tC.user)
			t.Logf("%#v %#v", family, err)
			if tC.returnError != nil {
				assert.Error(t, tC.returnError, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tC.family.ID, family.ID)
			}
		})
	}
}
