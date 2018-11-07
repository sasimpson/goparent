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
	var feedings []*goparent.Feeding
	familyKey := datastore.NewKey(ctx, FamilyKind, family.ID, 0, nil)

	daysBack := int(0 - days)
	start := time.Now().AddDate(0, 0, daysBack)

	q := datastore.NewQuery(FeedingKind).Ancestor(familyKey).Filter("TimeStamp > ", start).Order("-TimeStamp")
	itx := q.Run(ctx)
	for {
		var feeding goparent.Feeding
		_, err := itx.Next(&feeding)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		feedings = append(feedings, &feeding)
	}
	return feedings, nil
}

//Stats -
func (s *FeedingService) Stats(ctx context.Context, child *goparent.Child) (*goparent.FeedingSummary, error) {
	var feedings []goparent.Feeding
	end := time.Now()
	start := end.AddDate(0, 0, -1)

	familyKey := datastore.NewKey(ctx, FamilyKind, child.FamilyID, 0, nil)
	childKey := datastore.NewKey(ctx, ChildKind, child.ID, 0, familyKey)

	q := datastore.NewQuery(FeedingKind).Filter("ChildID = ", childKey).Filter("TimeStamp >=", start).Order("-TimeStamp")
	itx := q.Run(ctx)
	for {
		var feeding goparent.Feeding
		_, err := itx.Next(&feeding)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		feedings = append(feedings, feeding)
	}

	summary := &goparent.FeedingSummary{
		Data:  feedings,
		Total: make(map[string]float32),
		Mean:  make(map[string]float32),
		Range: make(map[string]int),
	}

	for _, x := range feedings {
		if _, ok := summary.Total[x.Type]; !ok {
			summary.Total[x.Type] = 0.0
		}
		summary.Total[x.Type] += x.Amount
		summary.Range[x.Type]++
	}
	for k := range summary.Total {
		summary.Mean[k] = summary.Total[k] / float32(summary.Range[k])
	}
	return summary, nil
}

//GraphData -
func (s *FeedingService) GraphData(ctx context.Context, child *goparent.Child) (*goparent.FeedingChartData, error) {
	var feedings []goparent.Feeding
	feedingCounts := make(map[time.Time][]goparent.Feeding)

	end := time.Now()
	start := end.AddDate(0, 0, -7)
	q := datastore.NewQuery(FeedingKind).Filter("ChildID =", child.ID).Filter("TimeStamp >", start).Filter("TimeStamp <=", end).Order("-TimeStamp")
	//get each item from the query organize them by day.
	itx := q.Run(ctx)
	for {
		var feeding goparent.Feeding
		_, err := itx.Next(&feeding)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		roundedDate := RoundToDay(feeding.TimeStamp, false)
		feedingCounts[roundedDate] = append(feedingCounts[roundedDate], feeding)
		feedings = append(feedings, feeding)
	}

	chartData := &goparent.FeedingChartData{
		Start:   start,
		End:     end,
		Dataset: make([]goparent.FeedingChartDataset, 1),
	}
	for day, feedings := range feedingCounts {
		counts := make(map[string]int)
		for _, t := range feedings {
			counts[t.Type]++
		}
		for feedingType, count := range counts {
			chartData.Dataset = append(chartData.Dataset, goparent.FeedingChartDataset{
				Date:  day,
				Type:  feedingType,
				Count: count,
			})
		}
	}
	return chartData, nil
}
