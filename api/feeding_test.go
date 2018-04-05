package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/config"
	"github.com/sasimpson/goparent/mock"
	"github.com/stretchr/testify/assert"
)

func TestFeedingGetHandler(t *testing.T) {
	var testEnv *config.Env
	mockHandler := Handler{
		FeedingService: &mock.MockFeedingService{Env: testEnv},
		UserService: &mock.MockUserService{
			Env: testEnv,
			Family: &goparent.Family{
				ID:          "1",
				Admin:       "1",
				Members:     []string{"1"},
				CreatedAt:   time.Now(),
				LastUpdated: time.Now(),
			},
		},
		FamilyService: &mock.MockFamilyService{Env: testEnv},
	}

	//setup request
	req, err := http.NewRequest("GET", "/feeding", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := mockHandler.feedingGetHandler()
	rr := httptest.NewRecorder()

	ctx := req.Context()
	ctx = context.WithValue(ctx, userContextKey, goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

// func TestFeedingNewHandler(t *testing.T) {
// 	var testEnv config.Env
// 	timestamp := time.Now()

// 	mock := r.NewMock()
// 	mock.
// 		On(
// 			r.Table("family").Filter(
// 				func(row r.Term) r.Term {
// 					return row.Field("members").Contains("1")
// 				},
// 			),
// 		).
// 		Return(map[string]interface{}{
// 			"id":           "1",
// 			"admin":        "1",
// 			"members":      []string{"1"},
// 			"created_at":   time.Now(),
// 			"last_updated": time.Now(),
// 		}, nil).
// 		On(
// 			r.Table("feeding").MockAnything(),
// 		).
// 		Return(r.WriteResponse{
// 			Inserted:      1,
// 			Errors:        0,
// 			GeneratedKeys: []string{"1"},
// 		}, nil)
// 	testEnv.DB.Session = mock

// 	f := FeedingRequest{
// 		FeedingData: goparent.Feeding{
// 			Type:      "bottle",
// 			Amount:    3.5,
// 			Side:      "",
// 			UserID:    "1",
// 			ChildID:   "1",
// 			FamilyID:  "1",
// 			TimeStamp: timestamp}}
// 	js, err := json.Marshal(&f)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req, err := http.NewRequest("POST", "/feeding", bytes.NewReader(js))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	handler := feedingNewHandler(&testEnv)
// 	rr := httptest.NewRecorder()
// 	ctx := req.Context()
// 	ctx = context.WithValue(ctx, userContextKey, goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
// 	req = req.WithContext(ctx)
// 	handler.ServeHTTP(rr, req)
// 	t.Log(rr.Body)
// 	assert.Equal(t, http.StatusOK, rr.Code)
// }

// func TestFeedingViewHandler(t *testing.T) {
// 	var testEnv config.Env

// 	req, err := http.NewRequest("GET", "/feeding/1", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	handler := feedingViewHandler(&testEnv)
// 	rr := httptest.NewRecorder()

// 	ctx := req.Context()
// 	ctx = context.WithValue(ctx, userContextKey, goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
// 	req = req.WithContext(ctx)
// 	handler.ServeHTTP(rr, req)
// 	assert.Equal(t, http.StatusOK, rr.Code)
// }

// func TestFeedingEditHandler(t *testing.T) {
// 	var testEnv config.Env

// 	req, err := http.NewRequest("GET", "/feeding/1", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	handler := feedingEditHandler(&testEnv)
// 	rr := httptest.NewRecorder()

// 	ctx := req.Context()
// 	ctx = context.WithValue(ctx, userContextKey, goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
// 	req = req.WithContext(ctx)
// 	handler.ServeHTTP(rr, req)
// 	assert.Equal(t, http.StatusOK, rr.Code)
// }

// func TestFeedingDeleteHandler(t *testing.T) {
// 	var testEnv config.Env

// 	req, err := http.NewRequest("DELETE", "/feeding/1", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	handler := feedingDeleteHandler(&testEnv)
// 	rr := httptest.NewRecorder()

// 	ctx := req.Context()
// 	ctx = context.WithValue(ctx, userContextKey, goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
// 	req = req.WithContext(ctx)
// 	handler.ServeHTTP(rr, req)
// 	assert.Equal(t, http.StatusOK, rr.Code)
// }

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

	var testEnv *config.Env
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
