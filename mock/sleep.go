package mock

import (
	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/config"
)

type MockSleepService struct {
	Env     *config.Env
	Sleeps  []*goparent.Sleep
	Stat    *goparent.SleepSummary
	GetErr  error
	StatErr error
}

func (m *MockSleepService) Save(*goparent.Sleep) error {
	if m.GetErr != nil {
		return m.GetErr
	}
	return nil
}

func (m *MockSleepService) Sleep(*goparent.Family) ([]*goparent.Sleep, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	if m.Sleeps != nil {
		return m.Sleeps, nil
	}
	return nil, nil
}

func (m *MockSleepService) Stats(*goparent.Child) (*goparent.SleepSummary, error) {
	if m.StatErr != nil {
		return nil, m.StatErr
	}
	if m.Stat != nil {
		return m.Stat, nil
	}

	return nil, nil
}

func (m *MockSleepService) Status(*goparent.Family, *goparent.Child) (bool, error) {
	panic("not implemented")
}

func (m *MockSleepService) Start(*goparent.Sleep, *goparent.Family, *goparent.Child) error {
	panic("not implemented")
}

func (m *MockSleepService) End(*goparent.Sleep, *goparent.Family, *goparent.Child) error {
	panic("not implemented")
}
