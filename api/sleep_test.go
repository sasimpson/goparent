package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sasimpson/goparent"

	"github.com/gorilla/mux"
	"github.com/sasimpson/goparent/mock"
	"github.com/stretchr/testify/assert"
)

func TestSleepGetHandler(t *testing.T) {
	testCases := []struct {
		desc          string
		env           *goparent.Env
		userService   goparent.UserService
		familyService goparent.FamilyService
		sleepService  goparent.SleepService
		contextUser   *goparent.User
		responseCode  int
		resultLength  int
	}{
		{
			desc:          "returns auth error",
			env:           &goparent.Env{},
			userService:   &mock.UserService{},
			familyService: &mock.FamilyService{},
			sleepService:  &mock.SleepService{},
			contextUser:   nil,
			responseCode:  http.StatusUnauthorized,
		},
		{
			desc: "returns family error",
			env:  &goparent.Env{},
			userService: &mock.UserService{
				FamilyErr: errors.New("test error"),
			},
			familyService: &mock.FamilyService{},
			sleepService:  &mock.SleepService{},
			contextUser:   &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			responseCode:  http.StatusInternalServerError,
		},
		{
			desc: "returns sleep error",
			env:  &goparent.Env{},
			userService: &mock.UserService{
				Family: &goparent.Family{
					ID:          "1",
					Admin:       "1",
					Members:     []string{"1"},
					CreatedAt:   time.Now(),
					LastUpdated: time.Now(),
				},
			},

			sleepService: &mock.SleepService{
				GetErr: errors.New("test error"),
			},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			responseCode: http.StatusInternalServerError,
		},
		{
			desc: "returns no sleep",
			env:  &goparent.Env{},
			userService: &mock.UserService{
				Family: &goparent.Family{
					ID:          "1",
					Admin:       "1",
					Members:     []string{"1"},
					CreatedAt:   time.Now(),
					LastUpdated: time.Now(),
				},
			},

			sleepService: &mock.SleepService{
				Sleeps: []*goparent.Sleep{},
			},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			responseCode: http.StatusOK,
			resultLength: 0,
		},
		{
			desc: "returns one sleep",
			env:  &goparent.Env{},
			userService: &mock.UserService{
				Family: &goparent.Family{
					ID:          "1",
					Admin:       "1",
					Members:     []string{"1"},
					CreatedAt:   time.Now(),
					LastUpdated: time.Now(),
				},
			},

			sleepService: &mock.SleepService{
				Sleeps: []*goparent.Sleep{
					&goparent.Sleep{ID: "1"},
				},
			},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			responseCode: http.StatusOK,
			resultLength: 1,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockHandler := Handler{
				Env:           tC.env,
				UserService:   tC.userService,
				FamilyService: tC.familyService,
				SleepService:  tC.sleepService,
			}

			req, err := http.NewRequest("GET", "/sleep", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := mockHandler.sleepGetHandler()
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
			if tC.responseCode == http.StatusOK {
				var result SleepResponse
				decoder := json.NewDecoder(rr.Body)
				err := decoder.Decode(&result)
				assert.Nil(t, err)
				assert.Equal(t, tC.resultLength, len(result.SleepData))
			}
		})
	}
}

func TestSleepNewHandler(t *testing.T) {
	timestamp := time.Now()
	testCases := []struct {
		desc          string
		env           *goparent.Env
		sleepRequest  SleepRequest
		userService   goparent.UserService
		familyService goparent.FamilyService
		sleepService  goparent.SleepService
		contextUser   *goparent.User
		contextError  bool
		responseCode  int
		resultLength  int
	}{
		{
			desc: "submit sleep",
			env:  &goparent.Env{},
			sleepRequest: SleepRequest{
				SleepData: goparent.Sleep{
					Start:    timestamp,
					End:      timestamp.Add(time.Hour),
					UserID:   "1",
					ChildID:  "1",
					FamilyID: "1"}},
			userService: &mock.UserService{
				Family: &goparent.Family{
					ID:          "1",
					Admin:       "1",
					Members:     []string{"1"},
					CreatedAt:   time.Now(),
					LastUpdated: time.Now(),
				},
			},
			familyService: &mock.FamilyService{},
			sleepService:  &mock.SleepService{},
			contextUser:   &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			contextError:  false,
			responseCode:  http.StatusOK,
			resultLength:  0,
		},
		{
			desc: "returns no family error",
			env:  &goparent.Env{},
			sleepRequest: SleepRequest{
				SleepData: goparent.Sleep{
					Start:    timestamp,
					End:      timestamp.Add(time.Hour),
					UserID:   "1",
					ChildID:  "1",
					FamilyID: "1"}},
			userService: &mock.UserService{
				FamilyErr: errors.New("user has no current family"),
			},
			familyService: &mock.FamilyService{},
			sleepService:  &mock.SleepService{},
			contextUser:   &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			contextError:  false,
			responseCode:  http.StatusInternalServerError,
		},
		{
			desc: "returns sleep error",
			env:  &goparent.Env{},
			sleepRequest: SleepRequest{
				SleepData: goparent.Sleep{
					Start:    timestamp,
					End:      timestamp.Add(time.Hour),
					UserID:   "1",
					ChildID:  "1",
					FamilyID: "1"}},
			userService: &mock.UserService{
				Family: &goparent.Family{
					ID:          "1",
					Admin:       "1",
					Members:     []string{"1"},
					CreatedAt:   time.Now(),
					LastUpdated: time.Now(),
				},
			},
			familyService: &mock.FamilyService{},
			sleepService: &mock.SleepService{
				GetErr: errors.New("unknown sleep error"),
			},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			contextError: false,
			responseCode: http.StatusConflict,
		},
		{
			desc: "returns auth error",
			env:  &goparent.Env{},
			sleepRequest: SleepRequest{
				SleepData: goparent.Sleep{
					Start:    timestamp,
					End:      timestamp.Add(time.Hour),
					UserID:   "1",
					ChildID:  "1",
					FamilyID: "1"}},
			userService:   &mock.UserService{},
			familyService: &mock.FamilyService{},
			sleepService:  &mock.SleepService{},
			contextUser:   &goparent.User{},
			contextError:  true,
			responseCode:  http.StatusUnauthorized,
		},
		{
			desc: "decode input error",
			env:  &goparent.Env{},
			userService: &mock.UserService{
				Family: &goparent.Family{
					ID:          "1",
					Admin:       "1",
					Members:     []string{"1"},
					CreatedAt:   time.Now(),
					LastUpdated: time.Now(),
				},
			},
			familyService: &mock.FamilyService{},
			sleepService:  &mock.SleepService{},
			contextUser:   &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			contextError:  false,
			responseCode:  http.StatusInternalServerError,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockHandler := Handler{
				Env:           tC.env,
				UserService:   tC.userService,
				FamilyService: tC.familyService,
				SleepService:  tC.sleepService,
			}

			js, err := json.Marshal(&tC.sleepRequest)
			if err != nil {
				t.Fatal(err)
			}
			if (tC.sleepRequest == SleepRequest{}) {
				js = []byte("this is a test")
			}
			req, err := http.NewRequest("POST", "/sleep", bytes.NewReader(js))
			if err != nil {
				t.Fatal(err)
			}

			handler := mockHandler.sleepNewHandler()
			rr := httptest.NewRecorder()
			ctx := req.Context()
			if tC.contextError == true {
				ctx = context.WithValue(ctx, userContextKey, "")
			} else {
				ctx = context.WithValue(ctx, userContextKey, tC.contextUser)
			}
			req = req.WithContext(ctx)
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tC.responseCode, rr.Code)
			// if tC.responseCode == http.StatusOK {

			// }
		})
	}
}

func TestSleepViewHandler(t *testing.T) {
	testCases := []struct {
		desc         string
		env          *goparent.Env
		route        string
		method       string
		responseCode int
		contextUser  *goparent.User
	}{
		{
			desc:         "sleepViewHandler unauthorized",
			env:          &goparent.Env{},
			route:        "/sleep/1",
			method:       "GET",
			responseCode: http.StatusUnauthorized,
			contextUser:  nil,
		},
		{
			desc:         "sleepViewHandler not impl",
			env:          &goparent.Env{},
			route:        "/sleep/1",
			method:       "GET",
			responseCode: http.StatusNotImplemented,
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockHandler := Handler{
				Env: tC.env,
			}
			req, err := http.NewRequest(tC.method, tC.route, nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := mockHandler.sleepViewHandler()
			rr := httptest.NewRecorder()
			ctx := req.Context()
			if tC.contextUser == nil {
				ctx = context.WithValue(ctx, userContextKey, "")
			} else {
				ctx = context.WithValue(ctx, userContextKey, tC.contextUser)
			}
			req = req.WithContext(ctx)
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tC.responseCode, rr.Code)

		})
	}
}

func TestSleepEditHandler(t *testing.T) {
	testCases := []struct {
		desc         string
		env          *goparent.Env
		route        string
		method       string
		responseCode int
		contextUser  *goparent.User
	}{
		{
			desc:         "sleepEditHandler unauthorized",
			env:          &goparent.Env{},
			route:        "/sleep/1",
			method:       "PUT",
			responseCode: http.StatusUnauthorized,
			contextUser:  nil,
		},
		{
			desc:         "sleepEditHandler not impl",
			env:          &goparent.Env{},
			route:        "/sleep/1",
			method:       "PUT",
			responseCode: http.StatusNotImplemented,
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockHandler := Handler{
				Env: tC.env,
			}
			req, err := http.NewRequest(tC.method, tC.route, nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := mockHandler.sleepEditHandler()
			rr := httptest.NewRecorder()
			ctx := req.Context()
			if tC.contextUser == nil {
				ctx = context.WithValue(ctx, userContextKey, "")
			} else {
				ctx = context.WithValue(ctx, userContextKey, tC.contextUser)
			}
			req = req.WithContext(ctx)
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tC.responseCode, rr.Code)

		})
	}
}

func TestSleepDeleteHandler(t *testing.T) {
	testCases := []struct {
		desc         string
		env          *goparent.Env
		route        string
		method       string
		responseCode int
		contextUser  *goparent.User
	}{
		{
			desc:         "sleepDeleteHandler unauthorized",
			env:          &goparent.Env{},
			route:        "/sleep/1",
			method:       "DELETE",
			responseCode: http.StatusUnauthorized,
			contextUser:  nil,
		},
		{
			desc:         "sleepDeleteHandler not impl",
			env:          &goparent.Env{},
			route:        "/sleep/1",
			method:       "DELETE",
			responseCode: http.StatusNotImplemented,
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockHandler := Handler{
				Env: tC.env,
			}
			req, err := http.NewRequest(tC.method, tC.route, nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := mockHandler.sleepDeleteHandler()
			rr := httptest.NewRecorder()
			ctx := req.Context()
			if tC.contextUser == nil {
				ctx = context.WithValue(ctx, userContextKey, "")
			} else {
				ctx = context.WithValue(ctx, userContextKey, tC.contextUser)
			}
			req = req.WithContext(ctx)
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tC.responseCode, rr.Code)

		})
	}
}

func TestSleepStartHandler(t *testing.T) {
	testCases := []struct {
		desc         string
		env          *goparent.Env
		route        string
		method       string
		responseCode int
		contextUser  *goparent.User
	}{
		{
			desc:         "sleepStartHandler unauthorized",
			env:          &goparent.Env{},
			route:        "/sleep/start",
			method:       "POST",
			responseCode: http.StatusUnauthorized,
			contextUser:  nil,
		},
		{
			desc:         "sleepStartHandler not impl",
			env:          &goparent.Env{},
			route:        "/sleep/start",
			method:       "POST",
			responseCode: http.StatusNotImplemented,
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockHandler := Handler{
				Env: tC.env,
			}
			req, err := http.NewRequest(tC.method, tC.route, nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := mockHandler.sleepStartHandler()
			rr := httptest.NewRecorder()
			ctx := req.Context()
			if tC.contextUser == nil {
				ctx = context.WithValue(ctx, userContextKey, "")
			} else {
				ctx = context.WithValue(ctx, userContextKey, tC.contextUser)
			}
			req = req.WithContext(ctx)
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tC.responseCode, rr.Code)

		})
	}
}

func TestSleepEndHandler(t *testing.T) {
	testCases := []struct {
		desc         string
		env          *goparent.Env
		route        string
		method       string
		responseCode int
		contextUser  *goparent.User
	}{
		{
			desc:         "sleepStartHandler unauthorized",
			env:          &goparent.Env{},
			route:        "/sleep/end",
			method:       "POST",
			responseCode: http.StatusUnauthorized,
			contextUser:  nil,
		},
		{
			desc:         "sleepStartHandler not impl",
			env:          &goparent.Env{},
			route:        "/sleep/end",
			method:       "POST",
			responseCode: http.StatusNotImplemented,
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockHandler := Handler{
				Env: tC.env,
			}
			req, err := http.NewRequest(tC.method, tC.route, nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := mockHandler.sleepEndHandler()
			rr := httptest.NewRecorder()
			ctx := req.Context()
			if tC.contextUser == nil {
				ctx = context.WithValue(ctx, userContextKey, "")
			} else {
				ctx = context.WithValue(ctx, userContextKey, tC.contextUser)
			}
			req = req.WithContext(ctx)
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tC.responseCode, rr.Code)

		})
	}
}

func TestSleepToggleStatusHandler(t *testing.T) {
	testCases := []struct {
		desc         string
		env          *goparent.Env
		route        string
		method       string
		responseCode int
		contextUser  *goparent.User
	}{
		{
			desc:         "sleepToggleStatusHandler unauthorized",
			env:          &goparent.Env{},
			route:        "/sleep/status",
			method:       "GET",
			responseCode: http.StatusUnauthorized,
			contextUser:  nil,
		},
		{
			desc:         "sleepToggleStatusHandler not impl",
			env:          &goparent.Env{},
			route:        "/sleep/status",
			method:       "GET",
			responseCode: http.StatusNotImplemented,
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockHandler := Handler{
				Env: tC.env,
			}
			req, err := http.NewRequest(tC.method, tC.route, nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := mockHandler.sleepToggleStatus()
			rr := httptest.NewRecorder()
			ctx := req.Context()
			if tC.contextUser == nil {
				ctx = context.WithValue(ctx, userContextKey, "")
			} else {
				ctx = context.WithValue(ctx, userContextKey, tC.contextUser)
			}
			req = req.WithContext(ctx)
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tC.responseCode, rr.Code)

		})
	}
}

func TestInitSleepHandlers(t *testing.T) {
	testCases := []struct {
		desc    string
		name    string
		path    string
		methods []string
	}{
		{
			desc:    "sleep get",
			name:    "SleepGet",
			path:    "/sleep",
			methods: []string{"GET"},
		},
		{
			desc:    "sleep new",
			name:    "SleepNew",
			path:    "/sleep",
			methods: []string{"POST"},
		},
		{
			desc:    "sleep view",
			name:    "SleepView",
			path:    "/sleep/{id}",
			methods: []string{"GET"},
		},
		{
			desc:    "sleep edit",
			name:    "SleepEdit",
			path:    "/sleep/{id}",
			methods: []string{"PUT"},
		},
		{
			desc:    "sleep delete",
			name:    "SleepDelete",
			path:    "/sleep/{id}",
			methods: []string{"DELETE"},
		},
		{
			desc:    "sleep status",
			name:    "SleepStatus",
			path:    "/sleep/status",
			methods: []string{"GET"},
		},
		{
			desc:    "sleep start",
			name:    "SleepStart",
			path:    "/sleep/start",
			methods: []string{"POST"},
		},
		{
			desc:    "sleep end",
			name:    "SleepEnd",
			path:    "/sleep/end",
			methods: []string{"POST"},
		},
	}

	var testEnv goparent.Env
	h := Handler{Env: &testEnv}
	routes := mux.NewRouter()
	h.initSleepHandlers(routes)

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
