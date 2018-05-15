package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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

func TestUserNewInviteHandler(t *testing.T) {
	testCases := []struct {
		desc              string
		env               *config.Env
		userInviteService goparent.UserInvitationService
		inviteUser        string
		contextUser       *goparent.User
		formErr           bool
		responseCode      int
	}{
		{
			desc:         "invite fails auth",
			env:          &config.Env{},
			responseCode: http.StatusUnauthorized,
		},
		{
			desc:         "invite parse form error",
			env:          &config.Env{},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			formErr:      true,
			responseCode: http.StatusInternalServerError,
		},
		{
			desc:         "invite no email",
			env:          &config.Env{},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			responseCode: http.StatusBadRequest,
		},
		{
			desc: "existing invite",
			env:  &config.Env{},
			userInviteService: &mock.MockUserInvitationService{
				InviteParentErr: goparent.ErrExistingInvitation,
			},
			inviteUser:   "invitedUser@test.com",
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			responseCode: http.StatusConflict,
		},
		{
			desc: "unknown invite error",
			env:  &config.Env{},
			userInviteService: &mock.MockUserInvitationService{
				InviteParentErr: errors.New("unknown error"),
			},
			inviteUser:   "invitedUser@test.com",
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			responseCode: http.StatusInternalServerError,
		},
		{
			desc:              "successful invite",
			env:               &config.Env{},
			userInviteService: &mock.MockUserInvitationService{},
			inviteUser:        "invitedUser@test.com",
			contextUser:       &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			responseCode:      http.StatusCreated,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockHandler := Handler{
				Env: tC.env,
				UserInvitationService: tC.userInviteService,
			}

			form := url.Values{}
			form.Add("email", tC.inviteUser)
			var req *http.Request
			if tC.formErr != true {
				req, _ = http.NewRequest("POST", "/user/invite", strings.NewReader(form.Encode()))
			} else {
				req, _ = http.NewRequest("POST", "/user/invite", nil)
			}
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			req.Form = form

			handler := mockHandler.userNewInviteHandler()
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

		})
	}
}

func TestListInviteHandler(t *testing.T) {
	testCases := []struct {
		desc                  string
		env                   *config.Env
		contextUser           *goparent.User
		userInvitationService goparent.UserInvitationService
		responseCode          int
		resultLength          int
	}{
		{
			desc:         "bad auth",
			responseCode: http.StatusUnauthorized,
		},
		{
			desc:        "get SentInvites error",
			contextUser: &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			userInvitationService: &mock.MockUserInvitationService{
				SentInvitesErr: errors.New("test error"),
			},
			responseCode: http.StatusInternalServerError,
		},
		{
			desc:        "get pending invites error",
			contextUser: &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			userInvitationService: &mock.MockUserInvitationService{
				InvitesErr: errors.New("test error"),
			},
			responseCode: http.StatusInternalServerError,
		},
		{
			desc:                  "empty response",
			contextUser:           &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			userInvitationService: &mock.MockUserInvitationService{},
			responseCode:          http.StatusOK,
			resultLength:          0,
		},
		{
			desc:        "non-empty response",
			contextUser: &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			userInvitationService: &mock.MockUserInvitationService{
				GetSentInvites: []*goparent.UserInvitation{
					&goparent.UserInvitation{
						ID:          "1",
						UserID:      "1",
						InviteEmail: "testuser@test.com",
						Timestamp:   time.Now(),
					},
				},
				GetInvites: []*goparent.UserInvitation{
					&goparent.UserInvitation{
						ID:          "2",
						UserID:      "2",
						InviteEmail: "testowner@test.com",
						Timestamp:   time.Now(),
					},
				},
			},
			responseCode: http.StatusOK,
			resultLength: 2,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockHandler := Handler{
				Env: tC.env,
				UserInvitationService: tC.userInvitationService,
			}

			handler := mockHandler.userListInviteHandler()
			req, _ := http.NewRequest("GET", "/user/invite", nil)
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
		})
	}
}

func TestAcceptInviteHandler(t *testing.T) {
	testCases := []struct {
		desc                  string
		env                   *config.Env
		contextUser           *goparent.User
		userInvitationService goparent.UserInvitationService
		responseCode          int
	}{
		{
			desc:         "bad auth",
			env:          &config.Env{},
			responseCode: http.StatusUnauthorized,
		},
		{
			desc:        "accept fail",
			env:         &config.Env{},
			contextUser: &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			userInvitationService: &mock.MockUserInvitationService{
				AcceptErr: errors.New("test error"),
			},
			responseCode: http.StatusInternalServerError,
		},
		{
			desc:                  "accept success",
			env:                   &config.Env{},
			contextUser:           &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			userInvitationService: &mock.MockUserInvitationService{},
			responseCode:          http.StatusNoContent,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockHandler := Handler{
				Env: tC.env,
				UserInvitationService: tC.userInvitationService,
			}

			handler := mockHandler.userAcceptInviteHandler()
			req, _ := http.NewRequest("GET", "/user/invite/accept/1", nil)
			req = mux.SetURLVars(req, map[string]string{"id": "1"})
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
		})
	}
}

func TestDeleteInviteHandler(t *testing.T) {
	testCases := []struct {
		desc                  string
		env                   *config.Env
		contextUser           *goparent.User
		userInvitationService goparent.UserInvitationService
		responseCode          int
	}{
		{
			desc:         "bad auth",
			env:          &config.Env{},
			responseCode: http.StatusUnauthorized,
		},
		{
			desc:        "get invite fail",
			env:         &config.Env{},
			contextUser: &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			userInvitationService: &mock.MockUserInvitationService{
				InviteErr: errors.New("not found"),
			},
			responseCode: http.StatusNotFound,
		},
		{
			desc:        "delete fail",
			env:         &config.Env{},
			contextUser: &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			userInvitationService: &mock.MockUserInvitationService{
				GetInvite: &goparent.UserInvitation{
					ID:          "1",
					UserID:      "1",
					InviteEmail: "testuser@test.com",
					Timestamp:   time.Now(),
				},
				DeleteErr: errors.New("test error"),
			},
			responseCode: http.StatusInternalServerError,
		},
		{
			desc:        "delete success",
			env:         &config.Env{},
			contextUser: &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			userInvitationService: &mock.MockUserInvitationService{
				GetInvite: &goparent.UserInvitation{
					ID:          "1",
					UserID:      "1",
					InviteEmail: "testuser@test.com",
					Timestamp:   time.Now(),
				},
			},
			responseCode: http.StatusNoContent,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockHandler := Handler{
				Env: tC.env,
				UserInvitationService: tC.userInvitationService,
			}

			handler := mockHandler.userDeleteInviteHandler()
			req, _ := http.NewRequest("DELETE", "/user/invite/1", nil)
			req = mux.SetURLVars(req, map[string]string{"id": "1"})
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
		})
	}
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
