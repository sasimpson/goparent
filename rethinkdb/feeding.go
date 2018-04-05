package rethinkdb

import (
	"time"

	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/config"
	gorethink "gopkg.in/gorethink/gorethink.v3"
)

type FeedingService struct {
	Env *config.Env
}

//Save - save the structure to the datastore
func (fs *FeedingService) Save(feeding *goparent.Feeding) error {
	session, err := fs.Env.DB.GetConnection()
	if err != nil {
		return err
	}

	res, err := gorethink.Table("feeding").Insert(feeding, gorethink.InsertOpts{Conflict: "replace"}).RunWrite(session)
	if err != nil {
		return err
	}

	if res.Inserted > 0 {
		feeding.ID = res.GeneratedKeys[0]
	}
	return nil
}

//Feeding - get all records for a user from the datastore
func (fs *FeedingService) Feeding(family *goparent.Family) ([]*goparent.Feeding, error) {
	session, err := fs.Env.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	res, err := gorethink.Table("feeding").
		Filter(map[string]interface{}{
			"familyID": family.ID,
		}).
		OrderBy(gorethink.Desc("timestamp")).
		Run(session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var rows []*goparent.Feeding
	err = res.All(&rows)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

//Stats - get feeding stats for one child for the last 24 hours.
func (fs *FeedingService) Stats(child *goparent.Child) (*goparent.FeedingSummary, error) {
	session, err := fs.Env.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	end := time.Now()
	start := end.AddDate(0, 0, -1)

	res, err := gorethink.Table("feeding").
		Filter(map[string]interface{}{
			"childID": child.ID,
		}).
		Filter(
			gorethink.Row.Field("timestamp").During(start, end),
		).
		OrderBy(gorethink.Desc("timestamp")).
		Run(session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var rows []goparent.Feeding
	err = res.All(&rows)
	if err != nil {
		return nil, err
	}
	//build summary output
	summary := &goparent.FeedingSummary{
		Data:  rows,
		Total: make(map[string]float32),
		Mean:  make(map[string]float32),
		Range: make(map[string]int),
	}

	for _, x := range rows {
		if _, ok := summary.Total[x.Type]; !ok {
			summary.Total[x.Type] = 0.0
		}
		summary.Total[x.Type] += x.Amount
		summary.Range[x.Type]++
	}
	for k := range summary.Total {
		summary.Mean[k] = summary.Total[k] / float32(summary.Range[k])
	}
	return summary, nil
}
