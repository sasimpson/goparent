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

func TestSleepNewHandler(t *testing.T) {
	var testEnv config.Env

	mock := r.NewMock()
	mock.On(
		r.Table("sleep").MockAnything(),
	).Return(r.WriteResponse{
		Inserted:      1,
		Errors:        0,
		GeneratedKeys: []string{"1"},
	}, nil)
	testEnv.DB.Session = mock

	s := SleepRequest{SleepData: models.Sleep{UserID: "1", SleepStart: time.Now().AddDate(0, 0, -1), SleepEnd: time.Now()}}
	js, err := json.Marshal(&s)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/feeding", bytes.NewReader(js))
	if err != nil {
		t.Fatal(err)
	}

	handler := sleepNewHandler(&testEnv)
	rr := httptest.NewRecorder()
	ctx := req.Context()
	ctx = context.WithValue(ctx, "user", models.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestSleepGetHandler(t *testing.T) {
	var testEnv config.Env

	mock := r.NewMock()
	mock.On(
		r.Table("sleep").MockAnything(),
	).Return(
		map[string]interface{}{
			"start":  time.Now().AddDate(0, 0, -1),
			"end":    time.Now(),
			"userid": "1",
		}, nil,
	)
	testEnv.DB.Session = mock

	req, err := http.NewRequest("GET", "/sleep", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := sleepGetHandler(&testEnv)
	rr := httptest.NewRecorder()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "user", models.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestSleepStatusHandler(t *testing.T) {
	testCases := []struct {
		desc   string
		status int
		ret    map[string]interface{}
		err    error
	}{
		{
			desc:   "false",
			status: http.StatusNotFound,
			ret:    nil,
			err:    r.ErrEmptyResult,
		},
		{
			desc:   "true",
			status: http.StatusOK,
			ret: map[string]interface{}{
				"start":  time.Now().AddDate(0, 0, -1),
				"userid": "1",
				"id":     "1",
			},
			err: nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			var testEnv config.Env

			mock := r.NewMock()
			mock.On(
				r.Table("sleep").MockAnything(),
			).Return(
				tC.ret, tC.err,
			)
			testEnv.DB.Session = mock

			req, err := http.NewRequest("GET", "/sleep/status", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := sleepToggleStatus(&testEnv)
			rr := httptest.NewRecorder()

			ctx := req.Context()
			ctx = context.WithValue(ctx, "user", models.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
			req = req.WithContext(ctx)
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tC.status, rr.Code)
		})
	}
}

func TestSleepStartHandler(t *testing.T) {
	testCases := []struct {
		desc   string
		status int
		ret    map[string]interface{}
		err    error
	}{
		{
			desc:   "false, no sleep active",
			status: http.StatusOK,
			ret:    nil,
			err:    r.ErrEmptyResult,
		},
		{
			desc:   "true, sleep active",
			status: http.StatusConflict,
			ret: map[string]interface{}{
				"start":  time.Now().AddDate(0, 0, -1),
				"userid": "1",
				"id":     "1",
			},
			err: nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			var testEnv config.Env

			mock := r.NewMock()
			mock.On(
				r.Table("sleep").MockAnything(),
			).Return(
				tC.ret, tC.err,
			)
			testEnv.DB.Session = mock

			req, err := http.NewRequest("GET", "/sleep/start", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := sleepStartHandler(&testEnv)
			rr := httptest.NewRecorder()

			ctx := req.Context()
			ctx = context.WithValue(ctx, "user", models.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
			req = req.WithContext(ctx)
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tC.status, rr.Code)
		})
	}
}

func TestSleepEndHandler(t *testing.T) {
	testCases := []struct {
		desc   string
		status int
		ret    map[string]interface{}
		err    error
	}{
		{
			desc:   "false, no sleep active",
			status: http.StatusOK,
			ret:    nil,
			err:    r.ErrEmptyResult,
		},
		{
			desc:   "true, sleep active",
			status: http.StatusConflict,
			ret: map[string]interface{}{
				"start":  time.Now().AddDate(0, 0, -1),
				"userid": "1",
				"id":     "1",
			},
			err: nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			var testEnv config.Env

			mock := r.NewMock()
			mock.On(
				r.Table("sleep").MockAnything(),
			).Return(
				tC.ret, tC.err,
			)
			testEnv.DB.Session = mock

			req, err := http.NewRequest("GET", "/sleep/end", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := sleepStartHandler(&testEnv)
			rr := httptest.NewRecorder()

			ctx := req.Context()
			ctx = context.WithValue(ctx, "user", models.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
			req = req.WithContext(ctx)
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tC.status, rr.Code)
		})
	}
}

func TestSleepViewHandler(t *testing.T) {
	var testEnv config.Env

	req, err := http.NewRequest("GET", "/sleep/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := sleepViewHandler(&testEnv)
	rr := httptest.NewRecorder()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "user", models.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestSleepEditHandler(t *testing.T) {
	var testEnv config.Env

	req, err := http.NewRequest("GET", "/sleep/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := sleepEditHandler(&testEnv)
	rr := httptest.NewRecorder()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "user", models.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestSleepDeleteHandler(t *testing.T) {
	var testEnv config.Env

	req, err := http.NewRequest("GET", "/sleep/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := sleepDeleteHandler(&testEnv)
	rr := httptest.NewRecorder()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "user", models.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
	req = req.WithContext(ctx)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
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

	var testEnv config.Env
	routes := mux.NewRouter()
	initSleepHandlers(&testEnv, routes)

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
