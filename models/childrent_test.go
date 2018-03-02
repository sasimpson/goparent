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
	//TODO: separate tests?
	var testEnv config.Env
	//test return children
	mock := r.NewMock()
	mock.
		On(
			r.Table("family").Filter(
				func(row r.Term) r.Term {
					return row.Field("members").Contains("1")
				},
			),
		).
		Return(map[string]interface{}{
			"id":           "1",
			"admin":        "1",
			"members":      []string{"1"},
			"created_at":   time.Now(),
			"last_updated": time.Now(),
		}, nil).
		On(
			r.Table("children").
				Filter(map[string]interface{}{
					"familyID": "1",
				}).
				OrderBy(r.Desc("birthday")),
		).
		Return([]interface{}{
			map[string]interface{}{
				"id":       "1",
				"name":     "kiddo1",
				"parentID": 1,
				"birthday": time.Now()},
		}, nil)
	testEnv.DB.Session = mock

	children, err := GetAllChildren(&testEnv, &User{ID: "1"})
	mock.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Len(t, children, 1)

	//test return no children
	mock = r.NewMock()
	mock.
		On(
			r.Table("family").Filter(
				func(row r.Term) r.Term {
					return row.Field("members").Contains("1")
				},
			),
		).
		Return(map[string]interface{}{
			"id":           "1",
			"admin":        "1",
			"members":      []string{"1"},
			"created_at":   time.Now(),
			"last_updated": time.Now(),
		}, nil).
		On(
			r.Table("children").
				Filter(map[string]interface{}{
					"familyID": "1",
				}).
				OrderBy(r.Desc("birthday")),
		).
		Return([]interface{}{}, nil)
	testEnv.DB.Session = mock
	children, err = GetAllChildren(&testEnv, &User{ID: "1"})
	mock.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Len(t, children, 0)

	//test return error
	mock = r.NewMock()
	mock.
		On(
			r.Table("family").Filter(
				func(row r.Term) r.Term {
					return row.Field("members").Contains("1")
				},
			),
		).
		Return(map[string]interface{}{
			"id":           "1",
			"admin":        "1",
			"members":      []string{"1"},
			"created_at":   time.Now(),
			"last_updated": time.Now(),
		}, nil).
		On(
			r.Table("children").
				Filter(map[string]interface{}{
					"familyID": "1",
				}).
				OrderBy(r.Desc("birthday")),
		).
		Return(nil, errors.New("Test Error"))
	testEnv.DB.Session = mock
	children, err = GetAllChildren(&testEnv, &User{ID: "1"})
	mock.AssertExpectations(t)
	assert.Error(t, err)
	assert.Len(t, children, 0)
}

func TestChildrenGetOne(t *testing.T) {
	var testEnv config.Env
	mock := r.NewMock()
	mock.
		On(
			r.Table("family").Filter(
				func(row r.Term) r.Term {
					return row.Field("members").Contains("1")
				},
			),
		).
		Return(map[string]interface{}{
			"id":           "1",
			"admin":        "1",
			"members":      []string{"1"},
			"created_at":   time.Now(),
			"last_updated": time.Now(),
		}, nil).
		On(
			r.Table("children").Filter(map[string]interface{}{"familyID": "1", "id": "1"}),
		).
		Return([]interface{}{
			map[string]interface{}{"id": "1", "name": "kiddo1", "parentID": "1", "familyID": 1, "birthday": time.Now()},
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
	mock.
		On(
			r.Table("children").Insert(
				map[string]interface{}{
					"parentID": "1",
					"familyID": "1",
					"birthday": timestamp,
					"name":     "joey",
				}, r.InsertOpts{Conflict: "replace"},
			),
		).
		Return(nil, errors.New("returned error"))
	testEnv.DB.Session = mock

	c := Child{Name: "joey", ParentID: "1", FamilyID: "1", Birthday: timestamp}
	err := c.Save(&testEnv)
	mock.AssertExpectations(t)
	assert.Error(t, err)
	assert.EqualError(t, err, "returned error")
}

func TestDeleteChild(t *testing.T) {
	var testEnv config.Env
	mock := r.NewMock()
	mock.
		On(
			r.Table("family").Filter(
				func(row r.Term) r.Term {
					return row.Field("members").Contains("1")
				},
			),
		).
		Return(map[string]interface{}{
			"id":           "1",
			"admin":        "1",
			"members":      []string{"1"},
			"created_at":   time.Now(),
			"last_updated": time.Now(),
		}, nil).
		On(
			r.Table("children").Filter(
				map[string]interface{}{
					"familyID": "1",
					"id":       "1",
				},
			).Delete(),
		).
		Return(
			r.WriteResponse{
				Deleted: 1,
			}, nil,
		)
	testEnv.DB.Session = mock

	c := Child{ID: "1", Name: "joey", ParentID: "1", FamilyID: "1", Birthday: time.Now()}
	resp, err := c.DeleteChild(&testEnv, &User{ID: "1"})
	mock.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, 1, resp)

}

func TestChildrenSave(t *testing.T) {
	var testEnv config.Env

	testCases := []struct {
		desc     string
		id       string
		name     string
		parentid string
		familyid string
		birthday time.Time
	}{
		{
			desc:     "joey",
			id:       "1",
			name:     "Joe",
			parentid: "1",
			familyid: "1",
			birthday: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			desc:     "amy",
			id:       "2",
			name:     "Amy",
			parentid: "1",
			familyid: "1",
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
						"familyID": tC.familyid,
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
				FamilyID: tC.familyid,
				Birthday: tC.birthday,
			}

			err := c.Save(&testEnv)
			mock.AssertExpectations(t)
			assert.Nil(t, err)
			assert.Equal(t, tC.id, c.ID)
			assert.Equal(t, tC.parentid, c.ParentID)
			assert.Equal(t, tC.familyid, c.FamilyID)
			assert.Equal(t, tC.birthday, c.Birthday)
		})
	}
}
