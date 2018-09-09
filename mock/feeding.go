package mock

import "github.com/sasimpson/goparent"

//FeedingService -
type FeedingService struct {
	Env      *goparent.Env
	Feedings []*goparent.Feeding
	Stat     *goparent.FeedingSummary
	Graph    *goparent.FeedingChartData
	GetErr   error
	StatErr  error
	GraphErr error
}

//Save -
func (m *FeedingService) Save(feeding *goparent.Feeding) error {
	if m.GetErr != nil {
		return m.GetErr
	}
	return nil
}

//Feeding -
func (m *FeedingService) Feeding(family *goparent.Family, days uint64) ([]*goparent.Feeding, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	if m.Feedings != nil {
		return m.Feedings, nil
	}

	return nil, nil
}

//Stats -
func (m *FeedingService) Stats(child *goparent.Child) (*goparent.FeedingSummary, error) {
	if m.StatErr != nil {
		return nil, m.StatErr
	}
	if m.Stat != nil {
		return m.Stat, nil
	}

	return nil, nil
}

//GraphData -
func (m *FeedingService) GraphData(*goparent.Child) (*goparent.FeedingChartData, error) {
	if m.GraphErr != nil {
		return nil, m.GraphErr
	}
	if m.Graph != nil {
		return m.Graph, nil
	}
	return nil, nil
}
