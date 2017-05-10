package models

import (
	"log"
	"time"

	gorethink "gopkg.in/gorethink/gorethink.v3"
)

//Waste - structure for holding waste data such as diapers
type Waste struct {
	ID        string    `json:"id" gorethink:"id,omitempty"`
	Type      int       `json:"wasteType" gorethink:"wasteType"`
	Notes     string    `json:"notes" gorethink:"notes"`
	UserID    string    `json:"userid" gorethink:"userid"`
	TimeStamp time.Time `json:"timestamp" gorethink:"timestamp"`
}

//WasteType - the type of waste, solid, liquid, solid & liquid
type WasteType struct {
	Name string `json:"name"`
}

var (
	Solid       = WasteType{Name: "Solid"}
	Liquid      = WasteType{Name: "Liquid"}
	SolidLiquid = WasteType{Name: "Solid & Liquid"}
)

func (waste *Waste) Save() error {
	session, err := GetConnection()
	if err != nil {
		return err
	}
	defer session.Close()
	resp, err := gorethink.Table("waste").Insert(waste, gorethink.InsertOpts{Conflict: "replace"}).RunWrite(session)
	if err != nil {
		log.Println("error with upsert from sleep upsert in waste.Save()")
		return err
	}
	if resp.Inserted > 0 {
		log.Println(resp.GeneratedKeys)
		waste.ID = resp.GeneratedKeys[0]
	}
	return nil
}

func (waste *Waste) GetAll(user *User) ([]Waste, error) {
	session, err := GetConnection()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	resp, err := gorethink.Table("waste").Filter(map[string]interface{}{"userid": user.ID}).OrderBy(gorethink.Desc("timestamp")).Run(session)
	if err != nil {
		log.Println("error with get in waste.GetAll()")
		return nil, err
	}
	var rows []Waste
	err = resp.All(&rows)
	if err != nil {
		log.Println("error with getting")
		return nil, err
	}
	return rows, nil
}

func (waste *Waste) GetByID(id string) error {
	session, err := GetConnection()
	if err != nil {
		return err
	}
	defer session.Close()
	resp, err := gorethink.Table("waste").Get(id).Run(session)
	if err != nil {
		return err
	}
	err = resp.One(&waste)
	if err != nil {
		return err
	}
	return nil
}
