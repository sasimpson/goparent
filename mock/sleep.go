package mock

import (
	"context"

	"github.com/sasimpson/goparent"
)

//SleepService -
type SleepService struct {
	Env     *goparent.Env
	Sleeps  []*goparent.Sleep
	Stat    *goparent.SleepSummary
	GetErr  error
	StatErr error
}

//Save -
func (m *SleepService) Save(context.Context, *goparent.Sleep) error {
	if m.GetErr != nil {
		return m.GetErr
	}
	return nil
}

//Sleep -
func (m *SleepService) Sleep(context.Context, *goparent.Family, uint64) ([]*goparent.Sleep, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	if m.Sleeps != nil {
		return m.Sleeps, nil
	}
	return nil, nil
}

//Stats -
func (m *SleepService) Stats(context.Context, *goparent.Child) (*goparent.SleepSummary, error) {
	if m.StatErr != nil {
		return nil, m.StatErr
	}
	if m.Stat != nil {
		return m.Stat, nil
	}

	return nil, nil
}

//Status -
func (m *SleepService) Status(context.Context, *goparent.Family, *goparent.Child) (bool, error) {
	panic("not implemented")
}

//Start -
func (m *SleepService) Start(context.Context, *goparent.Sleep, *goparent.Family, *goparent.Child) error {
	panic("not implemented")
}

//End -
func (m *SleepService) End(context.Context, *goparent.Sleep, *goparent.Family, *goparent.Child) error {
	panic("not implemented")
}

//GraphData -
func (m *SleepService) GraphData(context.Context, *goparent.Child) (*goparent.SleepChartData, error) {
	panic("not implemented")
}
