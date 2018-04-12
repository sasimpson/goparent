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

func TestGetSleep(t *testing.T) {
	timestamp := time.Now()
	testCases := []struct {
		desc         string
		env          *config.Env
		query        *r.MockQuery
		family       *goparent.Family
		resultLength int
		resultError  error
	}{
		{
			desc: "return 1 sleep",
			env:  &config.Env{},
			query: (&r.Mock{}).On(
				r.Table("sleep").Filter(
					map[string]interface{}{
						"familyID": "1",
					}).OrderBy(r.Desc("end")),
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
			env:  &config.Env{},
			query: (&r.Mock{}).On(
				r.Table("sleep").Filter(
					map[string]interface{}{
						"familyID": "1",
					}).OrderBy(r.Desc("end")),
			).
				Return([]interface{}{}, nil),
			family:       &goparent.Family{ID: "1"},
			resultLength: 0,
			resultError:  nil,
		},
		{
			desc: "return sleep error",
			env:  &config.Env{},
			query: (&r.Mock{}).On(
				r.Table("sleep").Filter(
					map[string]interface{}{
						"familyID": "1",
					}).OrderBy(r.Desc("end")),
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
			tC.env.DB = config.DBEnv{Session: mock}
			fs := SleepService{Env: tC.env}
			sleepResult, err := fs.Sleep(tC.family)
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
		env         *config.Env
		id          string
		query       *r.MockQuery
		timestamp   time.Time
		data        goparent.Sleep
		returnError error
	}{
		{
			desc:      "save data",
			env:       &config.Env{},
			id:        "1",
			timestamp: timestamp.Add(time.Hour),
			query: (&r.Mock{}).On(
				r.Table("sleep").Insert(
					map[string]interface{}{
						"userID":   "1",
						"familyID": "1",
						"childID":  "1",
						"start":    timestamp,
						"end":      timestamp.Add(time.Hour),
					}, r.InsertOpts{Conflict: "replace"},
				),
			).Return(
				r.WriteResponse{
					Inserted:      1,
					Errors:        0,
					GeneratedKeys: []string{"1"},
				}, nil),
			data: goparent.Sleep{
				FamilyID: "1",
				ChildID:  "1",
				UserID:   "1",
				Start:    timestamp,
				End:      timestamp.Add(time.Hour),
			},
		},
		{
			desc:      "save data error",
			env:       &config.Env{},
			timestamp: timestamp.Add(time.Hour),
			query: (&r.Mock{}).On(
				r.Table("sleep").Insert(
					map[string]interface{}{
						"userID":   "1",
						"familyID": "1",
						"childID":  "1",
						"start":    timestamp,
						"end":      timestamp.Add(time.Hour),
					}, r.InsertOpts{Conflict: "replace"},
				),
			).Return(nil, errors.New("returned error")),
			data: goparent.Sleep{
				FamilyID: "1",
				ChildID:  "1",
				UserID:   "1",
				Start:    timestamp,
				End:      timestamp.Add(time.Hour),
			},
			returnError: errors.New("returned error"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)
			tC.env.DB = config.DBEnv{Session: mock}
			fs := SleepService{Env: tC.env}
			err := fs.Save(&tC.data)
			if tC.returnError != nil {
				assert.Error(t, err, tC.returnError)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestStatus(t *testing.T) {
	testCases := []struct {
		desc        string
		env         *config.Env
		query       *r.MockQuery
		family      *goparent.Family
		child       *goparent.Child
		result      bool
		returnError error
	}{
		{
			desc: "get status true",
			env:  &config.Env{},
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
			result: true,
		},
		{
			desc: "get status false",
			env:  &config.Env{},
			query: (&r.Mock{}).On(
				r.Table("sleep").Filter(map[string]interface{}{
					"end":      time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
					"familyID": "1",
					"childID":  "1",
				}),
			).Return(nil, nil),
			family: &goparent.Family{ID: "1"},
			child:  &goparent.Child{ID: "1"},
			result: false,
		},
		{
			desc: "get status empty result",
			env:  &config.Env{},
			query: (&r.Mock{}).On(
				r.Table("sleep").Filter(map[string]interface{}{
					"end":      time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
					"familyID": "1",
					"childID":  "1",
				}),
			).Return(map[string]interface{}{}, r.ErrEmptyResult),
			family:      &goparent.Family{ID: "1"},
			child:       &goparent.Child{ID: "1"},
			result:      false,
			returnError: nil,
		},
		{
			desc: "get status err result",
			env:  &config.Env{},
			query: (&r.Mock{}).On(
				r.Table("sleep").Filter(map[string]interface{}{
					"end":      time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
					"familyID": "1",
					"childID":  "1",
				}),
			).Return(nil, errors.New("test error")),
			family:      &goparent.Family{ID: "1"},
			child:       &goparent.Child{ID: "1"},
			result:      false,
			returnError: errors.New("test error"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)
			tC.env.DB = config.DBEnv{Session: mock}
			fs := SleepService{Env: tC.env}
			status, err := fs.Status(tC.family, tC.child)
			if tC.returnError != nil {
				assert.Error(t, err, tC.returnError)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tC.result, status)
			}

		})
	}
}
