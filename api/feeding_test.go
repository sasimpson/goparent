package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/sasimpson/goparent/config"
	"github.com/sasimpson/goparent/models"
	"github.com/stretchr/testify/assert"

	r "gopkg.in/gorethink/gorethink.v3"
)

func TestFeedingGetHandler(t *testing.T) {
	var testEnv config.Env

	//mock out the db stuff and return call
	mock := r.NewMock()
	mock.On(
		r.Table("feeding").Filter(map[string]interface{}{"userid": "1"}).OrderBy(r.Desc("timestamp")),
	).Return([]interface{}{
		map[string]interface{}{
			"id":            "1",
			"feedingType":   "bottle",
			"feedingAmount": 1,
			"feedingSide":   "",
			"userid":        "1",
			"timestamp":     time.Now()},
	}, nil)
	testEnv.DB.Session = mock

	//setup request
	req, err := http.NewRequest("GET", "/feeding", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := FeedingGetHandler(&testEnv)
	rr := httptest.NewRecorder()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "user", models.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestFeedingNewHandler(t *testing.T) {
	var testEnv config.Env
	timestamp := time.Now()

	mock := r.NewMock()
	mock.On(
		r.Table("feeding").MockAnything(),
	).Return([]interface{}{
		map[string]interface{}{"id": "1", "feedingType": "bottle", "feedingAmount": 1, "feedingSide": "", "userid": "1", "timestamp": time.Now()},
	}, nil)
	testEnv.DB.Session = mock

	f := FeedingRequest{FeedingData: models.Feeding{Type: "bottle", Amount: 3.5, Side: "", UserID: "1", TimeStamp: timestamp}}
	js, err := json.Marshal(&f)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/feeding", bytes.NewReader(js))
	if err != nil {
		t.Fatal(err)
	}

	handler := FeedingNewHandler(&testEnv)
	rr := httptest.NewRecorder()
	ctx := req.Context()
	ctx = context.WithValue(ctx, "user", models.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestFeedingViewHandler(t *testing.T) {
	var testEnv config.Env

	req, err := http.NewRequest("GET", "/feeding/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := FeedingViewHandler(&testEnv)
	rr := httptest.NewRecorder()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "user", models.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestFeedingEditHandler(t *testing.T) {
	var testEnv config.Env

	req, err := http.NewRequest("GET", "/feeding/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := FeedingEditHandler(&testEnv)
	rr := httptest.NewRecorder()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "user", models.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestFeedingDeleteHandler(t *testing.T) {
	var testEnv config.Env

	req, err := http.NewRequest("GET", "/feeding/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := FeedingDeleteHandler(&testEnv)
	rr := httptest.NewRecorder()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "user", models.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
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
	}

	var testEnv config.Env
	r := mux.NewRouter()
	initFeedingHandlers(&testEnv, r)

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			route := r.Get(tC.name)
			path, _ := route.GetPathTemplate()
			methods, _ := route.GetMethods()
			assert.Equal(t, tC.name, route.GetName())
			assert.Equal(t, tC.path, path)
			assert.Equal(t, tC.methods, methods)
		})
	}
}

func gorillaWalkFn(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	name := route.GetName()
	methods, _ := route.GetMethods()
	path, _ := route.GetPathTemplate()
	log.Printf("%s: %s (%v)", name, path, methods)
	return nil
}
