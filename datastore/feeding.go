package datastore

import (
	"context"

	"github.com/sasimpson/goparent"
)

//FeedingService -
type FeedingService struct {
	Env *goparent.Env
}

//FeedingKind is the constant for the feeding entity kind in gcp datastore
const FeedingKind = "Feeding"

//Save -
func (s *FeedingService) Save(ctx context.Context, feeding *goparent.Feeding) error {
	// var feedKey *dataStore.feedKey
	// familyKey := datastore.NewKey(ctx, FamilyKind)
	panic("not implemented")
}

//Feeding -
func (s *FeedingService) Feeding(ctx context.Context, family *goparent.Family, days uint64) ([]*goparent.Feeding, error) {
	panic("not implemented")
}

//Stats -
func (s *FeedingService) Stats(ctx context.Context, child *goparent.Child) (*goparent.FeedingSummary, error) {
	panic("not implemented")
}

//GraphData -
func (s *FeedingService) GraphData(ctx context.Context, child *goparent.Child) (*goparent.FeedingChartData, error) {
	panic("not implemented")
}
