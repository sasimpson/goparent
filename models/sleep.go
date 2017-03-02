package models

import (
	"time"
)

//Sleep - tracks the baby's sleep start and end.
type Sleep struct {
	ID         string
	SleepStart time.Time
	SleepEnd   time.Time
}

//StartSleep - record start of sleep
func (s *Sleep) StartSleep() {
	s.SleepStart = time.Now()
}

//EndSleep - record end of sleep
func (s *Sleep) EndSleep() {
	s.SleepEnd = time.Now()
}
