package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/sasimpson/goparent/config"
	"github.com/sasimpson/goparent/models"
	"github.com/stretchr/testify/assert"

	r "gopkg.in/gorethink/gorethink.v3"
)

func TestLoginHandler(t *testing.T) {
	email := "testuser@test.com"
	password := "testpassword"

	var testEnv config.Env
	mock := r.NewMock()
	mock.On(
		r.Table("users").Filter(map[string]interface{}{
			"email":    "testuser@test.com",
			"password": "testpassword",
		}),
	).Return(
		map[string]interface{}{
			"id":       "1",
			"name":     "test user",
			"email":    "testuser@test.com",
			"username": "testuser",
		}, nil,
	)
	testEnv.DB.Session = mock
	params := url.Values{"username": {email}, "password": {password}}
	req, err := http.NewRequest("POST", "/user/login", bytes.NewBufferString(params.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	handler := loginHandler(&testEnv)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestInvalidLogin(t *testing.T) {
	email := "testuser@test.com"
	password := "testpassword"

	var testEnv config.Env
	mock := r.NewMock()
	mock.On(
		r.Table("users").Filter(map[string]interface{}{
			"email":    "testuser@test.com",
			"password": "testpassword",
		}),
	).Return(
		nil, nil,
	)
	testEnv.DB.Session = mock
	params := url.Values{"username": {email}, "password": {password}}
	req, err := http.NewRequest("POST", "/user/login", bytes.NewBufferString(params.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	handler := loginHandler(&testEnv)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestUserGetHandler(t *testing.T) {
	//TODO: verify output
	var testEnv config.Env
	mock := r.NewMock()
	mock.
		On(
			r.Table("family").Filter(
				func(row r.Term) r.Term {
					return row.Field("members").Contains("1")
				},
			),
		).
		Return(map[string]interface{}{
			"id":           "1",
			"admin":        "1",
			"members":      []string{"1"},
			"created_at":   time.Now(),
			"last_updated": time.Now(),
		}, nil)

	testEnv.DB.Session = mock
	req, err := http.NewRequest("GET", "/user/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := userGetHandler(&testEnv)
	rr := httptest.NewRecorder()
	ctx := req.Context()
	ctx = context.WithValue(ctx, userContextKey, models.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	var testResult UserResponse
	err = json.NewDecoder(rr.Body).Decode(&testResult)
	assert.Equal(t, "1", testResult.UserData.ID)
}

func TestUserNewHandler(t *testing.T) {
	//TODO: verify output
	var testEnv config.Env
	mock := r.NewMock()
	mock.
		On(
			r.Table("users").Filter(map[string]interface{}{
				"email": "testuser@test.com",
			}),
		).
		On(
			r.Table("users").Insert(
				map[string]interface{}{
					"name":     "test user",
					"email":    "testuser@test.com",
					"username": "testuser",
					"password": "testpassword",
				},
				r.InsertOpts{Conflict: "replace"},
			),
		).
		Return(
			r.WriteResponse{
				Inserted:      1,
				Errors:        0,
				GeneratedKeys: []string{"1"},
			}, nil,
		)
	testEnv.DB.Session = mock

	js := bytes.NewBufferString(`{ "userData": {"name": "test user", "email": "testuser@test.com", "username": "testuser", "password": "testpassword"}}`)
	req, err := http.NewRequest("POST", "/user", js)
	if err != nil {
		t.Fatal(err)
	}

	handler := userNewHandler(&testEnv)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestNewInviteHandler(t *testing.T) {
	var testEnv config.Env
	mock := r.NewMock()
	mock.
		On(
			r.Table("invites").MockAnything(),
		).
		Return(nil, nil).
		On(
			r.Table("invites").Insert(
				map[string]interface{}{
					"userID":      "1",
					"inviteEmail": "inviteuser@test.com",
					"timestamp":   time.Now(),
				},
				r.InsertOpts{Conflict: "replace"},
			),
		).
		Return(
			r.WriteResponse{
				Inserted:      1,
				Errors:        0,
				GeneratedKeys: []string{"1"},
			}, nil,
		)
	testEnv.DB.Session = mock

	form := url.Values{}
	form.Add("email", "inviteuser@test.com")
	req, _ := http.NewRequest("POST", "/user/invite", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	req.Form = form

	handler := userNewInviteHandler(&testEnv)
	rr := httptest.NewRecorder()
	ctx := req.Context()
	ctx = context.WithValue(ctx, userContextKey, models.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
	req = req.WithContext(ctx)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)

}

func TestListInviteHandler(t *testing.T) {
	var testEnv config.Env
	mock := r.NewMock()
	mock.
		On(
			r.Table("users").Get("1"),
		).
		Return([]map[string]interface{}{
			{"id": "1",
				"name":     "test user",
				"email":    "testuser@test.com",
				"username": "testuser",
				"password": "testpassword"},
		}, nil).
		On(
			r.Table("invites").Filter(map[string]interface{}{
				"userID": "1",
			}).OrderBy(r.Desc("timestamp")),
		).
		Return([]map[string]interface{}{
			{"id": "1", "userID": "1", "inviteEmail": "testinvite1@test.com", "timestamp": time.Now()},
			{"id": "2", "userID": "1", "inviteEmail": "testinvite2@test.com", "timestamp": time.Now()},
		}, nil)
	testEnv.DB.Session = mock

	req, _ := http.NewRequest("GET", "/invite", nil)
	req.Header.Add("Content-Type", jsonContentType)

	handler := userListInviteHandler(&testEnv)
	rr := httptest.NewRecorder()
	ctx := req.Context()
	ctx = context.WithValue(ctx, userContextKey, models.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
	req = req.WithContext(ctx)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteInviteHandler(t *testing.T) {
	var testEnv config.Env
	mock := r.NewMock()
	mock.
		On(
			r.Table("users").Get("1"),
		).
		Return([]map[string]interface{}{
			{
				"id":       "1",
				"name":     "test user",
				"email":    "testuser@test.com",
				"username": "testuser",
				"password": "testpassword"},
		}, nil).
		On(
			r.Table("invites").Filter(map[string]interface{}{
				"userID": "1",
				"id":     "1",
			}).Delete(),
		).
		Return(
			r.WriteResponse{
				Deleted: 1,
			}, nil)
	testEnv.DB.Session = mock

	req, err := http.NewRequest("DELETE", "/user/invite/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	//this is needed because apparently gorilla mux doesn't do ^ and we have to set it ahead of time to work.
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	handler := userDeleteInviteHandler(&testEnv)
	rr := httptest.NewRecorder()
	ctx := req.Context()
	ctx = context.WithValue(ctx, userContextKey, models.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
	req = req.WithContext(ctx)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNoContent, rr.Code)
}
func TestInitUsersHandlers(t *testing.T) {
	//TODO: update with new handler routes
	testCases := []struct {
		desc    string
		name    string
		path    string
		methods []string
	}{
		{
			desc:    "user new",
			name:    "UserNew",
			path:    "/user/",
			methods: []string{"POST"},
		},
		{
			desc:    "user login",
			name:    "UserLogin",
			path:    "/user/login",
			methods: []string{"POST"},
		},
	}

	var testEnv config.Env
	routes := mux.NewRouter()
	initUsersHandlers(&testEnv, routes)

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			route := routes.Get(tC.name)
			path, _ := route.GetPathTemplate()
			methods, _ := route.GetMethods()
			assert.Equal(t, tC.name, route.GetName())
			assert.Equal(t, tC.path, path)
			assert.Equal(t, tC.methods, methods)
		})
	}
}
