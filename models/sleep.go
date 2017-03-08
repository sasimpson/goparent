package models

import (
	"time"
)

//Sleep - tracks the baby's sleep start and end.
type Sleep struct {
	ID         string    `json:"id" gorethink:"id"`
	SleepStart time.Time `json:"start" gorethink:"start"`
	SleepEnd   time.Time `json:"end" gorethink:"end"`
	Owner      User      `json:"userData" gorethink:"userData"`
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

}
