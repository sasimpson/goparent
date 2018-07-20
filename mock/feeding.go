package mock

import (
	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/config"
)

type MockFeedingService struct {
	Env      *config.Env
	Feedings []*goparent.Feeding
	Stat     *goparent.FeedingSummary
	Graph    *goparent.FeedingChartData
	GetErr   error
	StatErr  error
	GraphErr error
}

func (m *MockFeedingService) Save(feeding *goparent.Feeding) error {
	if m.GetErr != nil {
		return m.GetErr
	}
	return nil
}

func (m *MockFeedingService) Feeding(family *goparent.Family) ([]*goparent.Feeding, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	if m.Feedings != nil {
		return m.Feedings, nil
	}

	return nil, nil
}

func (m *MockFeedingService) Stats(child *goparent.Child) (*goparent.FeedingSummary, error) {
	if m.StatErr != nil {
		return nil, m.StatErr
	}
	if m.Stat != nil {
		return m.Stat, nil
	}

	return nil, nil
}

func (m *MockFeedingService) GraphData(*goparent.Child) (*goparent.FeedingChartData, error) {
	if m.GraphErr != nil {
		return nil, m.GraphErr
	}
	if m.Graph != nil {
		return m.Graph, nil
	}
	return nil, nil
}
