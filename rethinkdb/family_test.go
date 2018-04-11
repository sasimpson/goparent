package rethinkdb

import (
	"errors"
	"testing"
	"time"

	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/config"
	"github.com/stretchr/testify/assert"

	r "gopkg.in/gorethink/gorethink.v3"
)

func TestFamilySave(t *testing.T) {
	testCases := []struct {
		desc        string
		family      *goparent.Family
		query       *r.MockQuery
		env         *config.Env
		returnError error
	}{
		{
			desc: "Save 1",
			family: &goparent.Family{
				Admin:       "1",
				Members:     []string{"1"},
				CreatedAt:   time.Now(),
				LastUpdated: time.Now(),
			},
			//one issue with the rethinkdb mocking is that you cannot mock out
			query: (&r.Mock{}).On(r.Table("family").MockAnything()).Once().Return(
				r.WriteResponse{
					Inserted:      1,
					Errors:        0,
					GeneratedKeys: []string{"1"},
				}, nil),
			env:         &config.Env{},
			returnError: nil,
		},
		{
			desc: "Save 2",
			family: &goparent.Family{
				ID:          "1",
				Admin:       "1",
				Members:     []string{"1"},
				CreatedAt:   time.Now(),
				LastUpdated: time.Now(),
			},
			query: (&r.Mock{}).On(r.Table("family").MockAnything()).Once().Return(
				r.WriteResponse{
					Replaced:      1,
					Updated:       0,
					Inserted:      0,
					Errors:        0,
					GeneratedKeys: []string{"1"}}, nil),
			env:         &config.Env{},
			returnError: nil,
		},
		{
			desc: "Save error",
			family: &goparent.Family{
				ID:          "1",
				Admin:       "1",
				Members:     []string{"1"},
				CreatedAt:   time.Now(),
				LastUpdated: time.Now(),
			},
			query: (&r.Mock{}).On(r.Table("family").MockAnything()).Once().Return(
				r.WriteResponse{
					Replaced:      0,
					Updated:       0,
					Inserted:      0,
					Errors:        1,
					GeneratedKeys: []string{"1"}}, errors.New("test error")),
			env:         &config.Env{},
			returnError: errors.New("test error"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)
			tC.env.DB = config.DBEnv{Session: mock}
			fs := FamilyService{Env: tC.env}
			err := fs.Save(tC.family)
			if tC.returnError != nil {
				assert.Error(t, tC.returnError, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestFamily(t *testing.T) {
	timestamp := time.Now()
	testCases := []struct {
		desc        string
		env         *config.Env
		id          string
		family      *goparent.Family
		query       *r.MockQuery
		returnError error
	}{
		{
			desc: "return family",
			env:  &config.Env{},
			id:   "family-1",
			family: &goparent.Family{
				ID:          "family-1",
				Admin:       "1",
				Members:     []string{"1"},
				CreatedAt:   timestamp,
				LastUpdated: timestamp,
			},
			query: (&r.Mock{}).On(
				r.Table("family").Get("family-1"),
			).Return(
				map[string]interface{}{
					"id":           "family-1",
					"admin":        "1",
					"members":      []string{"1"},
					"created_at":   timestamp,
					"last_updated": timestamp,
				}, nil),
		},
		{
			desc: "return error",
			env:  &config.Env{},
			id:   "family-1",
			family: &goparent.Family{
				ID:          "family-1",
				Admin:       "1",
				Members:     []string{"1"},
				CreatedAt:   timestamp,
				LastUpdated: timestamp,
			},
			query: (&r.Mock{}).On(
				r.Table("family").Get("family-1"),
			).Return(
				nil, r.ErrEmptyResult),
			returnError: r.ErrEmptyResult,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)
			tC.env.DB = config.DBEnv{Session: mock}
			fs := FamilyService{Env: tC.env}
			family, err := fs.Family(tC.id)
			if tC.returnError != nil {
				assert.Error(t, tC.returnError, err)
			} else {
				assert.Nil(t, err)
				assert.EqualValues(t, tC.family, family)
			}
		})
	}
}

func TestChildren(t *testing.T) {
	timestamp := time.Now()
	testCases := []struct {
		desc         string
		env          *config.Env
		family       *goparent.Family
		query        *r.MockQuery
		resultLength int
		returnError  error
	}{
		{
			desc: "",
			env:  &config.Env{},
			family: &goparent.Family{
				ID:          "family-1",
				Admin:       "1",
				Members:     []string{"1"},
				CreatedAt:   timestamp,
				LastUpdated: timestamp,
			},
			query: (&r.Mock{}).On(
				r.Table("children").Filter(
					map[string]interface{}{
						"familyID": "family-1",
					}).OrderBy(r.Desc("birthday")),
			).Return([]map[string]interface{}{
				{
					"id":       "child-1",
					"name":     "test-child-1",
					"userID":   "user-1",
					"familyID": "family-1",
					"birthday": timestamp.AddDate(0, -30, 0),
				},
			}, nil),
			resultLength: 1,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)
			tC.env.DB = config.DBEnv{Session: mock}
			fs := FamilyService{Env: tC.env}
			children, err := fs.Children(tC.family)
			if tC.returnError != nil {
				assert.Error(t, tC.returnError, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tC.resultLength, len(children))
			}
		})
	}
}

// func TestFamily_AddMember(t *testing.T) {
// 	type args struct {
// 		env       *config.Env
// 		newMember *User
// 	}
// 	tests := []struct {
// 		name    string
// 		family  *Family
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Already in family error",
// 			family: &Family{
// 				ID:          "1",
// 				Admin:       "1",
// 				Members:     []string{"1", "2"},
// 				CreatedAt:   time.Now(),
// 				LastUpdated: time.Now(),
// 			},
// 			args: args{
// 				env: &config.Env{},
// 				newMember: &User{
// 					ID:            "2",
// 					Name:          "test user jr",
// 					Username:      "testuserjr",
// 					Email:         "testuserjr@test.com",
// 					CurrentFamily: "1",
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		// {
// 		// 	name: "Added",
// 		// 	family: &Family{
// 		// 		ID:          "1",
// 		// 		Admin:       "1",
// 		// 		Members:     []string{"1"},
// 		// 		CreatedAt:   time.Now(),
// 		// 		LastUpdated: time.Now(),
// 		// 	},
// 		// 	args: args{
// 		// 		env: &config.Env{},
// 		// 		newMember: &User{
// 		// 			ID:            "2",
// 		// 			Name:          "test user jr",
// 		// 			Username:      "testuserjr",
// 		// 			Email:         "testuserjr@test.com",
// 		// 			CurrentFamily: "2",
// 		// 		},
// 		// 	},
// 		// 	wantErr: false,
// 		// },
// 	}
// 	for _, tt := range tests {
// 		if err := tt.family.AddMember(tt.args.env, tt.args.newMember); (err != nil) != tt.wantErr {
// 			t.Errorf("%q. Family.AddMember() error = %v, wantErr %v", tt.name, err, tt.wantErr)
// 		}
// 	}
// }
