package api

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestWasteGetHandler(t *testing.T) {
	var testEnv config.Env

	//mock out the db stuff and return call
	mock := r.NewMock()
	mock.On(
		r.Table("waste").MockAnything(),
	).Return([]interface{}{
		map[string]interface{}{
			"id":        "1",
			"wasteType": 1,
			"notes":     "test note",
			"userid":    "1",
			"timestamp": time.Now(),
		},
	}, nil)
	testEnv.DB.Session = mock

	//setup request
	req, err := http.NewRequest("GET", "/waste", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := wasteGetHandler(&testEnv)
	rr := httptest.NewRecorder()

	ctx := req.Context()
	ctx = context.WithValue(ctx, userContextKey, models.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestWasteNewHandler(t *testing.T) {
	var testEnv config.Env
	timestamp := time.Now()

	mock := r.NewMock()
	mock.On(
		r.Table("waste").MockAnything(),
	).Return(r.WriteResponse{
		Inserted:      1,
		Errors:        0,
		GeneratedKeys: []string{"1"},
	}, nil)
	testEnv.DB.Session = mock

	w := WasteRequest{WasteData: models.Waste{Type: 1, Notes: "some notes", UserID: "1", TimeStamp: timestamp}}
	js, err := json.Marshal(&w)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/waste", bytes.NewReader(js))
	if err != nil {
		t.Fatal(err)
	}

	handler := wasteNewHandler(&testEnv)
	rr := httptest.NewRecorder()
	ctx := req.Context()
	ctx = context.WithValue(ctx, userContextKey, models.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestWasteViewHandler(t *testing.T) {
	var testEnv config.Env

	mock := r.NewMock()
	mock.On(
		r.Table("waste").MockAnything(),
	).Return(map[string]interface{}{
		"id":        "1",
		"type":      1,
		"notes":     "some notes",
		"userid":    "1",
		"timestamp": time.Now(),
	}, nil)
	testEnv.DB.Session = mock

	req, err := http.NewRequest("GET", "/waste/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := wasteViewHandler(&testEnv)
	rr := httptest.NewRecorder()

	ctx := req.Context()
	ctx = context.WithValue(ctx, userContextKey, models.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestWasteEditHandler(t *testing.T) {
	var testEnv config.Env

	req, err := http.NewRequest("GET", "/waste/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := wasteEditHandler(&testEnv)
	rr := httptest.NewRecorder()

	ctx := req.Context()
	ctx = context.WithValue(ctx, userContextKey, models.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestWasteDeleteHandler(t *testing.T) {
	var testEnv config.Env

	req, err := http.NewRequest("DELETE", "/waste/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := wasteDeleteHandler(&testEnv)
	rr := httptest.NewRecorder()

	ctx := req.Context()
	ctx = context.WithValue(ctx, userContextKey, models.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
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

	var testEnv config.Env
	routes := mux.NewRouter()
	initWasteHandlers(&testEnv, routes)

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
