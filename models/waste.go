package models

import (
	"log"
	"time"

	"github.com/sasimpson/goparent/config"
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

func (waste *Waste) Save(env *config.Env) error {
	session, err := env.DB.GetConnection()
	if err != nil {
		return err
	}
	res, err := gorethink.Table("waste").Insert(waste, gorethink.InsertOpts{Conflict: "replace"}).RunWrite(session)
	if err != nil {
		log.Println("error with upsert from sleep upsert in waste.Save()")
		return err
	}
	if res.Inserted > 0 {
		log.Println(res.GeneratedKeys)
		waste.ID = res.GeneratedKeys[0]
	}
	return nil
}

func (waste *Waste) GetAll(env *config.Env, user *User) ([]Waste, error) {
	session, err := env.DB.GetConnection()
	if err != nil {
		return nil, err
	}
	res, err := gorethink.Table("waste").Filter(map[string]interface{}{"userid": user.ID}).OrderBy(gorethink.Desc("timestamp")).Run(session)
	if err != nil {
		log.Println("error with get in waste.GetAll()")
		return nil, err
	}
	defer res.Close()
	var rows []Waste
	err = res.All(&rows)
	if err != nil {
		log.Println("error with getting")
		return nil, err
	}
	return rows, nil
}

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
