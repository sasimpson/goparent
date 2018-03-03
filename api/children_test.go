package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/sasimpson/goparent/config"
	"github.com/sasimpson/goparent/models"
	"github.com/stretchr/testify/assert"

	r "gopkg.in/gorethink/gorethink.v3"
)

func Test_initChildrenHandlers(t *testing.T) {
	type args struct {
		env *config.Env
		r   *mux.Router
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		initChildrenHandlers(tt.args.env, tt.args.r)
	}
}

func TestChildSummary(t *testing.T) {
	type args struct {
		env *config.Env
	}
	tests := []struct {
		name string
		args args
		want http.Handler
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := childSummary(tt.args.env); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. ChildSummary() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestChildrenGetHandler(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		env      *config.Env
		user     map[string]string
		family   map[string]interface{}
		children []models.Child
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "get 1",
			args: args{
				env: &config.Env{},
				user: map[string]string{
					"id":       "1",
					"name":     "test user",
					"email":    "testuser@test.com",
					"username": "testuser",
				},
				family: map[string]interface{}{
					"id":           "1",
					"admin":        "1",
					"members":      []string{"1"},
					"created_at":   time.Now(),
					"last_updated": time.Now(),
				},
				children: []models.Child{
					models.Child{ID: "1", Name: "test child", ParentID: "1", FamilyID: "1", Birthday: time.Now()},
				},
			},
			want: map[string]interface{}{
				"responseCode": 200,
			},
		},
		{
			name: "get fail",
			args: args{
				env: &config.Env{},
				user: map[string]string{
					"id":       "1",
					"name":     "test user",
					"email":    "testuser@test.com",
					"username": "testuser",
				},
				family: nil,
			},
		},
	}
	for _, tt := range tests {
		mock := r.NewMock()
		mock.
			On(
				r.Table("family").Filter(
					func(row r.Term) r.Term {
						return row.Field("members").Contains(tt.args.family["id"].(string))
					},
				),
			).
			Return(tt.args.family, nil).
			On(
				r.Table("children").Filter(
					map[string]interface{}{
						"familyID": tt.args.family["id"].(string),
					},
				).OrderBy(r.Desc("birthday")),
			)
		tt.args.env.DB.Session = mock

		req, err := http.NewRequest("GET", "/children", nil)
		if err != nil {
			t.Fatal(err)
		}

		handler := childrenGetHandler(tt.args.env)
		rr := httptest.NewRecorder()

		ctx := req.Context()
		ctx = context.WithValue(ctx, userContextKey, models.User{
			ID:       tt.args.user["id"],
			Name:     tt.args.user["name"],
			Email:    tt.args.user["email"],
			Username: tt.args.user["username"],
		})
		req = req.WithContext(ctx)
		handler.ServeHTTP(rr, req)

		assert.Equal(tt.want["responseCode"], rr.Code)
	}
}

func TestChildNewHandler(t *testing.T) {
	type args struct {
		env *config.Env
	}
	tests := []struct {
		name string
		args args
		want http.Handler
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := childNewHandler(tt.args.env); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. ChildNewHandler() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestChildViewHandler(t *testing.T) {
	type args struct {
		env *config.Env
	}
	tests := []struct {
		name string
		args args
		want http.Handler
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := childViewHandler(tt.args.env); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. ChildViewHandler() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestChildEditHandler(t *testing.T) {
	type args struct {
		env *config.Env
	}
	tests := []struct {
		name string
		args args
		want http.Handler
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := childEditHandler(tt.args.env); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. ChildEditHandler() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestChildDeleteHandler(t *testing.T) {
	type args struct {
		env *config.Env
	}
	tests := []struct {
		name string
		args args
		want http.Handler
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := childDeleteHandler(tt.args.env); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. ChildDeleteHandler() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
