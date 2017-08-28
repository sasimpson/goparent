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
	UserID    string    `json:"userid" gorethink:"userid"`
	TimeStamp time.Time `json:"timestamp" gorethink:"timestamp"`
}

//Save - save the structure to the datastore
func (feeding *Feeding) Save(env *config.Env) error {
	session, err := env.DB.GetConnection()
	if err != nil {
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
	res, err := gorethink.Table("feeding").Filter(map[string]interface{}{"userid": user.ID}).OrderBy(gorethink.Desc("timestamp")).Run(session)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var rows []Feeding
	err = res.All(&rows)
	if err != nil {
		log.Println("error getting all")
		return nil, err
	}
	return rows, nil
}
