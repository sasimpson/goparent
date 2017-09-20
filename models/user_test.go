package models

import (
	"testing"

	"github.com/sasimpson/goparent/config"
	"github.com/stretchr/testify/assert"
	r "gopkg.in/gorethink/gorethink.v3"
)

func TestGetNoUser(t *testing.T) {
	var testEnv config.Env
	mock := r.NewMock()
	mock.On(
		r.Table("users").Get("1"),
	).Return([]interface{}{}, nil)

	testEnv.DB = config.DBEnv{Session: mock}

	var u User
	err := u.GetUser(&testEnv, "1")
	mock.AssertExpectations(t)
	assert.EqualError(t, err, "no result for that id")
}

func TestGetUser(t *testing.T) {
	var testEnv config.Env
	mock := r.NewMock()
	mock.On(
		r.Table("users").Get("1"),
	).Return([]interface{}{map[string]interface{}{
		"id":       "1",
		"name":     "test user",
		"email":    "testuser@test.com",
		"username": "testuser",
		"password": "testpassword",
	}}, nil)

	testEnv.DB = config.DBEnv{Session: mock}

	var u User
	err := u.GetUser(&testEnv, "1")
	mock.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestGetUserByLogin(t *testing.T) {
	var testEnv config.Env
	mock := r.NewMock()
	mock.On(
		r.Table("users").Filter(map[string]interface{}{
			"email":    "testuser@test.com",
			"password": "testpassword",
		}),
	).Return([]interface{}{map[string]interface{}{
		"id":       "1",
		"name":     "test user",
		"email":    "testuser@test.com",
		"username": "testuser",
		"password": "testpassword",
	}}, nil)

	testEnv.DB = config.DBEnv{Session: mock}

	var u User
	err := u.GetUserByLogin(&testEnv, "testuser@test.com", "testpassword")
	mock.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, "1", u.ID)
	assert.Equal(t, "test user", u.Name)
}

func TestGetUserByLoginError(t *testing.T) {
	var testEnv config.Env
	mock := r.NewMock()
	mock.On(
		r.Table("users").Filter(map[string]interface{}{
			"email":    "testuser@test.com",
			"password": "testpassword",
		}),
	).Return([]interface{}{}, nil)

	testEnv.DB = config.DBEnv{Session: mock}

	var u User
	err := u.GetUserByLogin(&testEnv, "testuser@test.com", "testpassword")
	mock.AssertExpectations(t)
	assert.EqualError(t, err, "no result for that username password combo")
}

func TestUserSave(t *testing.T) {
	var testEnv config.Env
	mock := r.NewMock()
	mock.On(
		r.Table("users").Filter(map[string]interface{}{
			"email": "testuser@test.com",
		}),
	).On(
		r.Table("users").Insert(
			map[string]interface{}{
				"name":     "test user",
				"email":    "testuser@test.com",
				"username": "testuser",
				"password": "testpassword",
			}, r.InsertOpts{Conflict: "replace"},
		),
	).Return(
		r.WriteResponse{
			Inserted:      1,
			Errors:        0,
			GeneratedKeys: []string{"1"},
		}, nil)

	testEnv.DB = config.DBEnv{Session: mock}

	u := User{
		Name:     "test user",
		Email:    "testuser@test.com",
		Username: "testuser",
		Password: "testpassword",
	}
	err := u.Save(&testEnv)
	mock.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, "1", u.ID)
}