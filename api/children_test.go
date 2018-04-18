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
	"github.com/sasimpson/goparent/config"
	"github.com/sasimpson/goparent/mock"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	testCases := []struct {
		desc    string
		name    string
		path    string
		methods []string
	}{
		{
			desc:    "children get",
			name:    "ChildrenGet",
			path:    "/children",
			methods: []string{"GET"},
		},
		{
			desc:    "child new",
			name:    "ChildNew",
			path:    "/children",
			methods: []string{"POST"},
		},
		{
			desc:    "child view",
			name:    "ChildView",
			path:    "/children/{id}",
			methods: []string{"GET"},
		},
		{
			desc:    "child edit",
			name:    "ChildEdit",
			path:    "/children/{id}",
			methods: []string{"PUT"},
		},
		{
			desc:    "child delete",
			name:    "ChildDelete",
			path:    "/children/{id}",
			methods: []string{"DELETE"},
		},
		{
			desc:    "child summary",
			name:    "ChildSummary",
			path:    "/children/{id}/summary",
			methods: []string{"GET"},
		},
	}

	var testEnv *config.Env
	h := Handler{Env: testEnv}
	routes := mux.NewRouter()
	h.initChildrenHandlers(routes)

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

func TestChildrenGetHandler(t *testing.T) {
	testCases := []struct {
		desc          string
		env           *config.Env
		userService   goparent.UserService
		familyService goparent.FamilyService
		childService  goparent.ChildService
		contextUser   *goparent.User
		contextError  bool
		responseCode  int
		resultLength  int
	}{
		{
			desc:          "returns auth error",
			env:           &config.Env{},
			userService:   &mock.MockUserService{},
			familyService: &mock.MockFamilyService{},
			childService:  &mock.MockChildService{},
			contextUser:   &goparent.User{},
			contextError:  true,
			responseCode:  http.StatusUnauthorized,
		},
		{
			desc: "returns family error",
			env:  &config.Env{},
			userService: &mock.MockUserService{
				GetErr: errors.New("test error"),
			},
			familyService: &mock.MockFamilyService{},
			childService:  &mock.MockChildService{},
			contextUser:   &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			contextError:  false,
			responseCode:  http.StatusInternalServerError,
		},
		{
			desc: "returns chilren error",
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
			familyService: &mock.MockFamilyService{
				GetErr: errors.New("test error"),
			},
			childService: &mock.MockChildService{},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			contextError: false,
			responseCode: http.StatusInternalServerError,
		},
		{
			desc: "return no children",
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
			familyService: &mock.MockFamilyService{
				Kids: []*goparent.Child{},
			},
			childService: &mock.MockChildService{},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			contextError: false,
			responseCode: http.StatusOK,
			resultLength: 0,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockHandler := Handler{
				Env:           tC.env,
				UserService:   tC.userService,
				FamilyService: tC.familyService,
				ChildService:  tC.childService,
			}

			req, err := http.NewRequest("GET", "/children", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := mockHandler.childrenGetHandler()
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
				var result ChildrenResponse
				decoder := json.NewDecoder(rr.Body)
				err := decoder.Decode(&result)
				assert.Nil(t, err)
				assert.Equal(t, tC.resultLength, len(result.Children))
			}
		})
	}
}

func TestChildrenNewHandler(t *testing.T) {
	testCases := []struct {
		desc          string
		env           *config.Env
		userService   goparent.UserService
		familyService goparent.FamilyService
		childService  goparent.ChildService
		childRequest  ChildRequest
		contextUser   *goparent.User
		contextError  bool
		responseCode  int
		resultLength  int
	}{
		{
			desc: "submit child",
			env:  &config.Env{},
			childRequest: ChildRequest{
				ChildData: goparent.Child{
					Name:     "test child",
					ID:       "1",
					FamilyID: "1",
					Birthday: time.Now()}},
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
			childService:  &mock.MockChildService{},
			contextUser:   &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			contextError:  false,
			responseCode:  http.StatusOK,
			resultLength:  0,
		},
		{
			desc: "returns no family error",
			env:  &config.Env{},
			childRequest: ChildRequest{
				ChildData: goparent.Child{
					Name:     "test child",
					ID:       "1",
					FamilyID: "1",
					Birthday: time.Now()}},
			userService: &mock.MockUserService{
				GetErr: errors.New("user has no current family"),
			},
			familyService: &mock.MockFamilyService{},
			childService:  &mock.MockChildService{},
			contextUser:   &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			contextError:  false,
			responseCode:  http.StatusInternalServerError,
		},
		{
			desc: "returns child error",
			env:  &config.Env{},
			childRequest: ChildRequest{
				ChildData: goparent.Child{
					Name:     "test child",
					ID:       "1",
					FamilyID: "1",
					Birthday: time.Now()}},
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
			childService: &mock.MockChildService{
				GetErr: errors.New("unknown child error"),
			},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			contextError: false,
			responseCode: http.StatusConflict,
		},
		{
			desc: "returns auth error",
			env:  &config.Env{},
			childRequest: ChildRequest{
				ChildData: goparent.Child{
					ID:       "1",
					FamilyID: "1",
					Name:     "Test Child",
					Birthday: time.Now()}},
			userService:   &mock.MockUserService{},
			familyService: &mock.MockFamilyService{},
			childService:  &mock.MockChildService{},
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
			childService:  &mock.MockChildService{},
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
				ChildService:  tC.childService,
			}

			js, err := json.Marshal(&tC.childRequest)
			if err != nil {
				t.Fatal(err)
			}
			if (tC.childRequest == ChildRequest{}) {
				js = []byte("this is a test")
			}
			req, err := http.NewRequest("POST", "/children", bytes.NewReader(js))
			if err != nil {
				t.Fatal(err)
			}

			handler := mockHandler.childNewHandler()
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
		})
	}
}

func TestChildViewHandler(t *testing.T) {
	testCases := []struct {
		desc          string
		env           *config.Env
		userService   goparent.UserService
		familyService goparent.FamilyService
		childService  goparent.ChildService
		contextUser   *goparent.User
		contextError  bool
		responseCode  int
	}{
		{

			desc: "get child",
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
			childService: &mock.MockChildService{
				Kid: &goparent.Child{
					Name:     "test child",
					ID:       "1",
					FamilyID: "1",
					Birthday: time.Now()},
			},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			responseCode: http.StatusOK,
		},
		{
			desc: "returns no family error",
			env:  &config.Env{},
			userService: &mock.MockUserService{
				GetErr: errors.New("user has no current family"),
			},
			familyService: &mock.MockFamilyService{},
			childService:  &mock.MockChildService{},
			contextUser:   &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			responseCode:  http.StatusInternalServerError,
		},
		{
			desc: "returns child error",
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
			childService: &mock.MockChildService{
				GetErr: errors.New("unknown child error"),
			},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			responseCode: http.StatusInternalServerError,
		},
		{
			desc:          "returns auth error",
			env:           &config.Env{},
			userService:   &mock.MockUserService{},
			familyService: &mock.MockFamilyService{},
			childService:  &mock.MockChildService{},
			contextUser:   &goparent.User{},
			contextError:  true,
			responseCode:  http.StatusUnauthorized,
		},
		{

			desc: "get not user's child",
			env:  &config.Env{},
			userService: &mock.MockUserService{
				Family: &goparent.Family{
					ID:          "2",
					Admin:       "1",
					Members:     []string{"1"},
					CreatedAt:   time.Now(),
					LastUpdated: time.Now(),
				},
			},
			familyService: &mock.MockFamilyService{},
			childService: &mock.MockChildService{
				Kid: &goparent.Child{
					Name:     "test child",
					ID:       "1",
					FamilyID: "1",
					Birthday: time.Now()},
			},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			responseCode: http.StatusNotFound,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockHandler := Handler{
				Env:           tC.env,
				UserService:   tC.userService,
				FamilyService: tC.familyService,
				ChildService:  tC.childService,
			}
			// req = mux.SetURLVars(req, map[string]string{"id": "1"})
			req, err := http.NewRequest("GET", "/children/1", nil)
			if err != nil {
				t.Fatal(err)
			}
			handler := mockHandler.childViewHandler()
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
		})
	}
}

func TestChildEditHandler(t *testing.T) {
	testCases := []struct {
		desc          string
		env           *config.Env
		userService   goparent.UserService
		familyService goparent.FamilyService
		childService  goparent.ChildService
		childRequest  ChildRequest
		contextUser   *goparent.User
		contextError  bool
		responseCode  int
	}{
		{
			desc: "submit child",
			env:  &config.Env{},
			childRequest: ChildRequest{
				ChildData: goparent.Child{
					Name:     "test child",
					ID:       "1",
					FamilyID: "1",
					Birthday: time.Now()}},
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
			childService: &mock.MockChildService{
				Kid: &goparent.Child{
					Name:     "test child",
					ID:       "1",
					FamilyID: "1",
					Birthday: time.Now()},
			},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			responseCode: http.StatusOK,
		},
		{
			desc: "submit child not in family",
			env:  &config.Env{},
			childRequest: ChildRequest{
				ChildData: goparent.Child{
					Name:     "test child",
					ID:       "1",
					FamilyID: "1",
					Birthday: time.Now()}},
			userService: &mock.MockUserService{
				Family: &goparent.Family{
					ID:          "2",
					Admin:       "1",
					Members:     []string{"1"},
					CreatedAt:   time.Now(),
					LastUpdated: time.Now(),
				},
			},
			familyService: &mock.MockFamilyService{},
			childService: &mock.MockChildService{
				Kid: &goparent.Child{
					Name:     "test child",
					ID:       "2",
					FamilyID: "1",
					Birthday: time.Now()},
			},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			responseCode: http.StatusBadRequest,
		},
		{
			desc: "returns no family error",
			env:  &config.Env{},
			childRequest: ChildRequest{
				ChildData: goparent.Child{
					Name:     "test child",
					ID:       "1",
					FamilyID: "1",
					Birthday: time.Now()}},
			userService: &mock.MockUserService{
				GetErr: errors.New("user has no current family"),
			},
			familyService: &mock.MockFamilyService{},
			childService:  &mock.MockChildService{},
			contextUser:   &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			contextError:  false,
			responseCode:  http.StatusInternalServerError,
		},
		{
			desc: "returns child error",
			env:  &config.Env{},
			childRequest: ChildRequest{
				ChildData: goparent.Child{
					Name:     "test child",
					ID:       "1",
					FamilyID: "1",
					Birthday: time.Now()}},
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
			childService: &mock.MockChildService{
				GetErr: errors.New("unknown child error"),
			},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			contextError: false,
			responseCode: http.StatusNotFound,
		},
		{
			desc: "returns auth error",
			env:  &config.Env{},
			childRequest: ChildRequest{
				ChildData: goparent.Child{
					ID:       "1",
					FamilyID: "1",
					Name:     "Test Child",
					Birthday: time.Now()}},
			userService:   &mock.MockUserService{},
			familyService: &mock.MockFamilyService{},
			childService:  &mock.MockChildService{},
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
			childService:  &mock.MockChildService{},
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
				ChildService:  tC.childService,
			}

			js, err := json.Marshal(&tC.childRequest)
			if err != nil {
				t.Fatal(err)
			}
			if (tC.childRequest == ChildRequest{}) {
				js = []byte("this is a test")
			}
			req, err := http.NewRequest("PUT", "/children/1", bytes.NewReader(js))
			if err != nil {
				t.Fatal(err)
			}

			handler := mockHandler.childEditHandler()
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
		})
	}
}

func TestChildDeleteHandler(t *testing.T) {
	testCases := []struct {
		desc          string
		env           *config.Env
		userService   goparent.UserService
		familyService goparent.FamilyService
		childService  goparent.ChildService
		childRequest  ChildRequest
		contextUser   *goparent.User
		contextError  bool
		responseCode  int
	}{
		{
			desc: "delete child",
			env:  &config.Env{},
			childRequest: ChildRequest{
				ChildData: goparent.Child{
					Name:     "test child",
					ID:       "1",
					FamilyID: "1",
					Birthday: time.Now()}},
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
			childService: &mock.MockChildService{
				Kid: &goparent.Child{
					Name:     "test child",
					ID:       "1",
					FamilyID: "1",
					Birthday: time.Now()},
				Deleted: 1,
			},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			responseCode: http.StatusAccepted,
		},
		{
			desc: "delete child, get child error",
			env:  &config.Env{},
			childRequest: ChildRequest{
				ChildData: goparent.Child{
					Name:     "test child",
					ID:       "1",
					FamilyID: "1",
					Birthday: time.Now()}},
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
			childService: &mock.MockChildService{
				Kid: &goparent.Child{
					Name:     "test child",
					ID:       "1",
					FamilyID: "1",
					Birthday: time.Now()},
				GetErr:  errors.New("test error"),
				Deleted: 0,
			},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			responseCode: http.StatusInternalServerError,
		},
		{
			desc: "delete child incorrect family",
			env:  &config.Env{},
			childRequest: ChildRequest{
				ChildData: goparent.Child{
					Name:     "test child",
					ID:       "1",
					FamilyID: "1",
					Birthday: time.Now()}},
			userService: &mock.MockUserService{
				Family: &goparent.Family{
					ID:          "2",
					Admin:       "1",
					Members:     []string{"1"},
					CreatedAt:   time.Now(),
					LastUpdated: time.Now(),
				},
			},
			familyService: &mock.MockFamilyService{},
			childService: &mock.MockChildService{
				Kid: &goparent.Child{
					Name:     "test child",
					ID:       "1",
					FamilyID: "1",
					Birthday: time.Now()},
				Deleted: 1,
			},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			responseCode: http.StatusNotFound,
		},
		{
			desc: "delete child, get delete error",
			env:  &config.Env{},
			childRequest: ChildRequest{
				ChildData: goparent.Child{
					Name:     "test child",
					ID:       "1",
					FamilyID: "1",
					Birthday: time.Now()}},
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
			childService: &mock.MockChildService{
				Kid: &goparent.Child{
					Name:     "test child",
					ID:       "1",
					FamilyID: "1",
					Birthday: time.Now()},
				DeleteErr: errors.New("test error"),
				Deleted:   0,
			},
			contextUser:  &goparent.User{ID: "1", Name: "test user", Email: "testuser@test.com", Username: "testuser"},
			responseCode: http.StatusInternalServerError,
		},
		{
			desc: "returns auth error",
			env:  &config.Env{},
			childRequest: ChildRequest{
				ChildData: goparent.Child{
					ID:       "1",
					FamilyID: "1",
					Name:     "Test Child",
					Birthday: time.Now()}},
			userService:   &mock.MockUserService{},
			familyService: &mock.MockFamilyService{},
			childService:  &mock.MockChildService{},
			contextUser:   &goparent.User{},
			contextError:  true,
			responseCode:  http.StatusUnauthorized,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockHandler := Handler{
				Env:           tC.env,
				UserService:   tC.userService,
				FamilyService: tC.familyService,
				ChildService:  tC.childService,
			}

			req, err := http.NewRequest("DELETE", "/children/1", nil)
			if err != nil {
				t.Fatal(err)
			}

			handler := mockHandler.childDeleteHandler()
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
		})
	}
}
