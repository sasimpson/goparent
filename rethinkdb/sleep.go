package rethinkdb

import (
	"errors"
	"time"

	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/config"
	"gopkg.in/gorethink/gorethink.v3"
)

//SleepService - struct for implementing the interface
type SleepService struct {
	Env *config.Env
}

//ErrExistingStart - already have a start for that sleep record
var ErrExistingStart = errors.New("already have a start record")

//ErrNoExistingSession - don't have a sleep record to end.
var ErrNoExistingSession = errors.New("no existing sleep session to end")

//Status - return the current status for a sleep session
func (ss *SleepService) Status(family *goparent.Family, child *goparent.Child) (bool, error) {
	session, err := ss.Env.DB.GetConnection()
	if err != nil {
		return false, err
	}

	//check to see if we already have an open sleep session
	res, err := gorethink.Table("sleep").Filter(map[string]interface{}{
		"end":      time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
		"familyID": family.ID,
		"childID":  child.ID,
	}).Run(session)
	if err != nil {
		if err == gorethink.ErrEmptyResult {
			return false, nil
		}
		return false, err
	}
	defer res.Close()
	var sleep goparent.Sleep
	err = res.One(&sleep)
	if err != nil {
		//if we don't, then set the sleep start as now and return
		if err == gorethink.ErrEmptyResult {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

//Start - record start of sleep
func (ss *SleepService) Start(sleep *goparent.Sleep, family *goparent.Family, child *goparent.Child) error {
	ok, err := ss.Status(family, child)
	if err != nil {
		return err
	}

	if !ok {
		sleep.Start = time.Now()
		return nil
	}
	return ErrExistingStart
}

//End - record end of sleep
func (ss *SleepService) End(sleep *goparent.Sleep, family *goparent.Family, child *goparent.Child) error {
	ok, err := ss.Status(family, child)
	if err != nil {
		return err
	}
	if ok {
		sleep.End = time.Now()
		return nil
	}
	return ErrNoExistingSession

}

//Save - creates/saves the record.  saves if there is an id filled in.
func (ss *SleepService) Save(sleep *goparent.Sleep) error {
	session, err := ss.Env.DB.GetConnection()
	if err != nil {
		return err
	}

	resp, err := gorethink.Table("sleep").Insert(sleep, gorethink.InsertOpts{Conflict: "replace"}).RunWrite(session)
	if err != nil {
		return err
	}
	if resp.Inserted > 0 {
		sleep.ID = resp.GeneratedKeys[0]
	}

	return nil
}

//Sleep - get all sleeps for a user (parent)
func (ss *SleepService) Sleep(family *goparent.Family, days uint64) ([]*goparent.Sleep, error) {
	session, err := ss.Env.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	daysBack := int(0 - days)
	res, err := gorethink.Table("sleep").
		Filter(map[string]interface{}{
			"familyID": family.ID,
		}).
		Filter(gorethink.Row.Field("timestamp").During(time.Now().AddDate(0, 0, daysBack), time.Now())).
		OrderBy(gorethink.Desc("end")).
		Run(session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var rows []*goparent.Sleep
	err = res.All(&rows)
	if err != nil {
		// log.Println("error getting all")
		return nil, err
	}
	return rows, nil
}

//Stats - get sleep stats for one child for the last 24 hours.
func (ss *SleepService) Stats(child *goparent.Child) (*goparent.SleepSummary, error) {
	session, err := ss.Env.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	end := time.Now()
	start := end.AddDate(0, 0, -1)

	res, err := gorethink.Table("sleep").
		Filter(map[string]interface{}{
			"childID": child.ID,
		}).
		Filter(gorethink.Row.Field("start").During(start, end)).
		OrderBy(gorethink.Desc("start")).
		Run(session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var rows []goparent.Sleep
	err = res.All(&rows)
	if err != nil {
		return nil, err
	}
	//build summary output
	summary := goparent.SleepSummary{
		Data: rows,
	}

	for _, x := range rows {
		summary.Total += x.End.Unix() - x.Start.Unix()
	}
	summary.Range = len(rows)
	summary.Mean = float64(summary.Total) / float64(summary.Range)
	return &summary, nil
}

//GraphData -
func (ss *SleepService) GraphData(child *goparent.Child) (*goparent.SleepChartData, error) {
	// session, err := ss.Env.DB.GetConnection()
	// if err != nil {
	// 	return nil, err
	// }

	// end := time.Now()
	// start := end.AddDate(0, 0, -7)

	// res, err := gorethink.Table("sleep").
	// 	Filter(gorethink.Row.Field("timestamp").During(start, end)).OrderBy("timestamp").
	// 	Group(
	// 		gorethink.Row.Field("timestamp").Year(),
	// 		gorethink.Row.Field("timestamp").Month(),
	// 		gorethink.Row.Field("timestamp").Day(),
	// 	).Run(session)
	return nil, nil
}
