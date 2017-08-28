package models

import (
	"errors"
	"testing"
	"time"

	"github.com/sasimpson/goparent/config"
	"github.com/stretchr/testify/assert"
	r "gopkg.in/gorethink/gorethink.v3"
)

var testEnv config.Env

func TestFeedingGetAll(t *testing.T) {
	mock := r.NewMock()
	mock.On(
		r.Table("feeding").Filter(map[string]interface{}{"userid": "1"}).OrderBy(r.Desc("timestamp")),
	).Return([]interface{}{
		map[string]interface{}{"id": "1", "feedingType": "bottle", "feedingAmount": 1, "feedingSide": "", "userid": "1", "timestamp": time.Now()},
	}, nil)
	testEnv.DB = config.DBEnv{Session: mock}

	var f Feeding
	feedings, err := f.GetAll(&testEnv, &User{ID: "1"})
	mock.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Len(t, feedings, 1)

	mock = r.NewMock()
	mock.On(
		r.Table("feeding").Filter(map[string]interface{}{"userid": "1"}).OrderBy(r.Desc("timestamp")),
	).Return([]interface{}{}, nil)
	testEnv.DB.Session = mock
	feedings, err = f.GetAll(&testEnv, &User{ID: "1"})
	mock.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Len(t, feedings, 0)

	mock = r.NewMock()
	mock.On(
		r.Table("feeding").Filter(map[string]interface{}{"userid": "1"}).OrderBy(r.Desc("timestamp")),
	).Return([]interface{}{}, errors.New("Test Error"))
	testEnv.DB.Session = mock
	feedings, err = f.GetAll(&testEnv, &User{ID: "1"})
	mock.AssertExpectations(t)
	assert.Error(t, err)
	assert.Len(t, feedings, 0)
}

func TestSaveError(t *testing.T) {
	timestamp := time.Now()
	mock := r.NewMock()
	mock.On(
		r.Table("feeding").Insert(
			map[string]interface{}{
				"userid":        "1",
				"timestamp":     timestamp,
				"feedingType":   "bottle",
				"feedingAmount": 3.5,
			}, r.InsertOpts{Conflict: "replace"},
		),
	).Return(nil, errors.New("returned error"))
	testEnv.DB.Session = mock

	f := Feeding{Type: "bottle", Amount: 3.5, Side: "", UserID: "1", TimeStamp: timestamp}
	err := f.Save(&testEnv)
	mock.AssertExpectations(t)
	assert.Error(t, err)
	assert.EqualError(t, err, "returned error")
}

func TestSave(t *testing.T) {
	testCases := []struct {
		desc          string
		recordID      string
		userID        string
		timestamp     time.Time
		feedingType   string
		feedingAmount float32
		feedingSide   string
	}{
		{
			desc:          "bottle, 3.5floz",
			recordID:      "1",
			userID:        "1",
			timestamp:     time.Now(),
			feedingType:   "bottle",
			feedingAmount: 3.5,
			feedingSide:   "",
		},
		{
			desc:          "breast, left side 20min",
			recordID:      "2",
			userID:        "1",
			timestamp:     time.Now(),
			feedingType:   "breast",
			feedingAmount: 20,
			feedingSide:   "left",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := r.NewMock()
			mock.On(
				r.Table("feeding").Insert(
					map[string]interface{}{
						"userid":        tC.userID,
						"timestamp":     tC.timestamp,
						"feedingType":   tC.feedingType,
						"feedingAmount": tC.feedingAmount,
						"feedingSide":   tC.feedingSide,
					}, r.InsertOpts{Conflict: "replace"},
				).MockAnything(),
			).Return(
				r.WriteResponse{
					Inserted:      1,
					Errors:        0,
					GeneratedKeys: []string{tC.recordID},
				}, nil)
			testEnv.DB.Session = mock

			f := Feeding{
				Type:      tC.feedingType,
				Amount:    tC.feedingAmount,
				Side:      tC.feedingSide,
				UserID:    tC.userID,
				TimeStamp: tC.timestamp,
			}
			err := f.Save(&testEnv)
			mock.AssertExpectations(t)
			assert.Nil(t, err)
			assert.Equal(t, tC.recordID, f.ID)
		})
	}
}
