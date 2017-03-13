package models

import (
	"errors"
	"log"
	"time"

	"gopkg.in/gorethink/gorethink.v3"
)

//Sleep - tracks the baby's sleep start and end.
type Sleep struct {
	ID         string    `json:"id" gorethink:"id"`
	SleepStart time.Time `json:"start" gorethink:"start"`
	SleepEnd   time.Time `json:"end" gorethink:"end"`
	OwnerID    string    `json:"user_id" gorethink:"user_id"`
}

//StartSleep - record start of sleep
func (sleep *Sleep) StartSleep() {
	sleep.SleepStart = time.Now()
}

//EndSleep - record end of sleep
func (sleep *Sleep) EndSleep() {
	sleep.SleepEnd = time.Now()
}

//Save - creates/saves the record.  saves if there is an id filled in.
func (sleep *Sleep) Save() error {
	session, err := GetConnection()
	if err != nil {
		return err
	}
	defer session.Close()
	log.Printf("sleep: %v", sleep.OwnerID)
	if sleep.OwnerID != "" {
		resp, err := gorethink.Table("users").Insert(sleep, gorethink.InsertOpts{Conflict: "replace"}).RunWrite(session)
		if err != nil {
			log.Println("error with upsert from sleep upsert in sleep.Save()")
			return err
		}
		if resp.Inserted > 0 {
			sleep.ID = resp.GeneratedKeys[0]
		}
	}

	return errors.New("an owner should be included")
}
