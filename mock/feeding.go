package mock

import (
	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/config"
)

type MockFeedingService struct {
	Env *config.Env
}

func (mfs *MockFeedingService) Save(feeding *goparent.Feeding) error {
	return nil
}

func (mfs *MockFeedingService) Feeding(family *goparent.Family) ([]*goparent.Feeding, error) {
	return nil, nil
}

func (mfs *MockFeedingService) Stats(child *goparent.Child) (*goparent.FeedingSummary, error) {
	return nil, nil
}
