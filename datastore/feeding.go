package datastore

import "github.com/sasimpson/goparent"

//FeedingService -
type FeedingService struct {
}

//Save -
func (s *FeedingService) Save(*goparent.Feeding) error {
	panic("not implemented")
}

//Feeding -
func (s *FeedingService) Feeding(*goparent.Family, uint64) ([]*goparent.Feeding, error) {
	panic("not implemented")
}

//Stats -
func (s *FeedingService) Stats(*goparent.Child) (*goparent.FeedingSummary, error) {
	panic("not implemented")
}

//GraphData -
func (s *FeedingService) GraphData(*goparent.Child) (*goparent.FeedingChartData, error) {
	panic("not implemented")
}
