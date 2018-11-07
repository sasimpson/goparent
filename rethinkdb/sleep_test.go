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

func TestGetSleep(t *testing.T) {
	timestamp := time.Now()
	testCases := []struct {
		desc         string
		env          *goparent.Env
		query        *r.MockQuery
		family       *goparent.Family
		resultLength int
		resultError  error
	}{
		{
			desc: "return 1 sleep",
			env:  &goparent.Env{DB: &mock.DBEnv{}},
			query: (&r.Mock{}).On(
				r.Table("sleep").MockAnything(),
			).
				Return([]interface{}{
					map[string]interface{}{
						"id":       "1",
						"start":    timestamp,
						"end":      timestamp.AddDate(0, 0, 1),
						"childID":  "1",
						"familyID": "1",
						"userID":   "1",
					},
				}, nil),
			family:       &goparent.Family{ID: "1"},
			resultLength: 1,
			resultError:  nil,
		},
		{
			desc: "return 0 sleep",
			env:  &goparent.Env{DB: &mock.DBEnv{}},
			query: (&r.Mock{}).On(
				r.Table("sleep").MockAnything(),
			).
				Return([]interface{}{}, nil),
			family:       &goparent.Family{ID: "1"},
			resultLength: 0,
			resultError:  nil,
		},
		{
			desc: "return sleep error",
			env:  &goparent.Env{DB: &mock.DBEnv{}},
			query: (&r.Mock{}).On(
				r.Table("sleep").MockAnything(),
			).
				Return([]interface{}{}, errors.New("unknown error")),
			family:       &goparent.Family{ID: "1"},
			resultLength: 0,
			resultError:  errors.New("unknown error"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctx := context.Background()
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)
			fs := SleepService{Env: tC.env, DB: &DBEnv{Session: mock}}
			sleepResult, err := fs.Sleep(ctx, tC.family, 7)
			if tC.resultError != nil {
				assert.Error(t, err, tC.resultError.Error())
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tC.resultLength, len(sleepResult))
			mock.AssertExpectations(t)
			mock.AssertExecuted(t, tC.query)
		})
	}
}

func TestSleepSave(t *testing.T) {
	timestamp := time.Now()
	testCases := []struct {
		desc        string
		env         *goparent.Env
		id          string
		query       *r.MockQuery
		timestamp   time.Time
		data        goparent.Sleep
		returnError error
	}{
		{
			desc:      "save data",
			env:       &goparent.Env{DB: &mock.DBEnv{}},
			id:        "1",
			timestamp: timestamp.Add(time.Hour),
			query: (&r.Mock{}).On(
				r.Table("sleep").Insert(
					map[string]interface{}{
						"userID":      "1",
						"familyID":    "1",
						"childID":     "1",
						"start":       timestamp,
						"end":         timestamp.Add(time.Hour),
						"createdAt":   timestamp,
						"lastUpdated": timestamp,
					}, r.InsertOpts{Conflict: "replace"},
				),
			).Return(
				r.WriteResponse{
					Inserted:      1,
					Errors:        0,
					GeneratedKeys: []string{"1"},
				}, nil),
			data: goparent.Sleep{
				FamilyID:    "1",
				ChildID:     "1",
				UserID:      "1",
				Start:       timestamp,
				End:         timestamp.Add(time.Hour),
				CreatedAt:   timestamp,
				LastUpdated: timestamp,
			},
		},
		{
			desc:      "save data error",
			env:       &goparent.Env{DB: &mock.DBEnv{}},
			timestamp: timestamp.Add(time.Hour),
			query: (&r.Mock{}).On(
				r.Table("sleep").Insert(
					map[string]interface{}{
						"userID":      "1",
						"familyID":    "1",
						"childID":     "1",
						"start":       timestamp,
						"end":         timestamp.Add(time.Hour),
						"createdAt":   timestamp,
						"lastUpdated": timestamp,
					}, r.InsertOpts{Conflict: "replace"},
				),
			).Return(nil, errors.New("returned error")),
			data: goparent.Sleep{
				FamilyID:    "1",
				ChildID:     "1",
				UserID:      "1",
				Start:       timestamp,
				End:         timestamp.Add(time.Hour),
				CreatedAt:   timestamp,
				LastUpdated: timestamp,
			},
			returnError: errors.New("returned error"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctx := context.Background()
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)

			fs := SleepService{Env: tC.env, DB: &DBEnv{Session: mock}}
			err := fs.Save(ctx, &tC.data)
			if tC.returnError != nil {
				assert.Error(t, tC.returnError, err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestStatus(t *testing.T) {
	testCases := []struct {
		desc        string
		env         *goparent.Env
		query       *r.MockQuery
		family      *goparent.Family
		child       *goparent.Child
		sleep       *goparent.Sleep
		result      bool
		returnError error
	}{
		{
			desc: "get status true",
			env:  &goparent.Env{DB: &mock.DBEnv{}},
			query: (&r.Mock{}).On(
				r.Table("sleep").Filter(map[string]interface{}{
					"end":      time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
					"familyID": "1",
					"childID":  "1",
				}),
			).Return(map[string]interface{}{
				"id": "1",
			}, nil),
			family: &goparent.Family{ID: "1"},
			child:  &goparent.Child{ID: "1"},
			sleep:  &goparent.Sleep{ID: "1"},
			result: true,
		},
		{
			desc: "get status false",
			env:  &goparent.Env{DB: &mock.DBEnv{}},
			query: (&r.Mock{}).On(
				r.Table("sleep").Filter(map[string]interface{}{
					"end":      time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
					"familyID": "1",
					"childID":  "1",
				}),
			).Return(nil, nil),
			family: &goparent.Family{ID: "1"},
			child:  &goparent.Child{ID: "1"},
			sleep:  &goparent.Sleep{ID: "1"},
			result: false,
		},
		{
			desc: "get status empty result",
			env:  &goparent.Env{DB: &mock.DBEnv{}},
			query: (&r.Mock{}).On(
				r.Table("sleep").Filter(map[string]interface{}{
					"end":      time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
					"familyID": "1",
					"childID":  "1",
				}),
			).Return(map[string]interface{}{}, r.ErrEmptyResult),
			family:      &goparent.Family{ID: "1"},
			child:       &goparent.Child{ID: "1"},
			sleep:       &goparent.Sleep{ID: "1"},
			result:      false,
			returnError: nil,
		},
		{
			desc: "get status err result",
			env:  &goparent.Env{DB: &mock.DBEnv{}},
			query: (&r.Mock{}).On(
				r.Table("sleep").Filter(map[string]interface{}{
					"end":      time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
					"familyID": "1",
					"childID":  "1",
				}),
			).Return(nil, errors.New("test error")),
			family:      &goparent.Family{ID: "1"},
			child:       &goparent.Child{ID: "1"},
			sleep:       &goparent.Sleep{ID: "1"},
			result:      false,
			returnError: errors.New("test error"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctx := context.Background()
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)

			ss := SleepService{Env: tC.env, DB: &DBEnv{Session: mock}}
			status, err := ss.Status(ctx, tC.family, tC.child)
			sleepStart := ss.Start(ctx, tC.sleep, tC.family, tC.child)
			sleepEnd := ss.End(ctx, tC.sleep, tC.family, tC.child)
			if tC.returnError != nil {
				assert.EqualError(t, tC.returnError, err.Error())
				assert.Error(t, sleepStart)
				assert.Error(t, sleepEnd)

			} else {
				assert.Nil(t, err)
				assert.Equal(t, tC.result, status)

				switch tC.result {
				case true:
					assert.Error(t, sleepStart)
					assert.Nil(t, sleepEnd)
					break
				case false:
					assert.Error(t, sleepEnd)
					assert.Nil(t, sleepStart)
				}
			}

		})
	}
}
