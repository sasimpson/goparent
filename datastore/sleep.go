package datastore

import (
	"context"

	"github.com/sasimpson/goparent"
)

//SleepService -
type SleepService struct {
	Env *goparent.Env
}

func (s *SleepService) Save(context.Context, *goparent.Sleep) error {
	panic("not implemented")
}

func (s *SleepService) Sleep(context.Context, *goparent.Family, uint64) ([]*goparent.Sleep, error) {
	panic("not implemented")
}

func (s *SleepService) Stats(context.Context, *goparent.Child) (*goparent.SleepSummary, error) {
	panic("not implemented")
}

func (s *SleepService) Status(context.Context, *goparent.Family, *goparent.Child) (bool, error) {
	panic("not implemented")
}

func (s *SleepService) Start(context.Context, *goparent.Sleep, *goparent.Family, *goparent.Child) error {
	panic("not implemented")
}

func (s *SleepService) End(context.Context, *goparent.Sleep, *goparent.Family, *goparent.Child) error {
	panic("not implemented")
}

func (s *SleepService) GraphData(context.Context, *goparent.Child) (*goparent.SleepChartData, error) {
	panic("not implemented")
}
