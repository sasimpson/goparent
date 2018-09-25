package goparent

import (
	"context"
	"errors"
	"net/http"
	"time"

	"gopkg.in/gorethink/gorethink.v3"

	jwt "github.com/dgrijalva/jwt-go"
)

//Env - container for all environment configuraitons
type Env struct {
	Service Service
	DB      Datastore
	Auth    Authentication
}

//Service - structure for service configurations
type Service struct {
	Host string
	Port int
}

//Authentication - structure for authentication configurations
type Authentication struct {
	SigningKey []byte
}

//DBEnv - Environment for DB settings
type DBEnv struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
	Session  gorethink.QueryExecutor
}

//Datastore -
type Datastore interface {
	GetConnection() error
	GetContext(*http.Request) context.Context
}

//ErrExistingInvitation -
var ErrExistingInvitation = errors.New("invitation already exists for that parent")

//User -
type User struct {
	ID            string `json:"id" gorethink:"id,omitempty"`
	Name          string `json:"name" gorethink:"name"`
	Email         string `json:"email" gorethink:"email"`
	Username      string `json:"username" gorethink:"username"`
	Password      string `json:"-" gorethink:"password"`
	CurrentFamily string `json:"currentFamily" gorethink:"currentFamily"`
}

//UserClaims - structure for inserting claims into a jwt auth token
type UserClaims struct {
	ID       string
	Name     string
	Email    string
	Username string
	Password string
	jwt.StandardClaims
}

//UserService -
type UserService interface {
	User(context.Context, string) (*User, error)
	UserByLogin(context.Context, string, string) (*User, error)
	Save(context.Context, *User) error
	GetToken(*User) (string, error)
	ValidateToken(context.Context, string) (*User, bool, error)
	GetFamily(context.Context, *User) (*Family, error)
	GetAllFamily(context.Context, *User) ([]*Family, error)
}

//UserInvitation - structure for storing invitations
type UserInvitation struct {
	ID          string    `json:"id" gorethink:"id,omitempty"`
	UserID      string    `json:"userID" gorethink:"userID"`
	InviteEmail string    `json:"inviteEmail" gorethink:"inviteEmail"`
	Timestamp   time.Time `json:"timestamp" gorethink:"timestamp"`
}

//UserInvitationService -
type UserInvitationService interface {
	InviteParent(*User, string, time.Time) error
	SentInvites(*User) ([]*UserInvitation, error)
	Invite(string) (*UserInvitation, error)
	Invites(*User) ([]*UserInvitation, error)
	Accept(*User, string) error
	Delete(*UserInvitation) error
}

//Family -
type Family struct {
	ID          string    `json:"id" gorethink:"id,omitempty"`
	Admin       string    `json:"admin" gorethink:"admin"`
	Members     []string  `json:"members" gorethink:"members"`
	CreatedAt   time.Time `json:"created_at" gorethink:"created_at"`
	LastUpdated time.Time `json:"last_updated" gorethink:"last_updated"`
}

//FamilyService -
type FamilyService interface {
	Save(context.Context, *Family) error
	Family(context.Context, string) (*Family, error)
	Children(*Family) ([]*Child, error)
	AddMember(context.Context, *Family, *User) error
	GetAdminFamily(context.Context, *User) (*Family, error)
	// Delete(*Family) (int, error)
}

//Child -
type Child struct {
	ID          string    `json:"id" gorethink:"id,omitempty"`
	Name        string    `json:"name" gorethink:"name"`
	ParentID    string    `json:"parentID" gorethink:"parentID"`
	FamilyID    string    `json:"familyID" gorethink:"familyID"`
	Birthday    time.Time `json:"birthday" gorethink:"birthday"`
	CreatedAt   time.Time `json:"created_at" gorethink:"created_at"`
	LastUpdated time.Time `json:"last_updated" gorethink:"last_updated"`
}

//ChildService -
type ChildService interface {
	Save(context.Context, *Child) error
	Child(context.Context, string) (*Child, error)
	Delete(context.Context, *Child) (int, error)
}

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

//FeedingChartData -
type FeedingChartData struct {
	Start   time.Time             `json:"start"`
	End     time.Time             `json:"end"`
	Dataset []FeedingChartDataset `json:"dataset"`
}

//FeedingChartDataset -
type FeedingChartDataset struct {
	Date  time.Time `json:"date"`
	Type  string    `json:"type"`
	Count int       `json:"count"`
	Sum   float32   `json:"sum"`
}

//FeedingService -
type FeedingService interface {
	Save(*Feeding) error
	Feeding(*Family, uint64) ([]*Feeding, error)
	Stats(*Child) (*FeedingSummary, error)
	GraphData(*Child) (*FeedingChartData, error)
}

//Sleep - tracks the baby's sleep start and end.
type Sleep struct {
	ID       string    `json:"id" gorethink:"id,omitempty"`
	Start    time.Time `json:"start" gorethink:"start"`
	End      time.Time `json:"end" gorethink:"end"`
	UserID   string    `json:"userid" gorethink:"userID"`
	FamilyID string    `json:"familyid" gorethink:"familyID"`
	ChildID  string    `json:"childID" gorethink:"childID"`
}

//SleepSummary - structure for the sleep summary data
type SleepSummary struct {
	Data  []Sleep `json:"data"`
	Total int64   `json:"total"`
	Mean  float64 `json:"mean"`
	Range int     `json:"range"`
}

//SleepChartData -
type SleepChartData struct {
	Start   time.Time             `json:"start"`
	End     time.Time             `json:"end"`
	Dataset []FeedingChartDataset `json:"dataset"`
}

//SleepChartDataset -
type SleepChartDataset struct {
	Date  time.Time `json:"date"`
	Type  string    `json:"type"`
	Count int       `json:"count"`
	Sum   float32   `json:"sum"`
}

//SleepService -
type SleepService interface {
	Save(*Sleep) error
	Sleep(*Family, uint64) ([]*Sleep, error)
	Stats(*Child) (*SleepSummary, error)
	Status(*Family, *Child) (bool, error)
	Start(*Sleep, *Family, *Child) error
	End(*Sleep, *Family, *Child) error
	GraphData(*Child) (*SleepChartData, error)
}

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

//WasteGraphData -
type WasteGraphData struct {
	Date  time.Time `json:"date"`
	Type  int       `json:"type"`
	Count int       `json:"count"`
}

//WasteType - the type of waste, solid, liquid, solid & liquid
type WasteType struct {
	Name string `json:"name"`
}

//WasteService -
type WasteService interface {
	Save(*Waste) error
	Waste(*Family, uint64) ([]*Waste, error)
	Stats(*Child) (*WasteSummary, error)
	GraphData(*Child) (*WasteChartData, error)
}

//wasteChartData -
type WasteChartData struct {
	Start   time.Time           `json:"start"`
	End     time.Time           `json:"end"`
	Dataset []WasteChartDataset `json:"dataset"`
}

//WasteChartDataset -
type WasteChartDataset struct {
	Date  time.Time `json:"date"`
	Type  int       `json:"type"`
	Count int       `json:"count"`
}
