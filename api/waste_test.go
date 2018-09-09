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

	"github.com/gorilla/mux"
	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/mock"
	"github.com/stretchr/testify/assert"
)

func TestWasteGetHandler(t *testing.T) {
	testCases := []struct {
		desc          string
		env           *goparent.Env
		userService   goparent.UserService
		familyService goparent.FamilyService
		wasteService  goparent.WasteService
		contextUser   *goparent.User
		responseCode  int
		resultLength  int
	}{
		{
			desc:          "returns auth error",
			env:           &goparent.Env{},
			userService:   &mock.UserService{},
			familyService: &mock.FamilyService{},
			wasteService:  &mock.WasteService{},
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
			wasteService:  &mock.WasteService{},
			contextUser:   &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			responseCode:  http.StatusInternalServerError,
		},
		{
			desc: "returns waste error",
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

			wasteService: &mock.WasteService{
				GetErr: errors.New("test error"),
			},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			responseCode: http.StatusInternalServerError,
		},
		{
			desc: "returns no waste",
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

			wasteService: &mock.WasteService{
				Wastes: []*goparent.Waste{},
			},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			responseCode: http.StatusOK,
			resultLength: 0,
		},
		{
			desc: "returns one waste",
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

			wasteService: &mock.WasteService{
				Wastes: []*goparent.Waste{
					&goparent.Waste{ID: "1"},
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
				WasteService:  tC.wasteService,
			}

			req, err := http.NewRequest("GET", "/waste", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := mockHandler.wasteGetHandler()
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
				var result WasteResponse
				decoder := json.NewDecoder(rr.Body)
				err := decoder.Decode(&result)
				assert.Nil(t, err)
				assert.Equal(t, tC.resultLength, len(result.WasteData))
			}
		})
	}
}

func TestWasteViewHandler(t *testing.T) {
	testCases := []struct {
		desc         string
		env          *goparent.Env
		route        string
		method       string
		responseCode int
		contextUser  *goparent.User
	}{
		{
			desc:         "wasteViewHandler unauthorized",
			env:          &goparent.Env{},
			route:        "/waste/1",
			method:       "GET",
			responseCode: http.StatusUnauthorized,
			contextUser:  nil,
		},
		{
			desc:         "wasteViewHandler not impl",
			env:          &goparent.Env{},
			route:        "/waste/1",
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

			handler := mockHandler.wasteViewHandler()
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

func TestWasteEditHandler(t *testing.T) {
	testCases := []struct {
		desc         string
		env          *goparent.Env
		route        string
		method       string
		responseCode int
		contextUser  *goparent.User
	}{
		{
			desc:         "wasteEditHandler unauthorized",
			env:          &goparent.Env{},
			route:        "/waste/1",
			method:       "PUT",
			responseCode: http.StatusUnauthorized,
			contextUser:  nil,
		},
		{
			desc:         "wasteEditHandler not impl",
			env:          &goparent.Env{},
			route:        "/waste/1",
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

			handler := mockHandler.wasteEditHandler()
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

func TestWasteDeleteHandler(t *testing.T) {
	testCases := []struct {
		desc         string
		env          *goparent.Env
		route        string
		method       string
		responseCode int
		contextUser  *goparent.User
	}{
		{
			desc:         "wasteDeleteHandler unauthorized",
			env:          &goparent.Env{},
			route:        "/waste/1",
			method:       "DELETE",
			responseCode: http.StatusUnauthorized,
			contextUser:  nil,
		},
		{
			desc:         "wasteDeleteHandler not impl",
			env:          &goparent.Env{},
			route:        "/waste/1",
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

			handler := mockHandler.wasteDeleteHandler()
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

func TestWasteNewHandler(t *testing.T) {
	timestamp := time.Now()
	testCases := []struct {
		desc          string
		env           *goparent.Env
		wasteRequest  WasteRequest
		userService   goparent.UserService
		familyService goparent.FamilyService
		wasteService  goparent.WasteService
		contextUser   *goparent.User
		contextError  bool
		responseCode  int
		resultLength  int
	}{
		{
			desc: "submit waste",
			env:  &goparent.Env{},
			wasteRequest: WasteRequest{
				WasteData: goparent.Waste{
					Type:      1,
					Notes:     "",
					UserID:    "1",
					ChildID:   "1",
					FamilyID:  "1",
					TimeStamp: timestamp}},
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
			wasteService:  &mock.WasteService{},
			contextUser:   &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			contextError:  false,
			responseCode:  http.StatusOK,
			resultLength:  0,
		},
		{
			desc: "returns no family error",
			env:  &goparent.Env{},
			wasteRequest: WasteRequest{
				WasteData: goparent.Waste{
					Type:      1,
					Notes:     "",
					UserID:    "1",
					ChildID:   "1",
					FamilyID:  "1",
					TimeStamp: timestamp}},
			userService: &mock.UserService{
				FamilyErr: errors.New("user has no current family"),
			},
			familyService: &mock.FamilyService{},
			wasteService:  &mock.WasteService{},
			contextUser:   &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			contextError:  false,
			responseCode:  http.StatusInternalServerError,
		},
		{
			desc: "returns waste error",
			env:  &goparent.Env{},
			wasteRequest: WasteRequest{
				WasteData: goparent.Waste{
					Type:      1,
					Notes:     "",
					UserID:    "1",
					ChildID:   "1",
					FamilyID:  "1",
					TimeStamp: timestamp}},
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
			wasteService: &mock.WasteService{
				GetErr: errors.New("unknown waste error"),
			},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			contextError: false,
			responseCode: http.StatusConflict,
		},
		{
			desc: "returns auth error",
			env:  &goparent.Env{},
			wasteRequest: WasteRequest{
				WasteData: goparent.Waste{
					Type:      1,
					Notes:     "",
					UserID:    "1",
					ChildID:   "1",
					FamilyID:  "1",
					TimeStamp: timestamp}},
			userService:   &mock.UserService{},
			familyService: &mock.FamilyService{},
			wasteService:  &mock.WasteService{},
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
			wasteService:  &mock.WasteService{},
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
				WasteService:  tC.wasteService,
			}

			js, err := json.Marshal(&tC.wasteRequest)
			if err != nil {
				t.Fatal(err)
			}
			if (tC.wasteRequest == WasteRequest{}) {
				js = []byte("this is a test")
			}
			req, err := http.NewRequest("POST", "/waste", bytes.NewReader(js))
			if err != nil {
				t.Fatal(err)
			}

			handler := mockHandler.wasteNewHandler()
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

func TestInitWasteHandler(t *testing.T) {
	testCases := []struct {
		desc    string
		name    string
		path    string
		methods []string
	}{
		{
			desc:    "waste get",
			name:    "WasteGet",
			path:    "/waste",
			methods: []string{"GET"},
		},
		{
			desc:    "waste new",
			name:    "WasteNew",
			path:    "/waste",
			methods: []string{"POST"},
		},
		{
			desc:    "waste graph data",
			name:    "WasteGraphData",
			path:    "/waste/graph/{id}",
			methods: []string{"GET"},
		},
		{
			desc:    "waste view",
			name:    "WasteView",
			path:    "/waste/{id}",
			methods: []string{"GET"},
		},
		{
			desc:    "waste edit",
			name:    "WasteEdit",
			path:    "/waste/{id}",
			methods: []string{"PUT"},
		},
		{
			desc:    "waste delete",
			name:    "WasteDelete",
			path:    "/waste/{id}",
			methods: []string{"DELETE"},
		},
	}

	var testEnv goparent.Env
	h := Handler{Env: &testEnv}
	routes := mux.NewRouter()
	h.initWasteHandlers(routes)

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
