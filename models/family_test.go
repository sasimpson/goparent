package models

import (
	"reflect"
	"testing"
	"time"

	"github.com/sasimpson/goparent/config"

	r "gopkg.in/gorethink/gorethink.v3"
)

func TestFamily_Save(t *testing.T) {
	type args struct {
		env *config.Env
	}
	tests := []struct {
		name    string
		family  *Family
		args    args
		wantErr bool
		wantDB  r.WriteResponse
	}{
		{
			name: "Save 1",
			family: &Family{
				Admin:       "1",
				Members:     []string{"1"},
				CreatedAt:   time.Now(),
				LastUpdated: time.Now(),
			},
			args:    args{env: &config.Env{}},
			wantErr: false,
			wantDB: r.WriteResponse{
				Inserted:      1,
				Errors:        0,
				GeneratedKeys: []string{"1"},
			},
		},
		{
			name: "Save 2",
			family: &Family{
				ID:          "1",
				Admin:       "1",
				Members:     []string{"1"},
				CreatedAt:   time.Now(),
				LastUpdated: time.Now(),
			},
			args:    args{env: &config.Env{}},
			wantErr: false,
			wantDB: r.WriteResponse{
				Replaced:      1,
				Updated:       0,
				Inserted:      0,
				Errors:        0,
				GeneratedKeys: []string{"1"},
			},
		},
	}
	for _, tt := range tests {
		mock := r.NewMock()
		mock.
			On(
				r.Table("family").MockAnything(),
			).
			Return(tt.wantDB, nil)
		tt.args.env.DB.Session = mock
		if err := tt.family.Save(tt.args.env); (err != nil) != tt.wantErr {
			t.Errorf("%q. Family.Save() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestFamily_GetFamily(t *testing.T) {
	type args struct {
		env *config.Env
		id  string
	}
	tests := []struct {
		name    string
		family  *Family
		args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if err := tt.family.GetFamily(tt.args.env, tt.args.id); (err != nil) != tt.wantErr {
			t.Errorf("%q. Family.GetFamily() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestFamily_GetAllChildren(t *testing.T) {
	type args struct {
		env *config.Env
	}
	tests := []struct {
		name    string
		family  *Family
		args    args
		want    []Child
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		got, err := tt.family.GetAllChildren(tt.args.env)
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Family.GetAllChildren() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Family.GetAllChildren() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestFamily_AddMember(t *testing.T) {
	type args struct {
		env       *config.Env
		newMember *User
	}
	tests := []struct {
		name    string
		family  *Family
		args    args
		wantErr bool
	}{
		{
			name: "Already in family error",
			family: &Family{
				ID:          "1",
				Admin:       "1",
				Members:     []string{"1", "2"},
				CreatedAt:   time.Now(),
				LastUpdated: time.Now(),
			},
			args: args{
				env: &config.Env{},
				newMember: &User{
					ID:            "2",
					Name:          "test user jr",
					Username:      "testuserjr",
					Email:         "testuserjr@test.com",
					CurrentFamily: "1",
				},
			},
			wantErr: true,
		},
		// {
		// 	name: "Added",
		// 	family: &Family{
		// 		ID:          "1",
		// 		Admin:       "1",
		// 		Members:     []string{"1"},
		// 		CreatedAt:   time.Now(),
		// 		LastUpdated: time.Now(),
		// 	},
		// 	args: args{
		// 		env: &config.Env{},
		// 		newMember: &User{
		// 			ID:            "2",
		// 			Name:          "test user jr",
		// 			Username:      "testuserjr",
		// 			Email:         "testuserjr@test.com",
		// 			CurrentFamily: "2",
		// 		},
		// 	},
		// 	wantErr: false,
		// },
	}
	for _, tt := range tests {
		if err := tt.family.AddMember(tt.args.env, tt.args.newMember); (err != nil) != tt.wantErr {
			t.Errorf("%q. Family.AddMember() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
