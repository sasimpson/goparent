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
		sleep.ID = u.String()
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

//Sleep gives back an array of sleep instances for the number of days back from today.
func (s *SleepService) Sleep(ctx context.Context, family *goparent.Family, days uint64) ([]*goparent.Sleep, error) {
	var sleeps []*goparent.Sleep
	familyKey := datastore.NewKey(ctx, FamilyKind, family.ID, 0, nil)

	daysBack := int(0 - days)
	start := time.Now().AddDate(0, 0, daysBack)
	q := datastore.NewQuery(SleepKind).Ancestor(familyKey).Filter("Start > ", start).Order("-Start")
	itx := q.Run(ctx)
	for {
		var sleep goparent.Sleep
		_, err := itx.Next(&sleep)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		sleeps = append(sleeps, &sleep)
	}
	return sleeps, nil
}

//Status should return the sleep status of a child, ie if they are asleep or not
func (s *SleepService) Status(context.Context, *goparent.Family, *goparent.Child) (bool, error) {
	panic("not implemented")
}

//Start will start a sleep session or error if there is a current one active.
func (s *SleepService) Start(context.Context, *goparent.Sleep, *goparent.Family, *goparent.Child) error {
	panic("not implemented")
}

//End will end a current sleep session or error if there isn't one
func (s *SleepService) End(context.Context, *goparent.Sleep, *goparent.Family, *goparent.Child) error {
	panic("not implemented")
}

//Stats returns sleep stats about a one day period for a child.
func (s *SleepService) Stats(ctx context.Context, child *goparent.Child) (*goparent.SleepSummary, error) {
	var sleeps []goparent.Sleep
	end := time.Now()
	start := end.AddDate(0, 0, -1)

	familyKey := datastore.NewKey(ctx, FamilyKind, child.FamilyID, 0, nil)
	childKey := datastore.NewKey(ctx, ChildKind, child.ID, 0, familyKey)

	q := datastore.NewQuery(SleepKind).Filter("ChildID = ", childKey).Filter("TimeStamp >= ", start).Order("-TimeStamp")
	itx := q.Run(ctx)
	for {
		var sleep goparent.Sleep
		_, err := itx.Next(&sleep)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		sleeps = append(sleeps, sleep)
	}

	summary := &goparent.SleepSummary{
		Data:  sleeps,
		Total: 0,
		Mean:  0.0,
		Range: 0,
	}

	for _, x := range sleeps {
		//if the sleep end is before the start then the sleep period hasn't stopped yet.  don't count
		if x.End.After(x.Start) {
			summary.Total += int64(x.End.Sub(x.Start).Seconds())
			summary.Range++
		}
	}
	// summary.Mean = float64(summary.Total / int64(summary.Range))

	return summary, nil
}

//GraphData returns data that a graph can be created from
func (s *SleepService) GraphData(ctx context.Context, child *goparent.Child) (*goparent.SleepChartData, error) {
	var sleeps []goparent.Sleep
	sleepCounts := make(map[time.Time][]goparent.Sleep)

	end := time.Now()
	start := end.AddDate(0, 0, -7)
	q := datastore.NewQuery(SleepKind).Filter("ChildID =", child.ID).Filter("TimeStamp >", start).Filter("TimeStamp <=", end).Order("-TimeStamp")
	//get each item from the query organize them by day.
	itx := q.Run(ctx)
	for {
		var sleep goparent.Sleep
		_, err := itx.Next(&sleep)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		roundedDate := RoundToDay(sleep.Start, false)
		sleepCounts[roundedDate] = append(sleepCounts[roundedDate], sleep)
		sleeps = append(sleeps, sleep)
	}

	chartData := &goparent.SleepChartData{
		Start:   start,
		End:     end,
		Dataset: make([]goparent.SleepChartDataset, 1),
	}

	//sleep data needs to be represented as each increment of sleep per day.

	for day, sleeps := range sleepCounts {
		var sleepTotals []time.Duration
		for _, sleepEntry := range sleeps {
			//if the entry hasn't stopped, don't count it, its still active
			if sleepEntry.End.After(sleepEntry.Start) {
				sleepTotals = append(sleepTotals, sleepEntry.End.Sub(sleepEntry.Start))
			}
		}
		chartData.Dataset = append(chartData.Dataset, goparent.SleepChartDataset{
			Date:   day,
			Totals: sleepTotals,
		})
	}

	return chartData, nil
}
