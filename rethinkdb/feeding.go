package rethinkdb

import (
	"context"
	"fmt"
	"time"

	"github.com/sasimpson/goparent"
	gorethink "gopkg.in/gorethink/gorethink.v3"
)

//FeedingService - struct for implementing interface
type FeedingService struct {
	Env *goparent.Env
	DB  *DBEnv
}

type feedingReductionData struct {
	Group     []interface{} `gorethink:"group"`
	Reduction []struct {
		FeedingAmount float32 `gorethink:"feedingAmount"`
	} `gorethink:"reduction"`
}

//Save - save the structure to the datastore
func (fs *FeedingService) Save(ctx context.Context, feeding *goparent.Feeding) error {
	err := fs.DB.GetConnection()
	if err != nil {
		return err
	}

	res, err := gorethink.Table("feeding").Insert(feeding, gorethink.InsertOpts{Conflict: "replace"}).RunWrite(fs.DB.Session)
	if err != nil {
		return err
	}

	if res.Inserted > 0 {
		feeding.ID = res.GeneratedKeys[0]
	}
	return nil
}

//Feeding - get all records for a user from the datastore
func (fs *FeedingService) Feeding(ctx context.Context, family *goparent.Family, days uint64) ([]*goparent.Feeding, error) {
	err := fs.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	daysBack := int(0 - days)

	res, err := gorethink.Table("feeding").
		Filter(map[string]interface{}{
			"familyID": family.ID,
		}).
		Filter(gorethink.Row.Field("timestamp").During(time.Now().AddDate(0, 0, daysBack), time.Now())).
		OrderBy(gorethink.Desc("timestamp")).
		Run(fs.DB.Session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var rows []*goparent.Feeding
	err = res.All(&rows)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

//Stats - get feeding stats for one child for the last 24 hours.
func (fs *FeedingService) Stats(ctx context.Context, child *goparent.Child) (*goparent.FeedingSummary, error) {
	err := fs.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	end := time.Now()
	start := end.AddDate(0, 0, -1)

	res, err := gorethink.Table("feeding").
		Filter(map[string]interface{}{
			"childID": child.ID,
		}).
		Filter(
			gorethink.Row.Field("timestamp").During(start, end),
		).
		OrderBy(gorethink.Desc("timestamp")).
		Run(fs.DB.Session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var rows []goparent.Feeding
	err = res.All(&rows)
	if err != nil {
		return nil, err
	}
	//build summary output
	summary := &goparent.FeedingSummary{
		Data:  rows,
		Total: make(map[string]float32),
		Mean:  make(map[string]float32),
		Range: make(map[string]int),
	}

	for _, x := range rows {
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
func (fs *FeedingService) GraphData(ctx context.Context, child *goparent.Child) (*goparent.FeedingChartData, error) {
	err := fs.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	end := time.Now()
	start := end.AddDate(0, 0, -7)
	/*
		r.db("goparent")
		.table("feeding")
		.filter(r.row("timestamp")
			.during(r.time(2018,7,13,"Z"),r.now()))
		.group([r.row("timestamp").year(), r.row("timestamp").month(), r.row("timestamp").day(), r.row("feedingType")])
		.pluck( "feedingAmount")
	*/
	res, err := gorethink.Table("feeding").
		Filter(gorethink.Row.Field("timestamp").During(start, end)).OrderBy("timestamp").
		Group(
			gorethink.Row.Field("timestamp").Year(),
			gorethink.Row.Field("timestamp").Month(),
			gorethink.Row.Field("timestamp").Day(),
			gorethink.Row.Field("feedingType"),
		).
		Pluck("feedingAmount").
		Run(fs.DB.Session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var data []feedingReductionData
	err = res.All(&data)
	if err != nil {
		return nil, err
	}

	chartData := &goparent.FeedingChartData{Start: start, End: end}
	// graph.Data = goparent.ChartData{Datasets: []goparent.ChartDataset{}}
	for _, line := range data {
		gdDate, err := time.Parse("2006-01-02", fmt.Sprintf("%.0f-%02.0f-%02.0f", line.Group[0], line.Group[1], line.Group[2]))
		if err != nil {
			return nil, err
		}

		var rC int
		var rS float32
		for _, reduction := range line.Reduction {
			rC++
			rS += reduction.FeedingAmount
		}
		dataset := goparent.FeedingChartDataset{Date: gdDate, Type: line.Group[3].(string), Count: rC, Sum: rS}
		chartData.Dataset = append(chartData.Dataset, dataset)
	}
	return chartData, nil
}
