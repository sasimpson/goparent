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

func TestGetWastes(t *testing.T) {
	testCases := []struct {
		desc         string
		env          *goparent.Env
		query        *r.MockQuery
		family       *goparent.Family
		resultLength int
		resultError  error
	}{
		{
			desc: "return 1 waste",
			env:  &goparent.Env{DB: &mock.DBEnv{}},
			query: (&r.Mock{}).On(
				r.Table("waste").MockAnything(),
				// Filter(
				// 	map[string]interface{}{
				// 		"familyID": "1",
				// 	}).OrderBy(r.Desc("timestamp")),
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
			env:  &goparent.Env{DB: &mock.DBEnv{}},
			query: (&r.Mock{}).On(
				r.Table("waste").MockAnything(),
				// r.Table("waste").Filter(
				// 	map[string]interface{}{
				// 		"familyID": "1",
				// 	}).OrderBy(r.Desc("timestamp")),
			).
				Return([]interface{}{}, nil),
			family:       &goparent.Family{ID: "1"},
			resultLength: 0,
			resultError:  nil,
		},
		{
			desc: "return 0 waste",
			env:  &goparent.Env{DB: &mock.DBEnv{}},
			query: (&r.Mock{}).On(
				r.Table("waste").MockAnything(),
				// r.Table("waste").Filter(
				// 	map[string]interface{}{
				// 		"familyID": "1",
				// 	}).OrderBy(r.Desc("timestamp")),
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
			fs := WasteService{Env: tC.env, DB: &DBEnv{Session: mock}}
			wasteResult, err := fs.Waste(ctx, tC.family, 7)
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
		env         *goparent.Env
		id          string
		query       *r.MockQuery
		timestamp   time.Time
		data        goparent.Waste
		returnError error
	}{
		{
			desc:      "save data",
			env:       &goparent.Env{DB: &mock.DBEnv{}},
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
			env:       &goparent.Env{DB: &mock.DBEnv{}},
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
			ctx := context.Background()
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)
			fs := WasteService{Env: tC.env, DB: &DBEnv{Session: mock}}
			err := fs.Save(ctx, &tC.data)
			if tC.returnError != nil {
				assert.EqualError(t, tC.returnError, err.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tC.id, tC.data.ID)
			}
		})
	}
}

func TestStats(t *testing.T) {
	testCases := []struct {
		desc        string
		env         *goparent.Env
		query       *r.MockQuery
		child       *goparent.Child
		returnError error
	}{
		{
			desc: "get stats",
			env:  &goparent.Env{DB: &mock.DBEnv{}},
			query: (&r.Mock{}).On(
				r.Table("waste").MockAnything(),
			).Return([]map[string]interface{}{
				{"id": "1"},
				{"id": "2"},
			}, nil),
			child:       &goparent.Child{ID: "1"},
			returnError: nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctx := context.Background()
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)
			fs := WasteService{Env: tC.env, DB: &DBEnv{Session: mock}}
			statsData, err := fs.Stats(ctx, tC.child)
			if tC.returnError != nil {
				assert.EqualError(t, tC.returnError, err.Error())
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, statsData)
			}
		})
	}
}

func TestWasteGraph(t *testing.T) {
	testCases := []struct {
		desc        string
		env         *goparent.Env
		query       *r.MockQuery
		child       *goparent.Child
		returnError *error
	}{
		{
			desc: "get graph data",
			env:  &goparent.Env{DB: &mock.DBEnv{}},
			query: (&r.Mock{}).On(
				r.Table("waste").MockAnything(),
			).Return([]map[string]interface{}{
				{"group": []int{2018, 1, 1, 1}, "reduction": 1},
				{"group": []int{2018, 1, 1, 2}, "reduction": 2},
			}, nil),
			child: &goparent.Child{ID: "1", Name: "Billy"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctx := context.Background()
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)
			fs := WasteService{Env: tC.env, DB: &DBEnv{Session: mock}}
			chartData, err := fs.GraphData(ctx, tC.child)
			if tC.returnError != nil {
				assert.EqualError(t, *tC.returnError, err.Error())
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, chartData)
			}
		})
	}
}
