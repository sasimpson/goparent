package models

import (
	"errors"
	"testing"
	"time"

	"github.com/sasimpson/goparent/config"
	"github.com/stretchr/testify/assert"
	r "gopkg.in/gorethink/gorethink.v3"
)

func TestChildrenGetAll(t *testing.T) {
	var testEnv config.Env
	mock := r.NewMock()
	mock.On(
		r.Table("children").Filter(map[string]interface{}{"userid": "1"}).OrderBy(r.Desc("birthday")),
	).Return([]interface{}{
		map[string]interface{}{"id": "1", "name": "kiddo1", "parentID": 1, "birthday": time.Now()},
	}, nil)
	testEnv.DB.Session = mock

	var c Child
	children, err := c.GetAll(&testEnv, &User{ID: "1"})
	mock.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Len(t, children, 1)

	mock = r.NewMock()
	mock.On(
		r.Table("children").Filter(map[string]interface{}{"userid": "1"}).OrderBy(r.Desc("birthday")),
	).Return([]interface{}{}, nil)
	testEnv.DB.Session = mock
	children, err = c.GetAll(&testEnv, &User{ID: "1"})
	mock.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Len(t, children, 0)

	mock = r.NewMock()
	mock.On(
		r.Table("children").Filter(map[string]interface{}{"userid": "1"}).OrderBy(r.Desc("birthday")),
	).Return([]interface{}{}, errors.New("Test Error"))
	testEnv.DB.Session = mock
	children, err = c.GetAll(&testEnv, &User{ID: "1"})
	mock.AssertExpectations(t)
	assert.Error(t, err)
	assert.Len(t, children, 0)
}

func TestChildrenGetOne(t *testing.T) {
	var testEnv config.Env
	mock := r.NewMock()
	mock.On(
		r.Table("children").Filter(map[string]interface{}{"parentID": "1", "id": "1"}),
	).Return([]interface{}{
		map[string]interface{}{"id": "1", "name": "kiddo1", "parentID": 1, "birthday": time.Now()},
	}, nil)
	testEnv.DB.Session = mock

	var child Child
	err := child.GetChild(&testEnv, &User{ID: "1"}, "1")
	mock.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, "1", child.ID)
}

func TestChildrenSaveError(t *testing.T) {
	var testEnv config.Env
	timestamp := time.Now()
	mock := r.NewMock()
	mock.On(
		r.Table("children").Insert(
			map[string]interface{}{
				"parentID": "1",
				"birthday": timestamp,
				"name":     "joey",
			}, r.InsertOpts{Conflict: "replace"},
		),
	).Return(nil, errors.New("returned error"))
	testEnv.DB.Session = mock

	c := Child{Name: "joey", ParentID: "1", Birthday: timestamp}
	err := c.Save(&testEnv)
	mock.AssertExpectations(t)
	assert.Error(t, err)
	assert.EqualError(t, err, "returned error")
}

func TestChildrenSave(t *testing.T) {
	var testEnv config.Env

	testCases := []struct {
		desc     string
		id       string
		name     string
		parentid string
		birthday time.Time
	}{
		{
			desc:     "joey",
			id:       "1",
			name:     "Joe",
			parentid: "1",
			birthday: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:     "amy",
			id:       "2",
			name:     "Amy",
			parentid: "1",
			birthday: time.Date(2001, 2, 2, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := r.NewMock()
			mock.On(
				r.Table("children").Insert(
					map[string]interface{}{
						"name":     tC.name,
						"parentID": tC.parentid,
						"birthday": tC.birthday,
					}, r.InsertOpts{Conflict: "replace"},
				),
			).Return(
				r.WriteResponse{
					Inserted:      1,
					Errors:        0,
					GeneratedKeys: []string{tC.id},
				}, nil)
			testEnv.DB.Session = mock

			c := Child{
				Name:     tC.name,
				ParentID: tC.parentid,
				Birthday: tC.birthday,
			}

			err := c.Save(&testEnv)
			mock.AssertExpectations(t)
			assert.Nil(t, err)
			assert.Equal(t, tC.id, c.ID)
			assert.Equal(t, tC.parentid, c.ParentID)
			assert.Equal(t, tC.birthday, c.Birthday)
		})
	}
}
