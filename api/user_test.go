package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/config"
	"github.com/sasimpson/goparent/mock"
	"github.com/stretchr/testify/assert"
)

func TestLoginHandler(t *testing.T) {
	testCases := []struct {
		desc         string
		env          *config.Env
		email        string
		password     string
		userService  goparent.UserService
		responseCode int
	}{
		{
			desc:     "user logged in and token issued",
			env:      &config.Env{},
			email:    "testuser@test.com",
			password: "testpassword",
			userService: &mock.MockUserService{
				ReturnedUser: &goparent.User{
					ID:       "1",
					Name:     "test user",
					Email:    "testuser@test.com",
					Username: "testuser",
				},
				Token: "this-is-a-token",
			},
			responseCode: http.StatusOK,
		},
		{
			desc:     "login error",
			env:      &config.Env{},
			email:    "testuser@test.com",
			password: "testpassword",
			userService: &mock.MockUserService{
				ReturnedUser: &goparent.User{
					ID:       "1",
					Name:     "test user",
					Email:    "testuser@test.com",
					Username: "testuser",
				},
				AuthErr: errors.New("invalid login"),
			},
			responseCode: http.StatusUnauthorized,
		},
		{
			desc:     "token issue error",
			env:      &config.Env{},
			email:    "testuser@test.com",
			password: "testpassword",
			userService: &mock.MockUserService{
				ReturnedUser: &goparent.User{
					ID:       "1",
					Name:     "test user",
					Email:    "testuser@test.com",
					Username: "testuser",
				},
				TokenErr: errors.New("token error"),
			},
			responseCode: http.StatusInternalServerError,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockHandler := Handler{
				Env:         tC.env,
				UserService: tC.userService,
			}
			params := url.Values{"username": {tC.email}, "password": {tC.password}}
			req, err := http.NewRequest("POST", "/user/login", bytes.NewBufferString(params.Encode()))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			if err != nil {
				t.Fatal(err)
			}
			handler := mockHandler.loginHandler()
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tC.responseCode, rr.Code)
		})
	}
}

func TestUserGetHandler(t *testing.T) {
	testCases := []struct {
		desc         string
		env          *config.Env
		userService  goparent.UserService
		contextUser  *goparent.User
		responseCode int
	}{
		{
			desc: "valid user",
			env:  &config.Env{},
			userService: &mock.MockUserService{
				Family: &goparent.Family{
					ID:          "1",
					Admin:       "1",
					Members:     []string{"1"},
					CreatedAt:   time.Now(),
					LastUpdated: time.Now(),
				},
			},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			responseCode: http.StatusOK,
		},
		{
			desc: "invalid user",
			env:  &config.Env{},
			userService: &mock.MockUserService{
				Family: &goparent.Family{
					ID:          "1",
					Admin:       "1",
					Members:     []string{"1"},
					CreatedAt:   time.Now(),
					LastUpdated: time.Now(),
				},
			},
			contextUser:  nil,
			responseCode: http.StatusUnauthorized,
		},
		{
			desc: "family error",
			env:  &config.Env{},
			userService: &mock.MockUserService{
				FamilyErr: errors.New("family error"),
			},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			responseCode: http.StatusInternalServerError,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockHandler := Handler{
				Env:         tC.env,
				UserService: tC.userService,
			}
			req, err := http.NewRequest("GET", "/user/", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := mockHandler.userGetHandler()
			rr := httptest.NewRecorder()

			ctx := req.Context()
			if tC.contextUser != nil {
				ctx = context.WithValue(ctx, userContextKey, tC.contextUser)
			} else {
				ctx = context.WithValue(ctx, userContextKey, "")
			}

			req = req.WithContext(ctx)
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tC.responseCode, rr.Code)
			if tC.responseCode/100 == 2 {
				var result UserResponse
				decoder := json.NewDecoder(rr.Body)
				err = decoder.Decode(&result)
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, tC.contextUser.ID, result.UserData.ID)
			}
		})
	}
}

func TestUserNewHandler(t *testing.T) {
	testCases := []struct {
		desc         string
		env          *config.Env
		userService  goparent.UserService
		userRequest  *UserRequest
		responseCode int
	}{
		{
			desc:         "invalid json",
			env:          &config.Env{},
			userService:  &mock.MockUserService{},
			responseCode: http.StatusInternalServerError,
		},
		{
			desc: "invalid user save",
			env:  &config.Env{},
			userService: &mock.MockUserService{
				SaveErr: errors.New("user exists"),
			},
			userRequest: &UserRequest{
				UserData: &goparent.User{
					Name:     "test user",
					Email:    "testuser@test.com",
					Username: "testuser",
					Password: "testpassword",
				},
			},
			responseCode: http.StatusConflict,
		},
		{
			desc: "valid user save",
			env:  &config.Env{},
			userService: &mock.MockUserService{
				UserID: "1",
			},
			userRequest: &UserRequest{
				UserData: &goparent.User{
					Name:     "test user",
					Email:    "testuser@test.com",
					Username: "testuser",
					Password: "testpassword",
				},
			},
			responseCode: http.StatusOK,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockHandler := Handler{
				Env:         tC.env,
				UserService: tC.userService,
			}

			var js []byte
			var err error
			switch tC.userRequest {
			case nil:
				js = []byte("this is a test")
			default:
				js, err = json.Marshal(&tC.userRequest)
				if err != nil {
					t.Fatal(err)
				}
			}
			req, err := http.NewRequest("POST", "/user", bytes.NewReader(js))
			if err != nil {
				t.Fatal(err)
			}

			handler := mockHandler.userNewHandler()
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tC.responseCode, rr.Code)
		})
	}
}

// func TestUserNewHandler(t *testing.T) {
// 	//TODO: verify output
// 	var testEnv config.Env
// 	mock := r.NewMock()
// 	mock.
// 		On(
// 			r.Table("users").Filter(map[string]interface{}{
// 				"email": "testuser@test.com",
// 			}),
// 		).
// 		On(
// 			r.Table("users").Insert(
// 				map[string]interface{}{
// 					"name":     "test user",
// 					"email":    "testuser@test.com",
// 					"username": "testuser",
// 					"password": "testpassword",
// 				},
// 				r.InsertOpts{Conflict: "replace"},
// 			),
// 		).
// 		Return(
// 			r.WriteResponse{
// 				Inserted:      1,
// 				Errors:        0,
// 				GeneratedKeys: []string{"1"},
// 			}, nil,
// 		)
// 	testEnv.DB.Session = mock

// 	js := bytes.NewBufferString(`{ "userData": {"name": "test user", "email": "testuser@test.com", "username": "testuser", "password": "testpassword"}}`)
// 	req, err := http.NewRequest("POST", "/user", js)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	handler := userNewHandler(&testEnv)
// 	rr := httptest.NewRecorder()

// 	handler.ServeHTTP(rr, req)
// 	assert.Equal(t, http.StatusOK, rr.Code)
// }

// func TestNewInviteHandler(t *testing.T) {
// 	var testEnv config.Env
// 	mock := r.NewMock()
// 	mock.
// 		On(
// 			r.Table("invites").MockAnything(),
// 		).
// 		Return(nil, nil).
// 		On(
// 			r.Table("invites").Insert(
// 				map[string]interface{}{
// 					"userID":      "1",
// 					"inviteEmail": "inviteuser@test.com",
// 					"timestamp":   time.Now(),
// 				},
// 				r.InsertOpts{Conflict: "replace"},
// 			),
// 		).
// 		Return(
// 			r.WriteResponse{
// 				Inserted:      1,
// 				Errors:        0,
// 				GeneratedKeys: []string{"1"},
// 			}, nil,
// 		)
// 	testEnv.DB.Session = mock

// 	form := url.Values{}
// 	form.Add("email", "inviteuser@test.com")
// 	req, _ := http.NewRequest("POST", "/user/invite", strings.NewReader(form.Encode()))
// 	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

// 	req.Form = form

// 	handler := userNewInviteHandler(&testEnv)
// 	rr := httptest.NewRecorder()
// 	ctx := req.Context()
// 	ctx = context.WithValue(ctx, userContextKey, goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
// 	req = req.WithContext(ctx)

// 	handler.ServeHTTP(rr, req)
// 	assert.Equal(t, http.StatusCreated, rr.Code)

// }

// func TestListInviteHandler(t *testing.T) {
// 	var testEnv config.Env
// 	mock := r.NewMock()
// 	mock.
// 		On(
// 			r.Table("users").Get("1"),
// 		).
// 		Return([]map[string]interface{}{
// 			{"id": "1",
// 				"name":     "test user",
// 				"email":    "testuser@test.com",
// 				"username": "testuser",
// 				"password": "testpassword"},
// 		}, nil).
// 		On(
// 			r.Table("invites").Filter(map[string]interface{}{
// 				"userID": "1",
// 			}).OrderBy(r.Desc("timestamp")),
// 		).
// 		Return([]map[string]interface{}{
// 			{"id": "1", "userID": "1", "inviteEmail": "testinvite1@test.com", "timestamp": time.Now()},
// 			{"id": "2", "userID": "1", "inviteEmail": "testinvite2@test.com", "timestamp": time.Now()},
// 		}, nil)
// 	testEnv.DB.Session = mock

// 	req, _ := http.NewRequest("GET", "/user/invite", nil)
// 	req.Header.Add("Content-Type", jsonContentType)

// 	handler := userListInviteHandler(&testEnv)
// 	rr := httptest.NewRecorder()
// 	ctx := req.Context()
// 	ctx = context.WithValue(ctx, userContextKey, goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
// 	req = req.WithContext(ctx)

// 	handler.ServeHTTP(rr, req)
// 	assert.Equal(t, http.StatusOK, rr.Code)
// }

// func TestDeleteInviteHandler(t *testing.T) {
// 	var testEnv config.Env
// 	mock := r.NewMock()
// 	mock.
// 		On(
// 			r.Table("users").Get("1"),
// 		).
// 		Return([]map[string]interface{}{
// 			{
// 				"id":       "1",
// 				"name":     "test user",
// 				"email":    "testuser@test.com",
// 				"username": "testuser",
// 				"password": "testpassword"},
// 		}, nil).
// 		On(
// 			r.Table("invites").Filter(map[string]interface{}{
// 				"userID": "1",
// 				"id":     "1",
// 			}).Delete(),
// 		).
// 		Return(
// 			r.WriteResponse{
// 				Deleted: 1,
// 			}, nil)
// 	testEnv.DB.Session = mock

// 	req, err := http.NewRequest("DELETE", "/user/invite/1", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	//this is needed because apparently gorilla mux doesn't do ^ and we have to set it ahead of time to work.
// 	req = mux.SetURLVars(req, map[string]string{"id": "1"})

// 	handler := userDeleteInviteHandler(&testEnv)
// 	rr := httptest.NewRecorder()
// 	ctx := req.Context()
// 	ctx = context.WithValue(ctx, userContextKey, goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
// 	req = req.WithContext(ctx)

// 	handler.ServeHTTP(rr, req)
// 	assert.Equal(t, http.StatusNoContent, rr.Code)
// }
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
	h := Handler{Env: &testEnv}
	routes := mux.NewRouter()
	h.initUsersHandlers(routes)

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
