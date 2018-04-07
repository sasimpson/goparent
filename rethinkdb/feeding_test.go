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

func TestGetFeedings(t *testing.T) {
	testCases := []struct {
		desc         string
		env          *config.Env
		query        *r.MockQuery
		family       *goparent.Family
		resultLength int
		resultError  error
	}{
		{
			desc: "return 1 feeding",
			env:  &config.Env{},
			query: (&r.Mock{}).On(
				r.Table("feeding").Filter(
					map[string]interface{}{
						"familyID": "1",
					}).OrderBy(r.Desc("timestamp")),
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
			query: (&r.Mock{}).On(
				r.Table("feeding").Filter(
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
			desc: "return 0 feeding",
			query: (&r.Mock{}).On(
				r.Table("feeding").Filter(
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
			fs := FeedingService{Env: tC.env}
			feedingResult, err := fs.Feeding(tC.family)
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

func Test(t *testing.T) {
	testCases := []struct {
		desc      string
		env       *config.Env
		timestamp time.Time
	}{
		{
			desc:      "",
			env:       &config.Env{},
			timestamp: time.Now(),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			fs := FeedingService{Env: tC.env}
			feeding := goparent.Feeding{
				Type:      "bottle",
				Amount:    3.5,
				Side:      "",
				FamilyID:  "1",
				UserID:    "1",
				ChildID:   "1",
				TimeStamp: tC.timestamp}
			fs.Save(&feeding)
		})
	}
}

// func TestFeedingSaveError(t *testing.T) {
// 	var testEnv config.Env
// 	timestamp := time.Now()
// 	mock := r.NewMock()
// 	mock.On(
// 		r.Table("feeding").Insert(
// 			map[string]interface{}{
// 				"userID":        "1",
// 				"familyID":      "1",
// 				"childID":       "1",
// 				"timestamp":     timestamp,
// 				"feedingType":   "bottle",
// 				"feedingAmount": 3.5,
// 			}, r.InsertOpts{Conflict: "replace"},
// 		),
// 	).Return(nil, errors.New("returned error"))
// 	testEnv.DB.Session = mock

// 	f := Feeding{
// 		Type:      "bottle",
// 		Amount:    3.5,
// 		Side:      "",
// 		FamilyID:  "1",
// 		UserID:    "1",
// 		ChildID:   "1",
// 		TimeStamp: timestamp}
// 	err := f.Save(&testEnv)
// 	mock.AssertExpectations(t)
// 	assert.Error(t, err)
// 	assert.EqualError(t, err, "returned error")
// }

// func TestFeedingSave(t *testing.T) {
// 	var testEnv config.Env

// 	testCases := []struct {
// 		desc          string
// 		recordID      string
// 		userID        string
// 		familyID      string
// 		childID       string
// 		timestamp     time.Time
// 		feedingType   string
// 		feedingAmount float32
// 		feedingSide   string
// 	}{
// 		{
// 			desc:          "bottle, 3.5floz",
// 			recordID:      "1",
// 			userID:        "1",
// 			familyID:      "1",
// 			childID:       "1",
// 			timestamp:     time.Now(),
// 			feedingType:   "bottle",
// 			feedingAmount: 3.5,
// 			feedingSide:   "",
// 		},
// 		{
// 			desc:          "breast, left side 20min",
// 			recordID:      "2",
// 			userID:        "1",
// 			familyID:      "1",
// 			childID:       "1",
// 			timestamp:     time.Now(),
// 			feedingType:   "breast",
// 			feedingAmount: 20,
// 			feedingSide:   "left",
// 		},
// 	}
// 	for _, tC := range testCases {
// 		t.Run(tC.desc, func(t *testing.T) {
// 			mock := r.NewMock()
// 			mock.On(
// 				r.Table("feeding").Insert(
// 					map[string]interface{}{
// 						"userID":        tC.userID,
// 						"familyID":      tC.familyID,
// 						"childID":       tC.childID,
// 						"timestamp":     tC.timestamp,
// 						"feedingType":   tC.feedingType,
// 						"feedingAmount": tC.feedingAmount,
// 						"feedingSide":   tC.feedingSide,
// 					}, r.InsertOpts{Conflict: "replace"},
// 				).MockAnything(),
// 			).Return(
// 				r.WriteResponse{
// 					Inserted:      1,
// 					Errors:        0,
// 					GeneratedKeys: []string{tC.recordID},
// 				}, nil)
// 			testEnv.DB.Session = mock

// 			f := Feeding{
// 				Type:      tC.feedingType,
// 				Amount:    tC.feedingAmount,
// 				Side:      tC.feedingSide,
// 				UserID:    tC.userID,
// 				ChildID:   tC.childID,
// 				FamilyID:  tC.familyID,
// 				TimeStamp: tC.timestamp,
// 			}
// 			err := f.Save(&testEnv)
// 			mock.AssertExpectations(t)
// 			assert.Nil(t, err)
// 			assert.Equal(t, tC.recordID, f.ID)
// 		})
// 	}
// }
