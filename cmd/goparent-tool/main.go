package main

import (
	"errors"
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/sasimpson/goparent/config"
	"github.com/sasimpson/goparent/models"
)

var (
	env        *config.Env
	createFlag bool
	genFlag    bool
	childID    string
	userID     string
	date       string
)

func main() {
	env = config.InitConfig()

	flag.BoolVar(&createFlag, "createTables", false, "create all the needed tables")
	flag.BoolVar(&genFlag, "generate", false, "generate test data")
	flag.StringVar(&childID, "child", "", "child id for test data")
	flag.StringVar(&userID, "user", "", "user id for test data")
	flag.StringVar(&date, "date", time.Now().Format("2006-01-02"), "day for test data")
	flag.Parse()

	//create tables in the database
	if createFlag {
		log.Println("creating tables")
		config.CreateTables(env)
		os.Exit(0)
	}

	//if generate, just run it and exit
	if genFlag {
		log.Println("generating some data")
		err := generateRandomData(env, childID, userID, date)
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}
}

func generateRandomData(env *config.Env, childID string, userID string, dateString string) error {
	var children []models.Child
	var child models.Child
	var user models.User
	var family models.Family

	if userID == "" {
		return errors.New("must have a user id")
	}

	err := user.GetUser(env, userID)
	if err != nil {
		return err
	}

	family, err = user.GetFamily(env)
	if err != nil {
		return err
	}

	switch childID {
	case "":
		children, err = family.GetAllChildren(env)
		if err != nil {
			return err
		}
	default:
		err = child.GetChild(env, &user, childID)
		if err != nil {
			return err
		}

		children = append(children, child)
	}

	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return err
	}
	for _, child := range children {
		generateRandomDiaper(env, child, user, family, date)
		generateRandomSleep(env, child, user, family, date)
		generateRandomFeeding(env, child, user, family, date)
	}

	return nil
}

func generateRandomDiaper(env *config.Env, child models.Child, user models.User, family models.Family, date time.Time) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var numberOfEntries = r.Intn(7) + 7
	log.Printf("number of diaper entries: %d", numberOfEntries)
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	for x := 0; x < numberOfEntries; x++ {
		randoTime := time.Unix(date.Unix()+r.Int63n(86400), 0)
		diaper := models.Waste{
			TimeStamp: randoTime,
			ChildID:   child.ID,
			UserID:    user.ID,
			FamilyID:  family.ID,
			Type:      r.Intn(2) + 1,
		}
		diaper.Save(env)
	}
}

func generateRandomSleep(env *config.Env, child models.Child, user models.User, family models.Family, date time.Time) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var numberOfEntries = 8
	log.Printf("number of sleep entries: %d", numberOfEntries)
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	var sleeps []models.Sleep
	for x := 0; x < numberOfEntries; x++ {
		if len(sleeps) > 0 {
			startDate = sleeps[len(sleeps)-1].End
		}
		randomInterval := (r.Int63n(60) + 60) * 60
		sleep := models.Sleep{
			ChildID:  child.ID,
			UserID:   user.ID,
			FamilyID: family.ID,
			Start:    time.Unix(startDate.Unix()+randomInterval, 0),
			End:      time.Unix(startDate.Unix()+randomInterval+(5400+r.Int63n(1800)), 0),
		}
		sleep.Save(env)
		sleeps = append(sleeps, sleep)
	}
}

func generateRandomFeeding(env *config.Env, child models.Child, user models.User, family models.Family, date time.Time) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var numberOfEntries = 6 + r.Intn(4)
	log.Printf("number of feeding entries: %d", numberOfEntries)
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	var feedings []models.Feeding

	var feedingType string
	switch r.Intn(2) {
	case 0:
		feedingType = "bottle"
		break
	case 1:
		feedingType = "breast"
		break
	}

	for x := 0; x < numberOfEntries; x++ {
		if len(feedings) > 0 {
			startDate = feedings[len(feedings)-1].TimeStamp
		}
		randomInterval := (r.Int63n(60) + 60) * 60
		if feedingType == "breast" {
			feeding := models.Feeding{
				TimeStamp: time.Unix(startDate.Unix()+randomInterval, 0),
				Type:      feedingType,
				Side:      "right",
				Amount:    float32(r.Intn(29) + 1),
				UserID:    user.ID,
				FamilyID:  family.ID,
				ChildID:   child.ID,
			}
			feeding2 := feeding
			feeding2.Side = "left"
			feeding.Save(env)
			feeding2.Save(env)
			feedings = append(feedings, feeding)
			feedings = append(feedings, feeding2)
		} else {
			feeding := models.Feeding{
				TimeStamp: time.Unix(startDate.Unix()+randomInterval, 0),
				Type:      feedingType,
				Amount:    float32(r.Intn(7) + 1),
				UserID:    user.ID,
				FamilyID:  family.ID,
				ChildID:   child.ID,
			}
			feeding.Save(env)
			feedings = append(feedings, feeding)
		}
	}
}
