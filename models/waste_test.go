package models

import (
	"errors"
	"testing"
	"time"

	"github.com/sasimpson/goparent/config"
	"github.com/stretchr/testify/assert"
	r "gopkg.in/gorethink/gorethink.v3"
)

func TestWasteGetAll(t *testing.T) {
	var testEnv config.Env
	mock := r.NewMock()
	mock.
		On(
			r.Table("family").Get("1"),
		).
		Return(map[string]interface{}{
			"id":           "1",
			"admin":        "1",
			"members":      []string{"1"},
			"created_at":   time.Now(),
			"last_updated": time.Now(),
		}, nil).
		On(
			r.Table("waste").Filter(
				map[string]interface{}{
					"familyID": "1",
				}).
				OrderBy(r.Desc("timestamp")),
		).
		Return([]interface{}{
			map[string]interface{}{
				"id":        "1",
				"wasteType": 1,
				"notes":     "test note",
				"userid":    "1",
				"timestamp": time.Now(),
			},
		}, nil)
	testEnv.DB.Session = mock
	var w Waste
	wastes, err := w.GetAll(&testEnv, &User{ID: "1", CurrentFamily: "1"})
	mock.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Len(t, wastes, 1)
}

func TestWasteGetByID(t *testing.T) {
	var testEnv config.Env
	mock := r.NewMock()
	mock.
		On(
			r.Table("waste").Get("1"),
		).Return([]interface{}{
		map[string]interface{}{
			"id":        "1",
			"wasteType": 1,
			"notes":     "test note",
			"userID":    "1",
			"familyID":  "1",
			"childID":   "1",
			"timestamp": time.Now(),
		},
	}, nil)
	testEnv.DB.Session = mock
	var w Waste
	err := w.GetByID(&testEnv, "1")
	mock.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, "1", w.ID)
}

func TestWasteSaveError(t *testing.T) {
	//TODO: verify output
	var testEnv config.Env
	timestamp := time.Now().UTC()
	mock := r.NewMock()
	mock.On(
		r.Table("waste").Insert(
			map[string]interface{}{
				"wasteType": 1,
				"notes":     "Some Notes",
				"userID":    "1",
				"familyID":  "1",
				"childID":   "1",
				"timestamp": timestamp,
			}, r.InsertOpts{Conflict: "replace"},
		).MockAnything(),
	).Return(nil, errors.New("returned error"))
	testEnv.DB.Session = mock

	w := Waste{UserID: "1", Type: 1, TimeStamp: timestamp, Notes: "Some Notes"}
	err := w.Save(&testEnv)
	mock.AssertExpectations(t)
	assert.Error(t, err)
	assert.EqualError(t, err, "returned error")
}

func TestWasteSave(t *testing.T) {
	//TODO: verify output
	var testEnv config.Env

	testCases := []struct {
		desc      string
		wasteType int
		notes     string
		userid    string
		familyid  string
		childid   string
		timestamp time.Time
		recordID  string
	}{
		{
			desc:      "waste type 1",
			wasteType: 1,
			notes:     "some waste test notes",
			userid:    "1",
			familyid:  "1",
			childid:   "1",
			timestamp: time.Now().UTC(),
			recordID:  "1",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := r.NewMock()
			mock.On(
				r.Table("waste").Insert(
					map[string]interface{}{
						"userID":    tC.userid,
						"childID":   tC.childid,
						"familyID":  tC.familyid,
						"wasteType": tC.wasteType,
						"notes":     tC.notes,
						"timestamp": tC.timestamp,
					},
				).MockAnything(),
			).Return(
				r.WriteResponse{
					Inserted:      1,
					Errors:        0,
					GeneratedKeys: []string{tC.recordID},
				}, nil,
			)

			testEnv.DB.Session = mock
			w := Waste{
				Type:      tC.wasteType,
				Notes:     tC.notes,
				UserID:    tC.userid,
				FamilyID:  tC.familyid,
				ChildID:   tC.childid,
				TimeStamp: tC.timestamp,
			}

			err := w.Save(&testEnv)
			mock.AssertExpectations(t)
			assert.Nil(t, err)
			assert.Equal(t, tC.recordID, w.ID)
		})
	}
}
