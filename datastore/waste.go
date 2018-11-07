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
	var wastes []*goparent.Waste
	familyKey := datastore.NewKey(ctx, FamilyKind, family.ID, 0, nil)

	daysBack := int(0 - days)
	start := time.Now().AddDate(0, 0, daysBack)

	q := datastore.NewQuery(WasteKind).Ancestor(familyKey).Filter("TimeStamp > ", start).Order("-TimeStamp")
	itx := q.Run(ctx)
	for {
		var waste goparent.Waste
		_, err := itx.Next(&waste)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		wastes = append(wastes, &waste)
	}
	return wastes, nil
}

//Stats returns the stats for a particular child
func (s *WasteService) Stats(ctx context.Context, child *goparent.Child) (*goparent.WasteSummary, error) {
	var wastes []goparent.Waste
	end := time.Now()
	start := end.AddDate(0, 0, -1)

	familyKey := datastore.NewKey(ctx, FamilyKind, child.FamilyID, 0, nil)
	childKey := datastore.NewKey(ctx, ChildKind, child.ID, 0, familyKey)

	q := datastore.NewQuery(WasteKind).Filter("ChildID=", childKey).Filter("TimeStamp >=", start).Order("-TimeStamp")
	itx := q.Run(ctx)
	for {
		var waste goparent.Waste
		_, err := itx.Next(&waste)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		wastes = append(wastes, waste)
	}

	summary := goparent.WasteSummary{
		Data:  wastes,
		Total: make(map[int]int),
	}

	for _, x := range wastes {
		if _, ok := summary.Total[x.Type]; !ok {
			summary.Total[x.Type] = 0
		}
		summary.Total[x.Type]++
	}
	return &summary, nil
}

//GraphData returns the data necessary for graphing information about the child
func (s *WasteService) GraphData(ctx context.Context, child *goparent.Child) (*goparent.WasteChartData, error) {
	var wastes []goparent.Waste
	wasteCounts := make(map[time.Time][]goparent.Waste)
	end := time.Now()
	start := end.AddDate(0, 0, -7)
	q := datastore.NewQuery(WasteKind).Filter("ChildID =", child.ID).Filter("TimeStamp >", start).Filter("TimeStamp <=", end).Order("-TimeStamp")
	//get each item from the query organize them by day.
	itx := q.Run(ctx)
	for {
		var waste goparent.Waste
		_, err := itx.Next(&waste)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		roundedDate := RoundToDay(waste.TimeStamp, false)
		wasteCounts[roundedDate] = append(wasteCounts[roundedDate], waste)
		wastes = append(wastes, waste)
	}

	chartData := &goparent.WasteChartData{
		Start:   start,
		End:     end,
		Dataset: make([]goparent.WasteChartDataset, 1),
	}

	//now organize each day by the total of each type. setup dataset
	for day, wastes := range wasteCounts {
		counts := make(map[int]int)
		for _, t := range wastes {
			counts[t.Type]++
		}
		for wasteType, count := range counts {
			chartData.Dataset = append(chartData.Dataset, goparent.WasteChartDataset{
				Date:  day,
				Type:  wasteType,
				Count: count,
			})
		}
	}

	return chartData, nil
}
