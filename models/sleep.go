package models

import (
	"errors"
	"time"

	"github.com/sasimpson/goparent/config"
	"gopkg.in/gorethink/gorethink.v3"
)

//Sleep - tracks the baby's sleep start and end.
type Sleep struct {
	ID      string    `json:"id" gorethink:"id,omitempty"`
	Start   time.Time `json:"start" gorethink:"start"`
	End     time.Time `json:"end" gorethink:"end"`
	UserID  string    `json:"userid" gorethink:"userid"`
	ChildID string    `json:"childID" gorethink:"childid"`
}

type SleepSummary struct {
	Data  []Sleep `json:"data"`
	Total int64   `json:"total"`
	Mean  float64 `json:"mean"`
	Range int     `json:"range"`
}

var ExistingStartErr = errors.New("already have a start record")
var NoExistingSessionErr = errors.New("no existing sleep session to end")

func (sleep *Sleep) Status(env *config.Env, user *User) (bool, error) {
	session, err := env.DB.GetConnection()
	if err != nil {
		return false, err
	}
	//check to see if we already have an open sleep session
	res, err := gorethink.Table("sleep").Filter(map[string]interface{}{
		"end":    time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
		"userid": user.ID,
	}).Run(session)
	if err != nil {
		if err == gorethink.ErrEmptyResult {
			return false, nil
		}
		return false, err
	}
	res.Close()
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
func (sleep *Sleep) SleepStart(env *config.Env, user *User) error {
	ok, err := sleep.Status(env, user)
	if err != nil {
		return err
	}
	if !ok {
		sleep.Start = time.Now()
		return nil
	}
	return ExistingStartErr

}

//End - record end of sleep
func (sleep *Sleep) SleepEnd(env *config.Env, user *User) error {
	ok, err := sleep.Status(env, user)
	if err != nil {
		return err
	}
	if ok {
		sleep.End = time.Now()
		return nil
	}
	return NoExistingSessionErr

}

//Save - creates/saves the record.  saves if there is an id filled in.
func (sleep *Sleep) Save(env *config.Env) error {
	session, err := env.DB.GetConnection()
	if err != nil {
		return err
	}

	resp, err := gorethink.Table("sleep").Insert(sleep, gorethink.InsertOpts{Conflict: "replace"}).RunWrite(session)
	if err != nil {
		// log.Println("error with upsert from sleep upsert in sleep.Save()")
		return err
	}
	if resp.Inserted > 0 {
		sleep.ID = resp.GeneratedKeys[0]
	}

	return nil
}

func (sleep *Sleep) GetAll(env *config.Env, user *User) ([]Sleep, error) {
	session, err := env.DB.GetConnection()
	if err != nil {
		return nil, err
	}
	res, err := gorethink.Table("sleep").
		Filter(map[string]interface{}{
			"userid": user.ID,
		}).
		OrderBy(gorethink.Desc("end")).
		Run(session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var rows []Sleep
	err = res.All(&rows)
	if err != nil {
		// log.Println("error getting all")
		return nil, err
	}
	return rows, nil
}

//SleepGetStats - get sleep stats for one child for the last 24 hours.
func SleepGetStats(env *config.Env, user *User, child *Child) (SleepSummary, error) {
	session, err := env.DB.GetConnection()
	if err != nil {
		return SleepSummary{}, err
	}

	end := time.Now()
	start := end.AddDate(0, 0, -1)

	res, err := gorethink.Table("sleep").
		Filter(map[string]interface{}{
			"userid":  user.ID,
			"childid": child.ID,
		}).
		Filter(gorethink.Row.Field("start").During(start, end)).
		OrderBy(gorethink.Desc("start")).
		Run(session)
	if err != nil {
		return SleepSummary{}, err
	}
	defer res.Close()

	var rows []Sleep
	err = res.All(&rows)
	if err != nil {
		return SleepSummary{}, err
	}
	//build summary output
	summary := SleepSummary{
		Data: rows,
	}

	for _, x := range rows {
		summary.Total += x.End.Unix() - x.Start.Unix()
	}
	summary.Range = len(rows)
	summary.Mean = float64(summary.Total) / float64(summary.Range)
	return summary, nil
}
