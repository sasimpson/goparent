package datastore

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sasimpson/goparent"
	"google.golang.org/appengine/datastore"
)

//WasteService -
type WasteService struct {
	Env *goparent.Env
}

//WasteKind is the constant for the waste entity kind in gcp datastore
const WasteKind = "Waste"

//Save the waste entry
func (s *WasteService) Save(ctx context.Context, waste *goparent.Waste) error {
	var wasteKey *datastore.Key
	familyKey := datastore.NewKey(ctx, FamilyKind, waste.FamilyID, 0, nil)
	childKey := datastore.NewKey(ctx, ChildKind, waste.ChildID, 0, familyKey)
	if waste.ID == "" {
		u := uuid.New()
		wasteKey = datastore.NewKey(ctx, WasteKind, u.String(), 0, childKey)
		waste.CreatedAt = time.Now()
		waste.LastUpdated = waste.CreatedAt
		waste.ID = u.String()
	} else {
		wasteKey = datastore.NewKey(ctx, WasteKind, waste.ID, 0, childKey)
		waste.LastUpdated = time.Now()
	}

	_, err := datastore.Put(ctx, wasteKey, waste)
	if err != nil {
		return NewError("WasteService.Save", err)
	}
	return nil
}

//Waste returns all waste entries by user and child id?
func (s *WasteService) Waste(ctx context.Context, family *goparent.Family, days uint64) ([]*goparent.Waste, error) {
	panic("not implemented")
}

//Stats -
func (s *WasteService) Stats(ctx context.Context, child *goparent.Child) (*goparent.WasteSummary, error) {
	panic("not implemented")
}

//GraphData -
func (s *WasteService) GraphData(ctx context.Context, child *goparent.Child) (*goparent.WasteChartData, error) {
	panic("not implemented")
}
