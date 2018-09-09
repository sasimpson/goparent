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

func TestFeedingGetHandler(t *testing.T) {
	testCases := []struct {
		desc           string
		env            *goparent.Env
		userService    goparent.UserService
		familyService  goparent.FamilyService
		feedingService goparent.FeedingService
		contextUser    *goparent.User
		contextError   bool
		responseCode   int
		resultLength   int
	}{
		{
			desc: "returns no feedings",
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
			feedingService: &mock.FeedingService{
				Feedings: []*goparent.Feeding{},
			},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			contextError: false,
			responseCode: http.StatusOK,
			resultLength: 0,
		},
		{
			desc: "returns some feedings",
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
			feedingService: &mock.FeedingService{
				Feedings: []*goparent.Feeding{
					&goparent.Feeding{ID: "1", Type: "bottle", Amount: 4.5, UserID: "1", FamilyID: "1", TimeStamp: time.Now(), ChildID: "1"},
					&goparent.Feeding{ID: "2", Type: "bottle", Amount: 5.5, UserID: "1", FamilyID: "1", TimeStamp: time.Now().Add(time.Hour), ChildID: "1"},
				},
			},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			contextError: false,
			responseCode: http.StatusOK,
			resultLength: 2,
		},
		{
			desc: "returns no family error",
			env:  &goparent.Env{},
			userService: &mock.UserService{
				FamilyErr: errors.New("user has no current family"),
			},
			familyService:  &mock.FamilyService{},
			feedingService: &mock.FeedingService{},
			contextUser:    &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			contextError:   false,
			responseCode:   http.StatusInternalServerError,
		},
		{
			desc:          "returns feeding error",
			env:           &goparent.Env{},
			userService:   &mock.UserService{},
			familyService: &mock.FamilyService{},
			feedingService: &mock.FeedingService{
				GetErr: errors.New("unknown feeding error"),
			},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			contextError: false,
			responseCode: http.StatusInternalServerError,
		},
		{
			desc:           "returns auth error",
			env:            &goparent.Env{},
			userService:    &mock.UserService{},
			familyService:  &mock.FamilyService{},
			feedingService: &mock.FeedingService{},
			contextUser:    &goparent.User{},
			contextError:   true,
			responseCode:   http.StatusUnauthorized,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockHandler := Handler{
				Env:            tC.env,
				UserService:    tC.userService,
				FamilyService:  tC.familyService,
				FeedingService: tC.feedingService,
			}

			req, err := http.NewRequest("GET", "/feeding", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := mockHandler.feedingGetHandler()
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
			if tC.responseCode == http.StatusOK {
				var result FeedingResponse
				decoder := json.NewDecoder(rr.Body)
				err := decoder.Decode(&result)
				assert.Nil(t, err)
				assert.Equal(t, tC.resultLength, len(result.FeedingData))
			}
		})
	}
}

func TestFeedingNewHandler(t *testing.T) {
	timestamp := time.Now()
	testCases := []struct {
		desc           string
		env            *goparent.Env
		feedingRequest FeedingRequest
		userService    goparent.UserService
		familyService  goparent.FamilyService
		feedingService goparent.FeedingService
		contextUser    *goparent.User
		contextError   bool
		responseCode   int
		resultLength   int
	}{
		{
			desc: "submit feeding",
			env:  &goparent.Env{},
			feedingRequest: FeedingRequest{
				FeedingData: goparent.Feeding{
					Type:      "bottle",
					Amount:    3.5,
					Side:      "",
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
			familyService:  &mock.FamilyService{},
			feedingService: &mock.FeedingService{},
			contextUser:    &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			contextError:   false,
			responseCode:   http.StatusOK,
			resultLength:   0,
		},
		{
			desc: "returns no family error",
			env:  &goparent.Env{},
			feedingRequest: FeedingRequest{
				FeedingData: goparent.Feeding{
					Type:      "bottle",
					Amount:    3.5,
					Side:      "",
					UserID:    "1",
					ChildID:   "1",
					FamilyID:  "1",
					TimeStamp: timestamp}},
			userService: &mock.UserService{
				FamilyErr: errors.New("user has no current family"),
			},
			familyService:  &mock.FamilyService{},
			feedingService: &mock.FeedingService{},
			contextUser:    &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			contextError:   false,
			responseCode:   http.StatusInternalServerError,
		},
		{
			desc: "returns feeding error",
			env:  &goparent.Env{},
			feedingRequest: FeedingRequest{
				FeedingData: goparent.Feeding{
					Type:      "bottle",
					Amount:    3.5,
					Side:      "",
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
			feedingService: &mock.FeedingService{
				GetErr: errors.New("unknown feeding error"),
			},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			contextError: false,
			responseCode: http.StatusConflict,
		},
		{
			desc: "returns auth error",
			env:  &goparent.Env{},
			feedingRequest: FeedingRequest{
				FeedingData: goparent.Feeding{
					Type:      "bottle",
					Amount:    3.5,
					Side:      "",
					UserID:    "1",
					ChildID:   "1",
					FamilyID:  "1",
					TimeStamp: timestamp}},
			userService:    &mock.UserService{},
			familyService:  &mock.FamilyService{},
			feedingService: &mock.FeedingService{},
			contextUser:    &goparent.User{},
			contextError:   true,
			responseCode:   http.StatusUnauthorized,
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
			familyService:  &mock.FamilyService{},
			feedingService: &mock.FeedingService{},
			contextUser:    &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			contextError:   false,
			responseCode:   http.StatusInternalServerError,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockHandler := Handler{
				Env:            tC.env,
				UserService:    tC.userService,
				FamilyService:  tC.familyService,
				FeedingService: tC.feedingService,
			}

			js, err := json.Marshal(&tC.feedingRequest)
			if err != nil {
				t.Fatal(err)
			}
			if (tC.feedingRequest == FeedingRequest{}) {
				js = []byte("this is a test")
			}
			req, err := http.NewRequest("POST", "/feeding", bytes.NewReader(js))
			if err != nil {
				t.Fatal(err)
			}

			handler := mockHandler.feedingNewHandler()
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

func TestFeedingViewHandler(t *testing.T) {

	mockHandler := Handler{
		Env: &goparent.Env{},
		UserService: &mock.UserService{
			Family: &goparent.Family{
				ID:          "1",
				Admin:       "1",
				Members:     []string{"1"},
				CreatedAt:   time.Now(),
				LastUpdated: time.Now(),
			},
		},
		FamilyService:  &mock.FamilyService{},
		FeedingService: &mock.FeedingService{},
	}

	req, err := http.NewRequest("GET", "/feeding/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := mockHandler.feedingViewHandler()

	rr := httptest.NewRecorder()
	ctx := req.Context()
	ctx = context.WithValue(ctx, userContextKey, "")
	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	rr = httptest.NewRecorder()
	ctx = req.Context()
	ctx = context.WithValue(ctx, userContextKey, &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotImplemented, rr.Code)
}

func TestFeedingEditHandler(t *testing.T) {
	mockHandler := Handler{
		Env: &goparent.Env{},
		UserService: &mock.UserService{
			Family: &goparent.Family{
				ID:          "1",
				Admin:       "1",
				Members:     []string{"1"},
				CreatedAt:   time.Now(),
				LastUpdated: time.Now(),
			},
		},
		FamilyService:  &mock.FamilyService{},
		FeedingService: &mock.FeedingService{},
	}

	req, err := http.NewRequest("GET", "/feeding/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := mockHandler.feedingEditHandler()
	rr := httptest.NewRecorder()
	ctx := req.Context()
	ctx = context.WithValue(ctx, userContextKey, "")
	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	rr = httptest.NewRecorder()
	ctx = req.Context()
	ctx = context.WithValue(ctx, userContextKey, &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotImplemented, rr.Code)
}

func TestFeedingDeleteHandler(t *testing.T) {
	mockHandler := Handler{
		Env: &goparent.Env{},
		UserService: &mock.UserService{
			Family: &goparent.Family{
				ID:          "1",
				Admin:       "1",
				Members:     []string{"1"},
				CreatedAt:   time.Now(),
				LastUpdated: time.Now(),
			},
		},
		FamilyService:  &mock.FamilyService{},
		FeedingService: &mock.FeedingService{},
	}

	req, err := http.NewRequest("DELETE", "/feeding/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := mockHandler.feedingDeleteHandler()
	rr := httptest.NewRecorder()
	ctx := req.Context()
	ctx = context.WithValue(ctx, userContextKey, "")
	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	rr = httptest.NewRecorder()
	ctx = req.Context()
	ctx = context.WithValue(ctx, userContextKey, &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotImplemented, rr.Code)
}

func TestInitFeedingHandlers(t *testing.T) {
	testCases := []struct {
		desc    string
		name    string
		path    string
		methods []string
	}{
		{
			desc:    "feeding get",
			name:    "FeedingGet",
			path:    "/feeding",
			methods: []string{"GET"},
		},
		{
			desc:    "feeding new",
			name:    "FeedingNew",
			path:    "/feeding",
			methods: []string{"POST"},
		},
		{
			desc:    "feeding view",
			name:    "FeedingView",
			path:    "/feeding/{id}",
			methods: []string{"GET"},
		},
		{
			desc:    "feeding edit",
			name:    "FeedingEdit",
			path:    "/feeding/{id}",
			methods: []string{"PUT"},
		},
		{
			desc:    "feeding delete",
			name:    "FeedingDelete",
			path:    "/feeding/{id}",
			methods: []string{"DELETE"},
		},
	}

	var testEnv *goparent.Env
	h := Handler{Env: testEnv}
	routes := mux.NewRouter()
	h.initFeedingHandlers(routes)

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
