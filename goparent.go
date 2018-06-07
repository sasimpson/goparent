package goparent

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var ErrExistingInvitation = errors.New("invitation already exists for that parent")

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

type UserService interface {
	User(string) (*User, error)
	UserByLogin(string, string) (*User, error)
	Save(*User) error
	GetToken(*User) (string, error)
	ValidateToken(string) (*User, bool, error)
	GetFamily(*User) (*Family, error)
	GetAllFamily(*User) ([]*Family, error)
}

//UserInvitation - structure for storing invitations
type UserInvitation struct {
	ID          string    `json:"id" gorethink:"id,omitempty"`
	UserID      string    `json:"userID" gorethink:"userID"`
	InviteEmail string    `json:"inviteEmail" gorethink:"inviteEmail"`
	Timestamp   time.Time `json:"timestamp" gorethink:"timestamp"`
}

type UserInvitationService interface {
	InviteParent(*User, string, time.Time) error
	SentInvites(*User) ([]*UserInvitation, error)
	Invite(string) (*UserInvitation, error)
	Invites(*User) ([]*UserInvitation, error)
	Accept(*User, string) error
	Delete(*UserInvitation) error
}

type Family struct {
	ID          string    `json:"id" gorethink:"id,omitempty"`
	Admin       string    `json:"admin" gorethink:"admin"`
	Members     []string  `json:"members" gorethink:"members"`
	CreatedAt   time.Time `json:"created_at" gorethink:"created_at"`
	LastUpdated time.Time `json:"last_updated" gorethink:"last_updated"`
}

type FamilyService interface {
	Save(*Family) error
	Family(string) (*Family, error)
	Children(*Family) ([]*Child, error)
	AddMember(*Family, *User) error
	GetAdminFamily(*User) (*Family, error)
	// Delete(*Family) (int, error)
}

type Child struct {
	ID       string    `json:"id" gorethink:"id,omitempty"`
	Name     string    `json:"name" gorethink:"name"`
	ParentID string    `json:"parentID" gorethink:"parentID"`
	FamilyID string    `json:"familyID" gorethink:"familyID"`
	Birthday time.Time `json:"birthday" gorethink:"birthday"`
}

type ChildService interface {
	Save(*Child) error
	Child(string) (*Child, error)
	Delete(*Child) (int, error)
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

type FeedingService interface {
	Save(*Feeding) error
	Feeding(*Family) ([]*Feeding, error)
	Stats(*Child) (*FeedingSummary, error)
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

type SleepService interface {
	Save(*Sleep) error
	Sleep(*Family) ([]*Sleep, error)
	Stats(*Child) (*SleepSummary, error)
	Status(*Family, *Child) (bool, error)
	Start(*Sleep, *Family, *Child) error
	End(*Sleep, *Family, *Child) error
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

//WasteType - the type of waste, solid, liquid, solid & liquid
type WasteType struct {
	Name string `json:"name"`
}

type WasteService interface {
	Save(*Waste) error
	Waste(*Family, uint64) ([]*Waste, error)
	Stats(*Child) (*WasteSummary, error)
}
