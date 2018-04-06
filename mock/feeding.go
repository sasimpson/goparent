package mock

import (
	"log"

	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/config"
)

type MockFeedingService struct {
	Env      *config.Env
	Feedings []*goparent.Feeding
	GetErr   error
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
	log.Println("mock feeding returning ", mfs.Feedings)
	return mfs.Feedings, nil
}

func (mfs *MockFeedingService) Stats(child *goparent.Child) (*goparent.FeedingSummary, error) {
	return nil, nil
}
