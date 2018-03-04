package models

import (
	"errors"
	"testing"
	"time"

	"github.com/sasimpson/goparent/config"
	"github.com/stretchr/testify/assert"
	r "gopkg.in/gorethink/gorethink.v3"
)

func TestSleepStatusFalse(t *testing.T) {
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
			r.Table("sleep").Filter(map[string]interface{}{
				"end":      time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
				"familyID": "1",
			}),
		).Return(nil, nil)
	testEnv.DB.Session = mock
	var s Sleep
	status, err := s.Status(&testEnv, &User{ID: "1", CurrentFamily: "1"})
	mock.AssertExpectations(t)
	assert.Equal(t, false, status)
	assert.Nil(t, err)
}

func TestSleepStatusTrue(t *testing.T) {
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
			r.Table("sleep").Filter(map[string]interface{}{
				"end":      time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
				"familyID": "1",
			}),
		).
		Return(map[string]interface{}{
			"id":     "1",
			"start":  time.Now().AddDate(0, 0, -1),
			"end":    time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
			"userid": "1",
		}, nil)
	testEnv.DB.Session = mock
	var s Sleep
	status, err := s.Status(&testEnv, &User{ID: "1", CurrentFamily: "1"})
	mock.AssertExpectations(t)
	assert.Equal(t, true, status)
	assert.Nil(t, err)
}

func TestSleepStatusError(t *testing.T) {
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
			r.Table("sleep").Filter(map[string]interface{}{
				"end":      time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
				"familyID": "1",
			}),
		).
		Return(map[string]interface{}{}, errors.New("raised error"))
	testEnv.DB.Session = mock
	var s Sleep
	status, err := s.Status(&testEnv, &User{ID: "1", CurrentFamily: "1"})
	mock.AssertExpectations(t)
	assert.Error(t, err)
	assert.EqualError(t, err, "raised error")
	assert.Equal(t, false, status)
}

func TestSleepStatusEmpty(t *testing.T) {
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
			r.Table("sleep").Filter(map[string]interface{}{
				"end":      time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
				"familyID": "1",
			}),
		).Return(nil, nil)
	testEnv.DB.Session = mock
	var s Sleep
	status, err := s.Status(&testEnv, &User{ID: "1", CurrentFamily: "1"})
	mock.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, false, status)
}

func TestSleepStart(t *testing.T) {
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
			r.Table("sleep").Filter(map[string]interface{}{
				"end":      time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
				"familyID": "1",
			}),
		).
		Return(nil, nil)
	testEnv.DB.Session = mock
	var s Sleep
	err := s.SleepStart(&testEnv, &User{ID: "1", CurrentFamily: "1"})
	assert.Nil(t, err)
}

func TestSleepStartError(t *testing.T) {
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
			r.Table("sleep").Filter(map[string]interface{}{
				"end":      time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
				"familyID": "1",
			}),
		).
		Return(map[string]interface{}{
			"id":     "1",
			"start":  time.Now().AddDate(0, 0, -1),
			"end":    time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
			"userid": "1",
		}, nil)
	testEnv.DB.Session = mock
	var s Sleep
	err := s.SleepStart(&testEnv, &User{ID: "1", CurrentFamily: "1"})
	mock.AssertExpectations(t)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrExistingStart.Error())
}

func TestSleepEnd(t *testing.T) {
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
			r.Table("sleep").Filter(map[string]interface{}{
				"end":      time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
				"familyID": "1",
			}),
		).
		Return(map[string]interface{}{
			"id":       "1",
			"start":    time.Now().AddDate(0, 0, -1),
			"end":      time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
			"familyID": "1",
			"userID":   "1",
			"childID":  "1",
		}, nil)
	testEnv.DB.Session = mock
	var s Sleep
	err := s.SleepEnd(&testEnv, &User{ID: "1", CurrentFamily: "1"})
	mock.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestSleepEndError(t *testing.T) {
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
			r.Table("sleep").Filter(map[string]interface{}{
				"end":      time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
				"familyID": "1",
			}),
		).
		Return(nil, nil)
	testEnv.DB.Session = mock
	var s Sleep
	err := s.SleepEnd(&testEnv, &User{ID: "1", CurrentFamily: "1"})
	mock.AssertExpectations(t)
	assert.Error(t, err)
	assert.EqualError(t, err, ErrNoExistingSession.Error())

	mock = r.NewMock()
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
			r.Table("sleep").Filter(map[string]interface{}{
				"end":      time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
				"familyID": "1",
			}),
		).
		Return(nil, errors.New("raised error"))
	testEnv.DB.Session = mock
	err = s.SleepEnd(&testEnv, &User{ID: "1", CurrentFamily: "1"})
	mock.AssertExpectations(t)
	assert.Error(t, err)
	assert.EqualError(t, err, "raised error")
}

func TestSleepSave(t *testing.T) {
	var testEnv config.Env

	testCases := []struct {
		desc     string
		start    time.Time
		end      time.Time
		userid   string
		childid  string
		familyid string
		id       string
	}{
		{
			desc:     "test start save",
			start:    time.Now(),
			userid:   "1",
			childid:  "1",
			familyid: "1",
			id:       "1",
		},
		{
			desc:     "test end save",
			start:    time.Now().AddDate(0, 0, -1),
			end:      time.Now(),
			userid:   "1",
			childid:  "1",
			familyid: "1",
			id:       "1",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := r.NewMock()
			mock.On(
				r.Table("sleep").Insert(map[string]interface{}{
					"start":    tC.start,
					"end":      tC.end,
					"userID":   tC.userid,
					"childID":  tC.childid,
					"familyID": tC.familyid,
				}, r.InsertOpts{Conflict: "replace"}),
			).Return(r.WriteResponse{
				Inserted:      1,
				Errors:        0,
				GeneratedKeys: []string{tC.id},
			}, nil)

			testEnv.DB.Session = mock
			s := Sleep{Start: tC.start, End: tC.end, UserID: tC.userid, FamilyID: tC.familyid, ChildID: tC.childid}
			err := s.Save(&testEnv)
			mock.AssertExpectations(t)
			assert.Nil(t, err)
		})
	}
}

func TestSleepSaveError(t *testing.T) {
	var testEnv config.Env
	startTime := time.Now().AddDate(0, 0, -1)
	endTime := time.Now()
	mock := r.NewMock()
	mock.
		On(
			r.Table("sleep").Insert(map[string]interface{}{
				"start":    startTime,
				"end":      endTime,
				"familyID": "1",
				"childID":  "1",
				"userID":   "1",
			}, r.InsertOpts{Conflict: "replace"}),
		).Return(nil, errors.New("raised error"))

	testEnv.DB.Session = mock
	s := Sleep{Start: startTime, End: endTime, UserID: "1", ChildID: "1", FamilyID: "1"}
	err := s.Save(&testEnv)
	mock.AssertExpectations(t)
	assert.Error(t, err)
	assert.EqualError(t, err, "raised error")
}

func TestSleepGetAll(t *testing.T) {
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
			r.Table("sleep").Filter(map[string]interface{}{
				"familyID": "1",
			}).OrderBy(r.Desc("end")),
		).
		Return(
			map[string]interface{}{
				"start":    time.Now().AddDate(0, 0, -1),
				"end":      time.Now(),
				"userID":   "1",
				"childID":  "1",
				"familyID": "1",
			}, nil,
		)

	testEnv.DB.Session = mock
	var s Sleep
	sleeps, err := s.GetAll(&testEnv, &User{ID: "1", CurrentFamily: "1"})
	mock.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Len(t, sleeps, 1)
}

func TestSleepGetAllError(t *testing.T) {
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
			r.Table("sleep").Filter(map[string]interface{}{
				"familyID": "1",
			}).OrderBy(r.Desc("end")),
		).
		Return(nil, errors.New("raised error"))

	testEnv.DB.Session = mock
	var s Sleep
	sleeps, err := s.GetAll(&testEnv, &User{ID: "1", CurrentFamily: "1"})
	mock.AssertExpectations(t)
	assert.Error(t, err)
	assert.EqualError(t, err, "raised error")
	assert.Nil(t, sleeps)
}
