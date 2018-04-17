package rethinkdb

import (
	"testing"

	"github.com/sasimpson/goparent"
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

	us := UserService{Env: &testEnv}
	user, err := us.User("1")
	mock.AssertExpectations(t)
	assert.EqualError(t, err, "no result for that id")
	assert.Nil(t, user)
}

func TestGetUser(t *testing.T) {
	var testEnv config.Env
	mock := r.NewMock()
	mock.On(
		r.Table("users").Get("1"),
	).Return([]interface{}{map[string]interface{}{
		"id":            "1",
		"name":          "test user",
		"email":         "testuser@test.com",
		"username":      "testuser",
		"password":      "testpassword",
		"currentFamily": "1",
	}}, nil)

	testEnv.DB = config.DBEnv{Session: mock}

	us := UserService{Env: &testEnv}
	user, err := us.User("1")
	mock.AssertExpectations(t)
	assert.Nil(t, err)
	assert.NotNil(t, user)
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
		"id":            "1",
		"name":          "test user",
		"email":         "testuser@test.com",
		"username":      "testuser",
		"password":      "testpassword",
		"currentFamily": "1",
	}}, nil)

	testEnv.DB = config.DBEnv{Session: mock}

	us := UserService{Env: &testEnv}
	user, err := us.UserByLogin("testuser@test.com", "testpassword")
	mock.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, "1", user.ID)
	assert.Equal(t, "test user", user.Name)
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
	us := UserService{Env: &testEnv}
	user, err := us.UserByLogin("testuser@test.com", "testpassword")
	mock.AssertExpectations(t)
	assert.EqualError(t, err, "no result for that username password combo")
	assert.Nil(t, user)
}

func TestNewUserSave(t *testing.T) {
	var testEnv config.Env
	mock := r.NewMock()
	mock.On(
		r.Table("users").Filter(map[string]interface{}{
			"email": "testuser@test.com",
		})).
		Return(nil, nil).
		On(
			r.Table("family").MockAnything(),
		).Once().
		Return(
			r.WriteResponse{
				Inserted:      1,
				Errors:        0,
				GeneratedKeys: []string{"1"},
			}, nil).
		On(
			r.Table("users").Insert(
				map[string]interface{}{
					"name":          "test user",
					"email":         "testuser@test.com",
					"username":      "testuser",
					"password":      "testpassword",
					"currentFamily": "1",
				}, r.InsertOpts{Conflict: "replace"},
			),
		).
		Return(
			r.WriteResponse{
				Inserted:      1,
				Errors:        0,
				GeneratedKeys: []string{"1"},
			}, nil).
		On(
			r.Table("family").MockAnything(),
		).Once().
		Return(
			r.WriteResponse{
				Updated: 1,
			}, nil)

	testEnv.DB = config.DBEnv{Session: mock}

	user := goparent.User{
		Name:     "test user",
		Email:    "testuser@test.com",
		Username: "testuser",
		Password: "testpassword",
	}

	us := UserService{Env: &testEnv}
	err := us.Save(&user)
	mock.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, "1", user.ID)
}

func TestUserSave(t *testing.T) {
	var testEnv config.Env
	mock := r.NewMock()
	mock.On(
		r.Table("users").Filter(map[string]interface{}{
			"email": "testuser@test.com",
		})).
		Return(map[string]interface{}{
			"id":    "1",
			"email": "testuser@test.com",
		}, nil).
		On(
			r.Table("family").MockAnything(),
		).Once().
		Return(
			map[string]interface{}{
				"id":      "1",
				"admin":   "1",
				"members": []string{"1"},
			}, nil).
		On(
			r.Table("users").Insert(
				map[string]interface{}{
					"name":          "test user",
					"email":         "testuser@test.com",
					"username":      "testuser",
					"password":      "testpassword",
					"currentFamily": "1",
					"id":            "1",
				}, r.InsertOpts{Conflict: "replace"},
			),
		).
		Return(
			r.WriteResponse{
				Updated: 1,
			}, nil)
		// On(
		// 	r.Table("family").MockAnything(),
		// ).Once().
		// Return(
		// 	r.WriteResponse{
		// 		Updated: 1,
		// 	}, nil)

	testEnv.DB = config.DBEnv{Session: mock}

	user := goparent.User{
		ID:       "1",
		Name:     "test user",
		Email:    "testuser@test.com",
		Username: "testuser",
		Password: "testpassword",
	}

	us := UserService{Env: &testEnv}
	err := us.Save(&user)
	mock.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, "1", user.ID)
}

func TestTokens(t *testing.T) {
	var testEnv config.Env
	u := goparent.User{
		ID:            "1",
		Name:          "test user",
		Email:         "testuser@test.com",
		Username:      "testuser",
		Password:      "testpassword",
		CurrentFamily: "1",
	}
	testEnv.Auth.SigningKey = []byte("testkey")
	us := UserService{Env: &testEnv}
	token, err := us.GetToken(&u)
	assert.Nil(t, err)
	if assert.NotNil(t, token) {
		assert.NotEqual(t, "", token)
	}

	mock := r.NewMock()
	mock.On(
		r.Table("users").Get("1"),
	).Return([]interface{}{map[string]interface{}{
		"id":            "1",
		"name":          "test user",
		"email":         "testuser@test.com",
		"username":      "testuser",
		"password":      "testpassword",
		"currentFamily": "1",
	}}, nil)
	testEnv.DB.Session = mock

	user, ok, err := us.ValidateToken(token)
	mock.AssertExpectations(t)
	assert.Nil(t, err)
	assert.True(t, ok)
	assert.EqualValues(t, &u, user)
}

func TestGetAllFamily(t *testing.T) {
	testCases := []struct {
		desc        string
		env         *config.Env
		query       *r.MockQuery
		user        *goparent.User
		resultSize  int
		resultError error
	}{
		{
			desc: "get all families",
			env:  &config.Env{},
			query: (&r.Mock{}).On(
				r.Table("family").Filter(
					func(row r.Term) r.Term {
						return row.Field("members").Contains("1")
					},
				),
			).Return(
				map[string]interface{}{"id": "1", "admin": "1", "members": []string{"1"}}, nil,
			),
			user:        &goparent.User{ID: "1"},
			resultSize:  1,
			resultError: nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)
			tC.env.DB = config.DBEnv{Session: mock}
			us := UserService{Env: tC.env}
			families, err := us.GetAllFamily(tC.user)
			if tC.resultError != nil {
				assert.Error(t, err, tC.resultError.Error())
			} else {
				assert.Nil(t, err)
				assert.Len(t, families, tC.resultSize)
			}
		})
	}
}

//move to api tests
// func TestUserFromContext(t *testing.T) {
// 	type contextKey string
// 	var userContextKey contextKey = "user"
// 	var ctx context.Context
// 	ctx = context.WithValue(ctx, userContextKey, User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
// 	user, err := UserFromContext(ctx)

// 	assert.Nil(t, err)
// 	assert.Equal(t, "1", user.ID)
// }
