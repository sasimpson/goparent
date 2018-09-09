package rethinkdb

import (
	"errors"
	"testing"
	"time"

	"github.com/sasimpson/goparent"
	"github.com/stretchr/testify/assert"
	r "gopkg.in/gorethink/gorethink.v3"
)

func TestGetFeedings(t *testing.T) {
	testCases := []struct {
		desc         string
		env          *goparent.Env
		query        *r.MockQuery
		family       *goparent.Family
		resultLength int
		resultError  error
	}{
		{
			desc: "return 1 feeding",
			env:  &goparent.Env{},
			query: (&r.Mock{}).On(
				r.Table("feeding").MockAnything(),
			).
				Return([]interface{}{
					map[string]interface{}{
						"id":            "1",
						"feedingType":   "bottle",
						"feedingAmount": 1,
						"feedingSide":   "",
						"userid":        "1",
						"timestamp":     time.Now(),
					},
				}, nil),
			family:       &goparent.Family{ID: "1"},
			resultLength: 1,
			resultError:  nil,
		},
		{
			desc: "return 0 feeding",
			env:  &goparent.Env{},
			query: (&r.Mock{}).On(
				r.Table("feeding").MockAnything(),
			).
				Return([]interface{}{}, nil),
			family:       &goparent.Family{ID: "1"},
			resultLength: 0,
			resultError:  nil,
		},
		{
			desc: "return 0 feeding",
			env:  &goparent.Env{},
			query: (&r.Mock{}).On(
				r.Table("feeding").MockAnything(),
			).
				Return([]interface{}{}, errors.New("unknown error")),
			family:       &goparent.Family{ID: "1"},
			resultLength: 0,
			resultError:  errors.New("unknown error"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)
			fs := FeedingService{Env: tC.env, DB: &DBEnv{Session: mock}}
			feedingResult, err := fs.Feeding(tC.family, 7)
			if tC.resultError != nil {
				assert.Error(t, err, tC.resultError.Error())
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tC.resultLength, len(feedingResult))
			mock.AssertExpectations(t)
			mock.AssertExecuted(t, tC.query)
		})
	}
}

func TestFeedingSave(t *testing.T) {
	timestamp := time.Now()
	testCases := []struct {
		desc        string
		env         *goparent.Env
		id          string
		query       *r.MockQuery
		timestamp   time.Time
		data        goparent.Feeding
		returnError error
	}{
		{
			desc:      "save data",
			env:       &goparent.Env{},
			id:        "1",
			timestamp: timestamp.Add(time.Hour),
			query: (&r.Mock{}).On(
				r.Table("feeding").Insert(
					map[string]interface{}{
						"userID":        "1",
						"familyID":      "1",
						"childID":       "1",
						"timestamp":     timestamp.Add(time.Hour),
						"feedingType":   "bottle",
						"feedingAmount": 3.5,
					}, r.InsertOpts{Conflict: "replace"},
				),
			).Return(
				r.WriteResponse{
					Inserted:      1,
					Errors:        0,
					GeneratedKeys: []string{"1"},
				}, nil),
			data: goparent.Feeding{
				Type:      "bottle",
				Amount:    3.5,
				Side:      "",
				FamilyID:  "1",
				UserID:    "1",
				ChildID:   "1",
				TimeStamp: timestamp.Add(time.Hour),
			},
		},
		{
			desc:      "save data",
			env:       &goparent.Env{},
			timestamp: timestamp.Add(time.Hour),
			query: (&r.Mock{}).On(
				r.Table("feeding").Insert(
					map[string]interface{}{
						"userID":        "1",
						"familyID":      "1",
						"childID":       "1",
						"timestamp":     timestamp.Add(time.Hour),
						"feedingType":   "bottle",
						"feedingAmount": 3.5,
					}, r.InsertOpts{Conflict: "replace"},
				),
			).Return(nil, errors.New("returned error")),
			data: goparent.Feeding{
				Type:      "bottle",
				Amount:    3.5,
				Side:      "",
				FamilyID:  "1",
				UserID:    "1",
				ChildID:   "1",
				TimeStamp: timestamp.Add(time.Hour),
			},
			returnError: errors.New("returned error"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)
			fs := FeedingService{Env: tC.env, DB: &DBEnv{Session: mock}}
			err := fs.Save(&tC.data)
			if tC.returnError != nil {
				assert.Error(t, err, tC.returnError)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tC.id, tC.data.ID)
			}

		})
	}
}
