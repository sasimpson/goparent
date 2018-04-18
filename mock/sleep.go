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

func (mss *MockSleepService) Save(*goparent.Sleep) error {
	panic("not implemented")
}

func (mss *MockSleepService) Sleep(*goparent.Family) ([]*goparent.Sleep, error) {
	panic("not implemented")
}

func (mss *MockSleepService) Stats(*goparent.Child) (*goparent.SleepSummary, error) {
	if mss.StatErr != nil {
		return nil, mss.StatErr
	}
	if mss.Stat != nil {
		return mss.Stat, nil
	}

	return nil, nil
}

func (mss *MockSleepService) Status(*goparent.Family, *goparent.Child) (bool, error) {
	panic("not implemented")
}

func (mss *MockSleepService) Start(*goparent.Sleep, *goparent.Family, *goparent.Child) error {
	panic("not implemented")
}

func (mss *MockSleepService) End(*goparent.Sleep, *goparent.Family, *goparent.Child) error {
	panic("not implemented")
}
