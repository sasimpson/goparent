package rethinkdb

import (
	"context"
	"fmt"
	"time"

	"github.com/sasimpson/goparent"
	gorethink "gopkg.in/gorethink/gorethink.v3"
)

//WasteService - struct for implmenting the interface
type WasteService struct {
	Env *goparent.Env
	DB  *DBEnv
}

type reductionData struct {
	Group     []int `gorethink:"group"`
	Reduction int   `gorethink:"reduction"`
}

//Save - save waste data
func (ws *WasteService) Save(ctx context.Context, waste *goparent.Waste) error {
	err := ws.DB.GetConnection()
	if err != nil {
		return err
	}

	res, err := gorethink.Table("waste").Insert(waste, gorethink.InsertOpts{Conflict: "replace"}).RunWrite(ws.DB.Session)
	if err != nil {
		return err
	}

	if res.Inserted > 0 {
		waste.ID = res.GeneratedKeys[0]
	}
	return nil
}

//Waste - get all waste by user and child id.
func (ws *WasteService) Waste(ctx context.Context, family *goparent.Family, days uint64) ([]*goparent.Waste, error) {
	err := ws.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	daysBack := int(0 - days)
	res, err := gorethink.Table("waste").
		Filter(
			map[string]interface{}{
				"familyID": family.ID,
			}).
		Filter(gorethink.Row.Field("timestamp").During(time.Now().AddDate(0, 0, daysBack), time.Now())).
		OrderBy(gorethink.Desc("timestamp")).Run(ws.DB.Session)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	var rows []*goparent.Waste
	err = res.All(&rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

//Stats - get waste stats for one child for the last 24 hours.
func (ws *WasteService) Stats(ctx context.Context, child *goparent.Child) (*goparent.WasteSummary, error) {
	err := ws.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	end := time.Now()
	start := end.AddDate(0, 0, -1)

	res, err := gorethink.Table("waste").
		Filter(map[string]interface{}{
			"childID": child.ID,
		}).
		Filter(gorethink.Row.Field("timestamp").During(start, end)).
		OrderBy(gorethink.Desc("timestamp")).
		Run(ws.DB.Session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var rows []goparent.Waste
	err = res.All(&rows)
	if err != nil {
		return nil, err
	}

	//build summary output
	summary := goparent.WasteSummary{
		Data:  rows,
		Total: make(map[int]int),
	}

	for _, x := range rows {
		if _, ok := summary.Total[x.Type]; !ok {
			summary.Total[x.Type] = 0
		}
		summary.Total[x.Type]++
	}
	return &summary, nil
}

//GraphData -
func (ws *WasteService) GraphData(ctx context.Context, child *goparent.Child) (*goparent.WasteChartData, error) {
	err := ws.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	end := time.Now()
	start := end.AddDate(0, 0, -7)
	/*
			r.db("goparent")
			.table("waste")
			.filter(r.row("timestamp")
				.during(r.time(2018,6,7,"Z"),r.now()))
			.group([r.row("timestamp").year(), r.row("timestamp").month(), r.row("timestamp").day()])
			.map(function(row){
				return row("wasteType")
		  })
	*/
	res, err := gorethink.Table("waste").
		Filter(gorethink.Row.Field("timestamp").During(start, end)).OrderBy("timestamp").
		Group(
			gorethink.Row.Field("timestamp").Year(),
			gorethink.Row.Field("timestamp").Month(),
			gorethink.Row.Field("timestamp").Day(),
			gorethink.Row.Field("wasteType"),
		).Count().
		Run(ws.DB.Session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var data []reductionData
	err = res.All(&data)
	if err != nil {
		return nil, err
	}

	chartData := &goparent.WasteChartData{Start: start, End: end, Dataset: make([]goparent.WasteChartDataset, 1)}
	// graph.Data = goparent.ChartData{Datasets: []goparent.ChartDataset{}}
	for _, line := range data {
		gdDate, err := time.Parse("2006-01-02", fmt.Sprintf("%d-%02d-%02d", line.Group[0], line.Group[1], line.Group[2]))
		if err != nil {
			return nil, err
		}

		dataset := goparent.WasteChartDataset{Date: gdDate, Type: line.Group[3], Count: line.Reduction}
		chartData.Dataset = append(chartData.Dataset, dataset)
	}
	return chartData, nil
}
