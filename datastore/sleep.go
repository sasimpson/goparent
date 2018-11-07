package datastore

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sasimpson/goparent"
	"google.golang.org/appengine/datastore"
)

//SleepService -
type SleepService struct {
	Env *goparent.Env
}

//SleepKind is the constant for the sleep entity kind in GCP datastore
const SleepKind = "Sleep"

//Save is the function that will save a record
func (s *SleepService) Save(ctx context.Context, sleep *goparent.Sleep) error {
	var sleepKey *datastore.Key
	familyKey := datastore.NewKey(ctx, FamilyKind, sleep.FamilyID, 0, nil)
	childKey := datastore.NewKey(ctx, ChildKind, sleep.ChildID, 0, familyKey)
	if sleep.ID == "" {
		u := uuid.New()
		sleepKey = datastore.NewKey(ctx, SleepKind, u.String(), 0, childKey)
		sleep.CreatedAt = time.Now()
		sleep.LastUpdated = sleep.CreatedAt
	} else {
		sleepKey = datastore.NewKey(ctx, SleepKind, sleep.ID, 0, childKey)
		sleep.LastUpdated = time.Now()
	}

	_, err := datastore.Put(ctx, sleepKey, sleep)
	if err != nil {
		return NewError("SleepService.Save", err)
	}
	return nil
}

func (s *SleepService) Sleep(context.Context, *goparent.Family, uint64) ([]*goparent.Sleep, error) {
	panic("not implemented")
}

func (s *SleepService) Stats(context.Context, *goparent.Child) (*goparent.SleepSummary, error) {
	panic("not implemented")
}

func (s *SleepService) Status(context.Context, *goparent.Family, *goparent.Child) (bool, error) {
	panic("not implemented")
}

func (s *SleepService) Start(context.Context, *goparent.Sleep, *goparent.Family, *goparent.Child) error {
	panic("not implemented")
}

func (s *SleepService) End(context.Context, *goparent.Sleep, *goparent.Family, *goparent.Child) error {
	panic("not implemented")
}

func (s *SleepService) GraphData(context.Context, *goparent.Child) (*goparent.SleepChartData, error) {
	panic("not implemented")
}
