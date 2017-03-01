package models

import "time"

//Feeding - main data structure for storing feeding data
type Feeding struct {
	Type      FoodType
	Amount    FoodAmount
	TimeStamp time.Time
}

//FoodType - structure for the types of food ie formula, breast milk, solids, etc.
type FoodType struct {
	Name string
}

//FoodAmount - structure for the amount of food consumed, number and unit.
type FoodAmount struct {
	Number float32
	Unit   string
}

