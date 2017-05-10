package models

import (
	"errors"
	"log"
	"time"

	"gopkg.in/gorethink/gorethink.v3"
)

//Sleep - tracks the baby's sleep start and end.
type Sleep struct {
	ID         string    `json:"id" gorethink:"id,omitempty"`
	SleepStart time.Time `json:"start" gorethink:"start"`
	SleepEnd   time.Time `json:"end" gorethink:"end"`
	UserID     string    `json:"userid" gorethink:"userid"`
}

var ExistingStartErr = errors.New("already have a start record")
var NoExistingSessionErr = errors.New("no existing sleep session to end")

func (sleep *Sleep) Status(user *User) (bool, error) {
	session, err := GetConnection()
	if err != nil {
		return false, err
	}
	defer session.Close()
	//check to see if we already have an open sleep session
	res, err := gorethink.Table("sleep").Filter(map[string]interface{}{
		"end":    time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
		"userid": user.ID,
	}).Run(session)
	if err != nil {
		return false, err
	}
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
func (sleep *Sleep) Start(user *User) error {
	ok, err := sleep.Status(user)
	if err != nil {
		return err
	}
	if !ok {
		sleep.SleepStart = time.Now()
		return nil
	}
	return ExistingStartErr

}

//End - record end of sleep
func (sleep *Sleep) End(user *User) error {
	ok, err := sleep.Status(user)
	if err != nil {
		return err
	}
	if ok {
		sleep.SleepEnd = time.Now()
		return nil
	}
	return NoExistingSessionErr

}

//Save - creates/saves the record.  saves if there is an id filled in.
func (sleep *Sleep) Save() error {
	session, err := GetConnection()
	if err != nil {
		return err
	}
	defer session.Close()
	// log.Printf("sleep: %v", sleep.OwnerID)
	// if sleep.OwnerID != "" {
	resp, err := gorethink.Table("sleep").Insert(sleep, gorethink.InsertOpts{Conflict: "replace"}).RunWrite(session)
	if err != nil {
		log.Println("error with upsert from sleep upsert in sleep.Save()")
		return err
	}
	if resp.Inserted > 0 {
		sleep.ID = resp.GeneratedKeys[0]
	}
	// }

	// return errors.New("an owner should be included")
	return nil
}

func (sleep *Sleep) GetAll(user *User) ([]Sleep, error) {
	session, err := GetConnection()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	resp, err := gorethink.Table("sleep").Filter(map[string]interface{}{"userid": user.ID}).OrderBy(gorethink.Desc("start")).Run(session)
	if err != nil {
		return nil, err
	}
	var rows []Sleep
	err = resp.All(&rows)
	if err != nil {
		log.Println("error getting all")
		return nil, err
	}
	return rows, nil
}
