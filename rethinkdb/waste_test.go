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

func TestGetWastes(t *testing.T) {
	testCases := []struct {
		desc         string
		env          *config.Env
		query        *r.MockQuery
		family       *goparent.Family
		resultLength int
		resultError  error
	}{
		{
			desc: "return 1 waste",
			env:  &config.Env{},
			query: (&r.Mock{}).On(
				r.Table("waste").Filter(
					map[string]interface{}{
						"familyID": "1",
					}).OrderBy(r.Desc("timestamp")),
			).
				Return([]interface{}{
					map[string]interface{}{
						"id":        "1",
						"wasteType": 1,
						"notes":     "",
						"userID":    "1",
						"familyID":  "1",
						"childID":   "1",
						"timestamp": time.Now(),
					},
				}, nil),
			family:       &goparent.Family{ID: "1"},
			resultLength: 1,
			resultError:  nil,
		},
		{
			desc: "return 0 waste",
			env:  &config.Env{},
			query: (&r.Mock{}).On(
				r.Table("waste").Filter(
					map[string]interface{}{
						"familyID": "1",
					}).OrderBy(r.Desc("timestamp")),
			).
				Return([]interface{}{}, nil),
			family:       &goparent.Family{ID: "1"},
			resultLength: 0,
			resultError:  nil,
		},
		{
			desc: "return 0 waste",
			env:  &config.Env{},
			query: (&r.Mock{}).On(
				r.Table("waste").Filter(
					map[string]interface{}{
						"familyID": "1",
					}).OrderBy(r.Desc("timestamp")),
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
			fs := WasteService{Env: tC.env}
			wasteResult, err := fs.Waste(tC.family)
			if tC.resultError != nil {
				assert.Error(t, err, tC.resultError.Error())
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tC.resultLength, len(wasteResult))
			mock.AssertExpectations(t)
			mock.AssertExecuted(t, tC.query)
		})
	}
}

func TestWasteSave(t *testing.T) {
	timestamp := time.Now()
	testCases := []struct {
		desc        string
		env         *config.Env
		id          string
		query       *r.MockQuery
		timestamp   time.Time
		data        goparent.Waste
		returnError error
	}{
		{
			desc:      "save data",
			env:       &config.Env{},
			id:        "1",
			timestamp: timestamp.Add(time.Hour),
			query: (&r.Mock{}).On(
				r.Table("waste").MockAnything(),
			).Return(
				r.WriteResponse{
					Inserted:      1,
					Errors:        0,
					GeneratedKeys: []string{"1"},
				}, nil),
			data: goparent.Waste{
				ID:        "1",
				Type:      1,
				Notes:     "",
				FamilyID:  "1",
				UserID:    "1",
				ChildID:   "1",
				TimeStamp: timestamp,
			},
		},
		{
			desc:      "save data error",
			env:       &config.Env{},
			timestamp: timestamp.Add(time.Hour),
			query: (&r.Mock{}).On(
				r.Table("waste").MockAnything(),
			).Return(nil, errors.New("returned error")),
			data: goparent.Waste{
				ID:        "1",
				Type:      1,
				Notes:     "",
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
			tC.env.DB = config.DBEnv{Session: mock}
			fs := WasteService{Env: tC.env}
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