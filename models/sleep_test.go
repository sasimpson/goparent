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
	mock.On(
		r.Table("sleep").Filter(map[string]interface{}{
			"end":    time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
			"userid": "1",
		}),
	).Return(nil, nil)
	testEnv.DB.Session = mock
	var s Sleep
	status, err := s.Status(&testEnv, &User{ID: "1"})
	mock.AssertExpectations(t)
	assert.Equal(t, false, status)
	assert.Nil(t, err)
}

func TestSleepStatusTrue(t *testing.T) {
	var testEnv config.Env
	mock := r.NewMock()
	mock.On(
		r.Table("sleep").Filter(map[string]interface{}{
			"end":    time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
			"userid": "1",
		}),
	).Return(map[string]interface{}{
		"id":     "1",
		"start":  time.Now().AddDate(0, 0, -1),
		"end":    time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
		"userid": "1",
	}, nil)
	testEnv.DB.Session = mock
	var s Sleep
	status, err := s.Status(&testEnv, &User{ID: "1"})
	mock.AssertExpectations(t)
	assert.Equal(t, true, status)
	assert.Nil(t, err)
}

func TestSleepStatusError(t *testing.T) {
	var testEnv config.Env
	mock := r.NewMock()
	mock.On(
		r.Table("sleep").Filter(map[string]interface{}{
			"end":    time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
			"userid": "1",
		}),
	).Return(map[string]interface{}{}, errors.New("raised error"))
	testEnv.DB.Session = mock
	var s Sleep
	status, err := s.Status(&testEnv, &User{ID: "1"})
	mock.AssertExpectations(t)
	assert.Error(t, err)
	assert.EqualError(t, err, "raised error")
	assert.Equal(t, false, status)
}

func TestSleepStatusEmpty(t *testing.T) {
	var testEnv config.Env
	mock := r.NewMock()
	mock.On(
		r.Table("sleep").Filter(map[string]interface{}{
			"end":    time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
			"userid": "1",
		}),
	).Return(nil, nil)
	testEnv.DB.Session = mock
	var s Sleep
	status, err := s.Status(&testEnv, &User{ID: "1"})
	mock.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, false, status)
}

func TestSleepStart(t *testing.T) {
	var testEnv config.Env
	mock := r.NewMock()
	mock.On(
		r.Table("sleep").Filter(map[string]interface{}{
			"end":    time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
			"userid": "1",
		}),
	).Return(nil, nil)
	testEnv.DB.Session = mock
	var s Sleep
	err := s.Start(&testEnv, &User{ID: "1"})
	assert.Nil(t, err)
}

func TestSleepStartError(t *testing.T) {
	var testEnv config.Env
	mock := r.NewMock()
	mock.On(
		r.Table("sleep").Filter(map[string]interface{}{
			"end":    time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
			"userid": "1",
		}),
	).Return(map[string]interface{}{
		"id":     "1",
		"start":  time.Now().AddDate(0, 0, -1),
		"end":    time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
		"userid": "1",
	}, nil)
	testEnv.DB.Session = mock
	var s Sleep
	err := s.Start(&testEnv, &User{ID: "1"})
	mock.AssertExpectations(t)
	assert.Error(t, err)
	assert.EqualError(t, err, ExistingStartErr.Error())
}

func TestSleepEnd(t *testing.T) {
	var testEnv config.Env
	mock := r.NewMock()
	mock.On(
		r.Table("sleep").Filter(map[string]interface{}{
			"end":    time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
			"userid": "1",
		}),
	).Return(map[string]interface{}{
		"id":     "1",
		"start":  time.Now().AddDate(0, 0, -1),
		"end":    time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
		"userid": "1",
	}, nil)
	testEnv.DB.Session = mock
	var s Sleep
	err := s.End(&testEnv, &User{ID: "1"})
	mock.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestSleepEndError(t *testing.T) {
	var testEnv config.Env
	mock := r.NewMock()
	mock.On(
		r.Table("sleep").Filter(map[string]interface{}{
			"end":    time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
			"userid": "1",
		}),
	).Return(nil, nil)
	testEnv.DB.Session = mock
	var s Sleep
	err := s.End(&testEnv, &User{ID: "1"})
	mock.AssertExpectations(t)
	assert.Error(t, err)
	assert.EqualError(t, err, NoExistingSessionErr.Error())

	mock = r.NewMock()
	mock.On(
		r.Table("sleep").Filter(map[string]interface{}{
			"end":    time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
			"userid": "1",
		}),
	).Return(nil, errors.New("raised error"))
	testEnv.DB.Session = mock
	err = s.End(&testEnv, &User{ID: "1"})
	mock.AssertExpectations(t)
	assert.Error(t, err)
	assert.EqualError(t, err, "raised error")
}

func TestSleepSave(t *testing.T) {
	var testEnv config.Env

	testCases := []struct {
		desc   string
		start  time.Time
		end    time.Time
		userid string
		id     string
	}{
		{
			desc:   "test start save",
			start:  time.Now(),
			userid: "1",
			id:     "1",
		},
		{
			desc:   "test end save",
			start:  time.Now().AddDate(0, 0, -1),
			end:    time.Now(),
			userid: "1",
			id:     "1",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := r.NewMock()
			mock.On(
				r.Table("sleep").Insert(map[string]interface{}{
					"start":  tC.start,
					"end":    tC.end,
					"userid": tC.userid,
				}, r.InsertOpts{Conflict: "replace"}),
			).Return(r.WriteResponse{
				Inserted:      1,
				Errors:        0,
				GeneratedKeys: []string{tC.id},
			}, nil)

			testEnv.DB.Session = mock
			s := Sleep{SleepStart: tC.start, SleepEnd: tC.end, UserID: tC.userid}
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
	mock.On(
		r.Table("sleep").Insert(map[string]interface{}{
			"start":  startTime,
			"end":    endTime,
			"userid": "1",
		}, r.InsertOpts{Conflict: "replace"}),
	).Return(nil, errors.New("raised error"))

	testEnv.DB.Session = mock
	s := Sleep{SleepStart: startTime, SleepEnd: endTime, UserID: "1"}
	err := s.Save(&testEnv)
	mock.AssertExpectations(t)
	assert.Error(t, err)
	assert.EqualError(t, err, "raised error")
}

func TestSleepGetAll(t *testing.T) {
	var testEnv config.Env
	mock := r.NewMock()
	mock.On(
		r.Table("sleep").Filter(map[string]interface{}{
			"userid": "1",
		}).OrderBy(r.Desc("end")),
	).Return(
		map[string]interface{}{
			"start":  time.Now().AddDate(0, 0, -1),
			"end":    time.Now(),
			"userid": "1",
		}, nil,
	)

	testEnv.DB.Session = mock
	var s Sleep
	sleeps, err := s.GetAll(&testEnv, &User{ID: "1"})
	mock.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Len(t, sleeps, 1)
}

func TestSleepGetAllError(t *testing.T) {
	var testEnv config.Env
	mock := r.NewMock()
	mock.On(
		r.Table("sleep").Filter(map[string]interface{}{
			"userid": "1",
		}).OrderBy(r.Desc("end")),
	).Return(nil, errors.New("raised error"))

	testEnv.DB.Session = mock
	var s Sleep
	sleeps, err := s.GetAll(&testEnv, &User{ID: "1"})
	mock.AssertExpectations(t)
	assert.Error(t, err)
	assert.EqualError(t, err, "raised error")
	assert.Nil(t, sleeps)
}
