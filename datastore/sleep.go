package datastore

import "github.com/sasimpson/goparent"

//SleepService -
type SleepService struct {
	Env *goparent.Env
}

//Save -
func (s *SleepService) Save(*goparent.Sleep) error {
	panic("not implemented")
}

//Sleep -
func (s *SleepService) Sleep(*goparent.Family, uint64) ([]*goparent.Sleep, error) {
	panic("not implemented")
}

//Stats -
func (s *SleepService) Stats(*goparent.Child) (*goparent.SleepSummary, error) {
	panic("not implemented")
}

//Status -
func (s *SleepService) Status(*goparent.Family, *goparent.Child) (bool, error) {
	panic("not implemented")
}

//Start -
func (s *SleepService) Start(*goparent.Sleep, *goparent.Family, *goparent.Child) error {
	panic("not implemented")
}

//End -
func (s *SleepService) End(*goparent.Sleep, *goparent.Family, *goparent.Child) error {
	panic("not implemented")
}

//GraphData -
func (s *SleepService) GraphData(*goparent.Child) (*goparent.SleepChartData, error) {
	panic("not implemented")
}
