package models

import (
	"log"
	"time"

	"github.com/sasimpson/goparent/config"
	gorethink "gopkg.in/gorethink/gorethink.v3"
)

//Feeding - main data structure for storing feeding data
type Feeding struct {
	ID        string    `json:"id" gorethink:"id,omitempty"`
	Type      string    `json:"feedingType" gorethink:"feedingType"`
	Amount    float32   `json:"feedingAmount" gorethink:"feedingAmount"`
	Side      string    `json:"feedingSide" gorethink:"feedingSide,omitempty"`
	UserID    string    `json:"userid" gorethink:"userID"`
	FamilyID  string    `json:"familyid" gorethink:"familyID"`
	TimeStamp time.Time `json:"timestamp" gorethink:"timestamp"`
	ChildID   string    `json:"childID" gorethink:"childID"`
}

//FeedingSummary - represents feeding summary data
type FeedingSummary struct {
	Data  []Feeding          `json:"data"`
	Total map[string]float32 `json:"total"`
	Mean  map[string]float32 `json:"mean"`
	Range map[string]int     `json:"range"`
}

//Save - save the structure to the datastore
func (feeding *Feeding) Save(env *config.Env) error {
	session, err := env.DB.GetConnection()
	if err != nil {
		log.Println("error getting db connection")
		return err
	}

	res, err := gorethink.Table("feeding").Insert(feeding, gorethink.InsertOpts{Conflict: "replace"}).RunWrite(session)
	if err != nil {
		log.Println("error with upsert from feeding upsert in feeding.Save()")
		return err
	}

	if res.Inserted > 0 {
		feeding.ID = res.GeneratedKeys[0]
	}
	return nil
}

//GetAll - get all records for a user from the datastore
func (feeding *Feeding) GetAll(env *config.Env, user *User) ([]Feeding, error) {
	session, err := env.DB.GetConnection()
	if err != nil {
		return nil, err
	}

	family, err := user.GetFamily(env)
	if err != nil {
		return nil, err
	}

	res, err := gorethink.Table("feeding").
		Filter(map[string]interface{}{
			"familyID": family.ID,
		}).
		OrderBy(gorethink.Desc("timestamp")).
		Run(session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var rows []Feeding
	err = res.All(&rows)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

//FeedingGetStats - get feeding stats for one child for the last 24 hours.
func FeedingGetStats(env *config.Env, user *User, child *Child) (FeedingSummary, error) {
	session, err := env.DB.GetConnection()
	if err != nil {
		return FeedingSummary{}, err
	}

	end := time.Now()
	start := end.AddDate(0, 0, -1)

	res, err := gorethink.Table("feeding").
		Filter(map[string]interface{}{
			"childid": child.ID,
		}).
		Filter(gorethink.Row.Field("timestamp").During(start, end)).
		OrderBy(gorethink.Desc("timestamp")).
		Run(session)
	if err != nil {
		return FeedingSummary{}, err
	}
	defer res.Close()

	var rows []Feeding
	err = res.All(&rows)
	if err != nil {
		return FeedingSummary{}, err
	}
	//build summary output
	summary := FeedingSummary{
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
