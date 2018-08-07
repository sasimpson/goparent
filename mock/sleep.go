package mock

import (
	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/config"
)

//SleepService -
type SleepService struct {
	Env     *config.Env
	Sleeps  []*goparent.Sleep
	Stat    *goparent.SleepSummary
	GetErr  error
	StatErr error
}

//Save -
func (m *SleepService) Save(*goparent.Sleep) error {
	if m.GetErr != nil {
		return m.GetErr
	}
	return nil
}

//Sleep -
func (m *SleepService) Sleep(*goparent.Family) ([]*goparent.Sleep, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	if m.Sleeps != nil {
		return m.Sleeps, nil
	}
	return nil, nil
}

//Stats -
func (m *SleepService) Stats(*goparent.Child) (*goparent.SleepSummary, error) {
	if m.StatErr != nil {
		return nil, m.StatErr
	}
	if m.Stat != nil {
		return m.Stat, nil
	}

	return nil, nil
}

//Status -
func (m *SleepService) Status(*goparent.Family, *goparent.Child) (bool, error) {
	panic("not implemented")
}

//Start -
func (m *SleepService) Start(*goparent.Sleep, *goparent.Family, *goparent.Child) error {
	panic("not implemented")
}

//End -
func (m *SleepService) End(*goparent.Sleep, *goparent.Family, *goparent.Child) error {
	panic("not implemented")
}

//GraphData -
func (m *SleepService) GraphData(*goparent.Child) (*goparent.SleepChartData, error) {
	panic("not implemented")
}
