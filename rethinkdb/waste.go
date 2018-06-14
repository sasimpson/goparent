package rethinkdb

import (
	"log"
	"time"

	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/config"
	gorethink "gopkg.in/gorethink/gorethink.v3"
)

//WasteService - struct for implmenting the interface
type WasteService struct {
	Env *config.Env
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
	log.Printf("days back: %d", daysBack)
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

func (ws *WasteService) GraphData(child *goparent.Child) error {
	// session, err := ws.Env.DB.GetConnection()
	// if err != nil {
	// 	return err
	// }

	/*
		r.db("goparent")
			.table("waste")
			.filter(r.row("timestamp")
				.during(r.time(2018,6,7,"Z"),r.now()))
			.group([r.row("timestamp").year(), r.row("timestamp").month(),r.row("timestamp").day(), r.row("wasteType")])
			.count()
	*/

	return nil

}
