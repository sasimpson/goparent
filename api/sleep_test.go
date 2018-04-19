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
	"github.com/sasimpson/goparent/config"
	"github.com/sasimpson/goparent/mock"
	"github.com/stretchr/testify/assert"
)

func TestSleepGetHandler(t *testing.T) {
	testCases := []struct {
		desc          string
		env           *config.Env
		userService   goparent.UserService
		familyService goparent.FamilyService
		sleepService  goparent.SleepService
		contextUser   *goparent.User
		responseCode  int
		resultLength  int
	}{
		{
			desc:          "returns auth error",
			env:           &config.Env{},
			userService:   &mock.MockUserService{},
			familyService: &mock.MockFamilyService{},
			sleepService:  &mock.MockSleepService{},
			contextUser:   nil,
			responseCode:  http.StatusUnauthorized,
		},
		{
			desc: "returns family error",
			env:  &config.Env{},
			userService: &mock.MockUserService{
				GetErr: errors.New("test error"),
			},
			familyService: &mock.MockFamilyService{},
			sleepService:  &mock.MockSleepService{},
			contextUser:   &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			responseCode:  http.StatusInternalServerError,
		},
		{
			desc: "returns sleep error",
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

			sleepService: &mock.MockSleepService{
				GetErr: errors.New("test error"),
			},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			responseCode: http.StatusInternalServerError,
		},
		{
			desc: "returns no sleep",
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

			sleepService: &mock.MockSleepService{
				Sleeps: []*goparent.Sleep{},
			},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			responseCode: http.StatusOK,
			resultLength: 0,
		},
		{
			desc: "returns one sleep",
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

			sleepService: &mock.MockSleepService{
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
		env           *config.Env
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
			env:  &config.Env{},
			sleepRequest: SleepRequest{
				SleepData: goparent.Sleep{
					Start:    timestamp,
					End:      timestamp.Add(time.Hour),
					UserID:   "1",
					ChildID:  "1",
					FamilyID: "1"}},
			userService: &mock.MockUserService{
				Family: &goparent.Family{
					ID:          "1",
					Admin:       "1",
					Members:     []string{"1"},
					CreatedAt:   time.Now(),
					LastUpdated: time.Now(),
				},
			},
			familyService: &mock.MockFamilyService{},
			sleepService:  &mock.MockSleepService{},
			contextUser:   &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			contextError:  false,
			responseCode:  http.StatusOK,
			resultLength:  0,
		},
		{
			desc: "returns no family error",
			env:  &config.Env{},
			sleepRequest: SleepRequest{
				SleepData: goparent.Sleep{
					Start:    timestamp,
					End:      timestamp.Add(time.Hour),
					UserID:   "1",
					ChildID:  "1",
					FamilyID: "1"}},
			userService: &mock.MockUserService{
				GetErr: errors.New("user has no current family"),
			},
			familyService: &mock.MockFamilyService{},
			sleepService:  &mock.MockSleepService{},
			contextUser:   &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			contextError:  false,
			responseCode:  http.StatusInternalServerError,
		},
		{
			desc: "returns sleep error",
			env:  &config.Env{},
			sleepRequest: SleepRequest{
				SleepData: goparent.Sleep{
					Start:    timestamp,
					End:      timestamp.Add(time.Hour),
					UserID:   "1",
					ChildID:  "1",
					FamilyID: "1"}},
			userService: &mock.MockUserService{
				Family: &goparent.Family{
					ID:          "1",
					Admin:       "1",
					Members:     []string{"1"},
					CreatedAt:   time.Now(),
					LastUpdated: time.Now(),
				},
			},
			familyService: &mock.MockFamilyService{},
			sleepService: &mock.MockSleepService{
				GetErr: errors.New("unknown sleep error"),
			},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			contextError: false,
			responseCode: http.StatusConflict,
		},
		{
			desc: "returns auth error",
			env:  &config.Env{},
			sleepRequest: SleepRequest{
				SleepData: goparent.Sleep{
					Start:    timestamp,
					End:      timestamp.Add(time.Hour),
					UserID:   "1",
					ChildID:  "1",
					FamilyID: "1"}},
			userService:   &mock.MockUserService{},
			familyService: &mock.MockFamilyService{},
			sleepService:  &mock.MockSleepService{},
			contextUser:   &goparent.User{},
			contextError:  true,
			responseCode:  http.StatusUnauthorized,
		},
		{
			desc: "decode input error",
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
			familyService: &mock.MockFamilyService{},
			sleepService:  &mock.MockSleepService{},
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

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	"github.com/gorilla/mux"
// 	"github.com/sasimpson/goparent"
// 	"github.com/sasimpson/goparent/config"
// 	"github.com/stretchr/testify/assert"

// 	r "gopkg.in/gorethink/gorethink.v3"
// )

// func TestSleepNewHandler(t *testing.T) {
// 	var testEnv config.Env

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
// 			r.Table("sleep").MockAnything(),
// 		).
// 		Return(r.WriteResponse{
// 			Inserted:      1,
// 			Errors:        0,
// 			GeneratedKeys: []string{"1"},
// 		}, nil)
// 	testEnv.DB.Session = mock

// 	s := SleepRequest{SleepData: goparent.Sleep{UserID: "1", ChildID: "1", Start: time.Now().AddDate(0, 0, -1), End: time.Now()}}
// 	js, err := json.Marshal(&s)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req, err := http.NewRequest("POST", "/sleep", bytes.NewReader(js))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	handler := sleepNewHandler(&testEnv)
// 	rr := httptest.NewRecorder()
// 	ctx := req.Context()
// 	ctx = context.WithValue(ctx, userContextKey, goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
// 	req = req.WithContext(ctx)
// 	handler.ServeHTTP(rr, req)
// 	assert.Equal(t, http.StatusOK, rr.Code)
// }

// func TestSleepStatusHandler(t *testing.T) {
// 	//TODO: verify returned values
// 	testCases := []struct {
// 		desc   string
// 		status int
// 		ret    map[string]interface{}
// 		err    error
// 	}{
// 		{
// 			desc:   "false",
// 			status: http.StatusNotFound,
// 			ret:    nil,
// 			err:    r.ErrEmptyResult,
// 		},
// 		{
// 			desc:   "true",
// 			status: http.StatusOK,
// 			ret: map[string]interface{}{
// 				"start":    time.Now().AddDate(0, 0, -1),
// 				"userID":   "1",
// 				"childID":  "1",
// 				"familyID": "1",
// 				"id":       "1",
// 			},
// 			err: nil,
// 		},
// 	}
// 	for _, tC := range testCases {
// 		t.Run(tC.desc, func(t *testing.T) {
// 			var testEnv config.Env

// 			mock := r.NewMock()
// 			mock.
// 				On(
// 					r.Table("family").Filter(
// 						func(row r.Term) r.Term {
// 							return row.Field("members").Contains("1")
// 						},
// 					),
// 				).
// 				Return(map[string]interface{}{
// 					"id":           "1",
// 					"admin":        "1",
// 					"members":      []string{"1"},
// 					"created_at":   time.Now(),
// 					"last_updated": time.Now(),
// 				}, nil).
// 				On(
// 					r.Table("sleep").MockAnything(),
// 				).
// 				Return(
// 					tC.ret, tC.err,
// 				)
// 			testEnv.DB.Session = mock

// 			req, err := http.NewRequest("GET", "/sleep/status", nil)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			handler := sleepToggleStatus(&testEnv)
// 			rr := httptest.NewRecorder()

// 			ctx := req.Context()
// 			ctx = context.WithValue(ctx, userContextKey, goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
// 			req = req.WithContext(ctx)
// 			handler.ServeHTTP(rr, req)
// 			assert.Equal(t, tC.status, rr.Code)
// 		})
// 	}
// }

// func TestSleepStartHandler(t *testing.T) {
// 	//TODO: verify returned values
// 	testCases := []struct {
// 		desc   string
// 		status int
// 		ret    map[string]interface{}
// 		err    error
// 	}{
// 		{
// 			desc:   "false, no sleep active",
// 			status: http.StatusOK,
// 			ret:    nil,
// 			err:    r.ErrEmptyResult,
// 		},
// 		{
// 			desc:   "true, sleep active",
// 			status: http.StatusConflict,
// 			ret: map[string]interface{}{
// 				"start":    time.Now().AddDate(0, 0, -1),
// 				"userID":   "1",
// 				"familyID": "1",
// 				"childID":  "1",
// 				"id":       "1",
// 			},
// 			err: nil,
// 		},
// 	}

// 	for _, tC := range testCases {
// 		t.Run(tC.desc, func(t *testing.T) {
// 			var testEnv config.Env

// 			mock := r.NewMock()
// 			mock.
// 				On(
// 					r.Table("family").Filter(
// 						func(row r.Term) r.Term {
// 							return row.Field("members").Contains("1")
// 						},
// 					),
// 				).
// 				Return(map[string]interface{}{
// 					"id":           "1",
// 					"admin":        "1",
// 					"members":      []string{"1"},
// 					"created_at":   time.Now(),
// 					"last_updated": time.Now(),
// 				}, nil).
// 				On(
// 					r.Table("sleep").MockAnything(),
// 				).
// 				Return(
// 					tC.ret, tC.err,
// 				)
// 			testEnv.DB.Session = mock

// 			req, err := http.NewRequest("GET", "/sleep/start", nil)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			handler := sleepStartHandler(&testEnv)
// 			rr := httptest.NewRecorder()

// 			ctx := req.Context()
// 			ctx = context.WithValue(ctx, userContextKey, goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
// 			req = req.WithContext(ctx)
// 			handler.ServeHTTP(rr, req)
// 			assert.Equal(t, tC.status, rr.Code)
// 		})
// 	}
// }

// func TestSleepEndHandler(t *testing.T) {
// 	//TODO: verify returned values
// 	testCases := []struct {
// 		desc   string
// 		status int
// 		ret    map[string]interface{}
// 		err    error
// 	}{
// 		{
// 			desc:   "false, no sleep active",
// 			status: http.StatusNotFound,
// 			ret:    nil,
// 			err:    r.ErrEmptyResult,
// 		},
// 		{
// 			desc:   "true, sleep active",
// 			status: http.StatusOK,
// 			ret: map[string]interface{}{
// 				"start":  time.Now().AddDate(0, 0, -1),
// 				"userid": "1",
// 				"id":     "1",
// 			},
// 			err: nil,
// 		},
// 	}

// 	for _, tC := range testCases {
// 		t.Run(tC.desc, func(t *testing.T) {
// 			var testEnv config.Env

// 			mock := r.NewMock()
// 			mock.
// 				On(
// 					r.Table("family").Filter(
// 						func(row r.Term) r.Term {
// 							return row.Field("members").Contains("1")
// 						},
// 					),
// 				).
// 				Return(map[string]interface{}{
// 					"id":           "1",
// 					"admin":        "1",
// 					"members":      []string{"1"},
// 					"created_at":   time.Now(),
// 					"last_updated": time.Now(),
// 				}, nil).
// 				On(
// 					r.Table("sleep").MockAnything(),
// 				).
// 				Return(
// 					tC.ret, tC.err,
// 				)
// 			testEnv.DB.Session = mock

// 			req, err := http.NewRequest("GET", "/sleep/end", nil)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			handler := sleepEndHandler(&testEnv)
// 			rr := httptest.NewRecorder()

// 			ctx := req.Context()
// 			ctx = context.WithValue(ctx, userContextKey, goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
// 			req = req.WithContext(ctx)
// 			handler.ServeHTTP(rr, req)
// 			assert.Equal(t, tC.status, rr.Code)
// 		})
// 	}
// }

// func TestSleepViewHandler(t *testing.T) {
// 	var testEnv config.Env

// 	req, err := http.NewRequest("GET", "/sleep/1", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	handler := sleepViewHandler(&testEnv)
// 	rr := httptest.NewRecorder()

// 	ctx := req.Context()
// 	ctx = context.WithValue(ctx, userContextKey, goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
// 	req = req.WithContext(ctx)
// 	handler.ServeHTTP(rr, req)
// 	assert.Equal(t, http.StatusOK, rr.Code)
// }

// func TestSleepEditHandler(t *testing.T) {
// 	var testEnv config.Env

// 	req, err := http.NewRequest("GET", "/sleep/1", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	handler := sleepEditHandler(&testEnv)
// 	rr := httptest.NewRecorder()

// 	ctx := req.Context()
// 	ctx = context.WithValue(ctx, userContextKey, goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
// 	req = req.WithContext(ctx)
// 	handler.ServeHTTP(rr, req)
// 	assert.Equal(t, http.StatusOK, rr.Code)
// }

// func TestSleepDeleteHandler(t *testing.T) {
// 	var testEnv config.Env

// 	req, err := http.NewRequest("GET", "/sleep/1", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	handler := sleepDeleteHandler(&testEnv)
// 	rr := httptest.NewRecorder()

// 	ctx := req.Context()
// 	ctx = context.WithValue(ctx, userContextKey, goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"})
// 	req = req.WithContext(ctx)
// 	handler.ServeHTTP(rr, req)
// 	assert.Equal(t, http.StatusOK, rr.Code)
// }

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
