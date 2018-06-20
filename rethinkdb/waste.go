package rethinkdb

import (
	"fmt"
	"time"

	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/config"
	gorethink "gopkg.in/gorethink/gorethink.v3"
)

//WasteService - struct for implmenting the interface
type WasteService struct {
	Env *config.Env
}

type graphData struct {
	Group     []int `gorethink:"group"`
	Reduction int   `gorethink:"reduction"`
}

//Save - save waste data
func (ws *WasteService) Save(waste *goparent.Waste) error {
	session, err := ws.Env.DB.GetConnection()
	if err != nil {
		return err
	}

	res, err := gorethink.Table("waste").Insert(waste, gorethink.InsertOpts{Conflict: "replace"}).RunWrite(session)
	if err != nil {
		return err
	}

	if res.Inserted > 0 {
		waste.ID = res.GeneratedKeys[0]
	}
	return nil
}

//Waste - get all waste by user and child id.
func (ws *WasteService) Waste(family *goparent.Family, days uint64) ([]*goparent.Waste, error) {
	session, err := ws.Env.DB.GetConnection()
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
		OrderBy(gorethink.Desc("timestamp")).Run(session)
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
func (ws *WasteService) Stats(child *goparent.Child) (*goparent.WasteSummary, error) {
	session, err := ws.Env.DB.GetConnection()
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
		Run(session)
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

func (ws *WasteService) GraphData(child *goparent.Child) (*[]goparent.WasteGraphData, error) {
	session, err := ws.Env.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	end := time.Now()
	start := end.AddDate(0, 0, -7)

	res, err := gorethink.Table("waste").
		Filter(gorethink.Row.Field("timestamp").During(start, end)).
		Group(
			gorethink.Row.Field("timestamp").Year(),
			gorethink.Row.Field("timestamp").Month(),
			gorethink.Row.Field("timestamp").Day(),
			gorethink.Row.Field("wasteType")).
		Count().Run(session)
	/*
		r.db("goparent")
			.table("waste")
			.filter(r.row("timestamp")
				.during(r.time(2018,6,7,"Z"),r.now()))
			.group([r.row("timestamp").year(), r.row("timestamp").month(),r.row("timestamp").day(), r.row("wasteType")])
			.count()
	*/
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var data []graphData
	err = res.All(&data)
	if err != nil {
		return nil, err
	}

	var wasteGraphData []goparent.WasteGraphData
	for _, line := range data {
		var wgd goparent.WasteGraphData
		wgdDate, err := time.Parse("2006-01-02", fmt.Sprintf("%d-%02d-%02d", line.Group[0], line.Group[1], line.Group[2]))
		if err != nil {
			return nil, err
		}
		wgd.Date = wgdDate
		wgd.Count = line.Reduction
		wgd.Type = line.Group[3]
		wasteGraphData = append(wasteGraphData, wgd)
	}

	return &wasteGraphData, nil

}
