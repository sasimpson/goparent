package models

import (
	"time"

	"github.com/sasimpson/goparent/config"
	gorethink "gopkg.in/gorethink/gorethink.v3"
)

//Waste - structure for holding waste data such as diapers
type Waste struct {
	ID        string    `json:"id" gorethink:"id,omitempty"`
	Type      int       `json:"wasteType" gorethink:"wasteType"`
	Notes     string    `json:"notes" gorethink:"notes"`
	UserID    string    `json:"userid" gorethink:"userID"`
	FamilyID  string    `json:"familyid" gorethink:"familyID"`
	ChildID   string    `json:"childid" gorethink:"childID"`
	TimeStamp time.Time `json:"timestamp" gorethink:"timestamp"`
}

//WasteSummary - structure for waste summary data
type WasteSummary struct {
	Data  []Waste     `json:"data"`
	Total map[int]int `json:"total"`
}

//WasteType - the type of waste, solid, liquid, solid & liquid
type WasteType struct {
	Name string `json:"name"`
}

//Save - save waste data
func (waste *Waste) Save(env *config.Env) error {
	session, err := env.DB.GetConnection()
	if err != nil {
		return err
	}

	res, err := gorethink.Table("waste").Insert(waste, gorethink.InsertOpts{Conflict: "replace"}).RunWrite(session)
	if err != nil {
		// log.Println("error with upsert from waste upsert in waste.Save()")
		return err
	}

	if res.Inserted > 0 {
		waste.ID = res.GeneratedKeys[0]
	}
	return nil
}

//GetAll - get all waste by user and child id.
func (waste *Waste) GetAll(env *config.Env, user *User) ([]Waste, error) {
	session, err := env.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	family, err := user.GetFamily(env)
	if err != nil {
		return nil, err
	}

	res, err := gorethink.Table("waste").
		Filter(
			map[string]interface{}{
				"familyID": family.ID,
			}).
		OrderBy(gorethink.Desc("timestamp")).Run(session)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	var rows []Waste
	err = res.All(&rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

//GetByID - get a waste record by id
func (waste *Waste) GetByID(env *config.Env, id string) error {
	session, err := env.DB.GetConnection()
	if err != nil {
		return err
	}
	res, err := gorethink.Table("waste").Get(id).Run(session)
	if err != nil {
		return err
	}
	defer res.Close()
	err = res.One(&waste)
	if err != nil {
		return err
	}
	return nil
}

//WasteGetStats - get waste stats for one child for the last 24 hours.
func WasteGetStats(env *config.Env, user *User, child *Child) (WasteSummary, error) {
	session, err := env.DB.GetConnection()
	if err != nil {
		return WasteSummary{}, err
	}

	end := time.Now()
	start := end.AddDate(0, 0, -1)

	res, err := gorethink.Table("waste").
		Filter(map[string]interface{}{
			"userid":  user.ID,
			"childid": child.ID,
		}).
		Filter(gorethink.Row.Field("timestamp").During(start, end)).
		OrderBy(gorethink.Desc("timestamp")).
		Run(session)
	if err != nil {
		return WasteSummary{}, err
	}
	defer res.Close()

	var rows []Waste
	err = res.All(&rows)
	if err != nil {
		return WasteSummary{}, err
	}

	//build summary output
	summary := WasteSummary{
		Data:  rows,
		Total: make(map[int]int),
	}

	for _, x := range rows {
		if _, ok := summary.Total[x.Type]; !ok {
			summary.Total[x.Type] = 0
		}
		summary.Total[x.Type]++
	}
	return summary, nil
}
