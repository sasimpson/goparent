package rethinkdb

import (
	"testing"
	"time"

	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/config"
	"github.com/stretchr/testify/assert"
	r "gopkg.in/gorethink/gorethink.v3"
)

func TestSave(t *testing.T) {
	timestamp := time.Now()
	testCases := []struct {
		desc        string
		env         *config.Env
		child       *goparent.Child
		query       *r.MockQuery
		returnError error
	}{
		{
			desc:  "save child",
			env:   &config.Env{},
			child: &goparent.Child{Name: "test child", ParentID: "1", FamilyID: "1", Birthday: timestamp.AddDate(-1, 0, 0)},
			query: (&r.Mock{}).On(r.Table("children").MockAnything()).Once().Return(
				r.WriteResponse{
					Replaced:      0,
					Updated:       0,
					Inserted:      1,
					Errors:        0,
					GeneratedKeys: []string{"1"}}, nil),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mock := r.NewMock()
			mock.ExpectedQueries = append(mock.ExpectedQueries, tC.query)
			tC.env.DB = config.DBEnv{Session: mock}
			cs := ChildService{Env: tC.env}
			err := cs.Save(tC.child)
			if tC.returnError != nil {
				assert.Equal(t, tC.returnError, err)
			} else {
				assert.Nil(t, nil)
			}
		})
	}
}
