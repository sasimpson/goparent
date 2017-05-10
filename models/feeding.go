package models

import (
	"log"
	"time"

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

func (feeding *Feeding) Save() error {
	session, err := GetConnection()
	if err != nil {
		return err
	}
	defer session.Close()
	resp, err := gorethink.Table("feeding").Insert(feeding, gorethink.InsertOpts{Conflict: "replace"}).RunWrite(session)
	if err != nil {
		log.Println("error with upsert from feeding upsert in feeding.Save()")
		return err
	}
	if resp.Inserted > 0 {
		feeding.ID = resp.GeneratedKeys[0]
	}
	return nil
}

func (feeding *Feeding) GetAll(user *User) ([]Feeding, error) {
	session, err := GetConnection()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	resp, err := gorethink.Table("feeding").Filter(map[string]interface{}{"userid": user.ID}).OrderBy(gorethink.Desc("timestamp")).Run(session)
	if err != nil {
		return nil, err
	}
	var rows []Feeding
	err = resp.All(&rows)
	if err != nil {
		log.Println("error getting all")
		return nil, err
	}
	return rows, nil
}
