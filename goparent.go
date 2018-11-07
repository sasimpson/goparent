package goparent

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

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

func (user User) String() string {
	return fmt.Sprintf("[%s] %s <%s>", user.ID, user.Name, user.Email)
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
	InviteParent(context.Context, *User, string, time.Time) error
	SentInvites(context.Context, *User) ([]*UserInvitation, error)
	Invite(context.Context, string) (*UserInvitation, error)
	Invites(context.Context, *User) ([]*UserInvitation, error)
	Accept(context.Context, *User, string) error
	Delete(context.Context, *UserInvitation) error
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
	Children(context.Context, *Family) ([]*Child, error)
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

func (child Child) String() string {
	return fmt.Sprintf("[%s] %s", child.ID, child.Name)
}

//ChildService -
type ChildService interface {
	Save(context.Context, *Child) error
	Child(context.Context, string) (*Child, error)
	Delete(context.Context, *Child) (int, error)
}

//Feeding - main data structure for storing feeding data
type Feeding struct {
	ID          string    `json:"id" gorethink:"id,omitempty"`
	Type        string    `json:"feedingType" gorethink:"feedingType"`
	Amount      float32   `json:"feedingAmount" gorethink:"feedingAmount"`
	Side        string    `json:"feedingSide" gorethink:"feedingSide,omitempty"`
	UserID      string    `json:"userid" gorethink:"userID"`
	FamilyID    string    `json:"familyid" gorethink:"familyID"`
	TimeStamp   time.Time `json:"timestamp" gorethink:"timestamp"`
	ChildID     string    `json:"childID" gorethink:"childID"`
	CreatedAt   time.Time `json:"createdAt" gorethink:"createdAt"`
	LastUpdated time.Time `json:"lastUpdated" gorethink:"lastUpdated"`
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
	Save(context.Context, *Feeding) error
	Feeding(context.Context, *Family, uint64) ([]*Feeding, error)
	Stats(context.Context, *Child) (*FeedingSummary, error)
	GraphData(context.Context, *Child) (*FeedingChartData, error)
}

//Sleep - tracks the baby's sleep start and end.
type Sleep struct {
	ID          string    `json:"id" gorethink:"id,omitempty"`
	Start       time.Time `json:"start" gorethink:"start"`
	End         time.Time `json:"end" gorethink:"end"`
	UserID      string    `json:"userid" gorethink:"userID"`
	FamilyID    string    `json:"familyid" gorethink:"familyID"`
	ChildID     string    `json:"childID" gorethink:"childID"`
	CreatedAt   time.Time `json:"createdAt" gorethink:"createdAt"`
	LastUpdated time.Time `json:"lastUpdated" gorethink:"lastUpdated"`
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
	Start   time.Time           `json:"start"`
	End     time.Time           `json:"end"`
	Dataset []SleepChartDataset `json:"dataset"`
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
	Save(context.Context, *Sleep) error
	Sleep(context.Context, *Family, uint64) ([]*Sleep, error)
	Stats(context.Context, *Child) (*SleepSummary, error)
	Status(context.Context, *Family, *Child) (bool, error)
	Start(context.Context, *Sleep, *Family, *Child) error
	End(context.Context, *Sleep, *Family, *Child) error
	GraphData(context.Context, *Child) (*SleepChartData, error)
}

//Waste - structure for holding waste data such as diapers
type Waste struct {
	ID          string    `json:"id" gorethink:"id,omitempty"`
	Type        int       `json:"wasteType" gorethink:"wasteType"`
	Notes       string    `json:"notes" gorethink:"notes"`
	UserID      string    `json:"userid" gorethink:"userID"`
	FamilyID    string    `json:"familyid" gorethink:"familyID"`
	ChildID     string    `json:"childid" gorethink:"childID"`
	TimeStamp   time.Time `json:"timestamp" gorethink:"timestamp"`
	CreatedAt   time.Time `json:"createdAt" gorethink:"createdAt"`
	LastUpdated time.Time `json:"lastUpdated" gorethink:"lastUpdated"`
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
	Save(context.Context, *Waste) error
	Waste(context.Context, *Family, uint64) ([]*Waste, error)
	Stats(context.Context, *Child) (*WasteSummary, error)
	GraphData(context.Context, *Child) (*WasteChartData, error)
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
