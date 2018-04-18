package mock

import (
	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/config"
)

type MockFeedingService struct {
	Env      *config.Env
	Feedings []*goparent.Feeding
	Stat     *goparent.FeedingSummary
	GetErr   error
	StatErr  error
}

func (mfs *MockFeedingService) Save(feeding *goparent.Feeding) error {
	if mfs.GetErr != nil {
		return mfs.GetErr
	}
	return nil
}

func (mfs *MockFeedingService) Feeding(family *goparent.Family) ([]*goparent.Feeding, error) {
	if mfs.GetErr != nil {
		return nil, mfs.GetErr
	}
	if mfs.Feedings != nil {
		return mfs.Feedings, nil
	}

	return nil, nil
}

func (mfs *MockFeedingService) Stats(child *goparent.Child) (*goparent.FeedingSummary, error) {
	if mfs.StatErr != nil {
		return nil, mfs.StatErr
	}
	if mfs.Stat != nil {
		return mfs.Stat, nil
	}

	return nil, nil
}
