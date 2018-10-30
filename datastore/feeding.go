package datastore

import (
	"context"
	"time"

	"google.golang.org/appengine/datastore"

	"github.com/google/uuid"
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
	var feedKey *datastore.Key
	familyKey := datastore.NewKey(ctx, FamilyKind, feeding.FamilyID, 0, nil)
	childKey := datastore.NewKey(ctx, ChildKind, feeding.ChildID, 0, familyKey)
	if feeding.ID == "" {
		u := uuid.New()
		feedKey = datastore.NewKey(ctx, FeedingKind, u.String(), 0, childKey)
		feeding.CreatedAt = time.Now()
		feeding.LastUpdated = feeding.CreatedAt
		feeding.ID = u.String()
	} else {
		feedKey = datastore.NewKey(ctx, FeedingKind, feeding.ID, 0, childKey)
		feeding.LastUpdated = time.Now()
	}

	_, err := datastore.Put(ctx, feedKey, feeding)
	if err != nil {
		return NewError("FeedingService.Save", err)
	}
	return nil
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
