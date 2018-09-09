package rethinkdb

import (
	"errors"
	"testing"
	"time"

	"github.com/sasimpson/goparent"
	"github.com/stretchr/testify/assert"
	r "gopkg.in/gorethink/gorethink.v3"
)

func TestSave(t *testing.T) {
	timestamp := time.Now()
	testCases := []struct {
		desc        string
		env         *goparent.Env
		child       *goparent.Child
		query       *r.MockQuery
		returnError error
	}{
		{
			desc:  "save child",
			env:   &goparent.Env{},
			child: &goparent.Child{Name: "test child", ParentID: "1", FamilyID: "1", Birthday: timestamp.AddDate(-1, 0, 0)},
			query: (&r.Mock{}).On(r.Table("children").MockAnything()).Once().Return(
				r.WriteResponse{
					Replaced:      0,
					Updated:       0,
					Inserted:      1,
					Errors:        0,
					GeneratedKeys: []string{"1"}}, nil),
		},
		{
			desc:  "error on save child",
			env:   &goparent.Env{},
			child: &goparent.Child{Name: "test child", ParentID: "1", FamilyID: "1", Birthday: timestamp.AddDate(-1, 0, 0)},
			query: (&r.Mock{}).On(r.Table("children").MockAnything()).Once().Return(
				r.WriteResponse{Errors: 1}, errors.New("test error")),
			returnError: errors.New("test error"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)
			cs := ChildService{Env: tC.env, DB: &DBEnv{Session: mock}}
			err := cs.Save(tC.child)
			if tC.returnError != nil {
				assert.Equal(t, tC.returnError, err)
			} else {
				assert.Nil(t, nil)
			}
		})
	}
}

func TestChild(t *testing.T) {
	timestamp := time.Now()
	testCases := []struct {
		desc        string
		env         *goparent.Env
		id          string
		child       *goparent.Child
		query       *r.MockQuery
		returnError error
	}{
		{
			desc: "get child",
			env:  &goparent.Env{},
			id:   "child-1",
			child: &goparent.Child{
				ID:       "child-1",
				Name:     "test child",
				ParentID: "1",
				FamilyID: "1",
				Birthday: timestamp.AddDate(-1, 0, 0)},
			query: (&r.Mock{}).On(
				r.Table("children").Get("child-1"),
			).Return(
				map[string]interface{}{
					"id":       "child-1",
					"name":     "test child",
					"parentID": "1",
					"familyID": "1",
					"birthday": timestamp.AddDate(-1, 0, 0),
				}, nil,
			),
		},
		{
			desc: "get child error",
			env:  &goparent.Env{},
			id:   "child-1",
			child: &goparent.Child{
				ID:       "child-1",
				Name:     "test child",
				ParentID: "1",
				FamilyID: "1",
				Birthday: timestamp.AddDate(-1, 0, 0)},
			query: (&r.Mock{}).On(
				r.Table("children").Get("child-1"),
			).Return(
				nil, errors.New("test error"),
			),
			returnError: errors.New("test error"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)
			cs := ChildService{Env: tC.env, DB: &DBEnv{Session: mock}}
			child, err := cs.Child(tC.id)
			if tC.returnError != nil {
				assert.Equal(t, tC.returnError, err)
			} else {
				assert.Nil(t, nil)
				assert.EqualValues(t, tC.child, child)
			}
		})
	}
}

func TestDeleteChild(t *testing.T) {
	timestamp := time.Now()
	testCases := []struct {
		desc        string
		env         *goparent.Env
		id          string
		child       *goparent.Child
		query       *r.MockQuery
		returnError error
		result      int
	}{
		{
			desc: "delete child",
			env:  &goparent.Env{},
			id:   "child-1",
			child: &goparent.Child{
				ID:       "child-1",
				Name:     "test child",
				ParentID: "1",
				FamilyID: "1",
				Birthday: timestamp.AddDate(-1, 0, 0)},
			query: (&r.Mock{}).On(
				r.Table("children").Get("child-1").Delete(),
			).Return(r.WriteResponse{Deleted: 1}, nil),
			result: 1,
		},
		{
			desc: "delete child error",
			env:  &goparent.Env{},
			id:   "child-1",
			child: &goparent.Child{
				ID:       "child-1",
				Name:     "test child",
				ParentID: "1",
				FamilyID: "1",
				Birthday: timestamp.AddDate(-1, 0, 0)},
			query: (&r.Mock{}).On(
				r.Table("children").Get("child-1").Delete(),
			).Return(r.WriteResponse{Errors: 1}, errors.New("test error")),
			returnError: errors.New("test error"),
			result:      0,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)
			cs := ChildService{Env: tC.env, DB: &DBEnv{Session: mock}}
			num, err := cs.Delete(tC.child)
			if tC.returnError != nil {
				assert.Equal(t, tC.returnError, err)
			} else {
				assert.Nil(t, nil)
			}
			assert.Equal(t, tC.result, num)

		})
	}
}
