package models

import "time"

//Waste - structure for holding waste data such as diapers
type Waste struct {
	Type      WasteType
	Amount    WasteAmount
	Notes     string
	TimeStamp time.Time
}

//WasteType - the type of waste, solid, liquid, solid & liquid
type WasteType struct {
	Name string
}

//WasteAmount - the size, S, M, L
type WasteAmount struct {
	Name string
}

var (
	Solid       = WasteType{Name: "Solid"}
	Liquid      = WasteType{Name: "Liquid"}
	SolidLiquid = WasteType{Name: "Solid & Liquid"}
	WasteS      = WasteAmount{Name: "Small"}
	WasteM      = WasteAmount{Name: "Medium"}
	WasteL      = WasteAmount{Name: "Large"}
)
