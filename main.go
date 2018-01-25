package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/sasimpson/goparent/api"
	"github.com/sasimpson/goparent/config"
	"github.com/sasimpson/goparent/models"
)

func main() {
	genFlag := flag.Bool("generate", false, "generate test data")
	childID := flag.String("child", "", "child id for test data")
	userID := flag.String("user", "", "user id for test data")
	date := flag.String("date", time.Now().Format("2006-01-02"), "day for test data")
	flag.Parse()

	env := config.InitConfig()

	//if generate, just run it and exit
	if *genFlag {
		log.Println("generating some data")
		generateRandomData(env, childID, userID, date)
		os.Exit(0)
	}

	config.CreateTables(env)
	api.RunService(env)
}

func generateRandomData(env *config.Env, childID *string, userID *string, dateString *string) {
	//TODO: handle error
	date, _ := time.Parse("2006-01-02", *dateString)
	generateRandomDiaper(env, childID, userID, date)
	generateRandomSleep(env, childID, userID, date)
	generateRandomFeeding(env, childID, userID, date)
}

func generateRandomDiaper(env *config.Env, childID *string, userID *string, date time.Time) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var numberOfEntries = r.Intn(7) + 7
	log.Printf("number of diaper entries: %d", numberOfEntries)

	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	for x := 0; x < numberOfEntries; x++ {
		randoTime := time.Unix(date.Unix()+r.Int63n(86400), 0)
		diaper := models.Waste{
			TimeStamp: randoTime,
			ChildID:   *childID,
			UserID:    *userID,
			Type:      r.Intn(2) + 1,
		}
		diaper.Save(env)
	}
}

func generateRandomSleep(env *config.Env, childID *string, userID *string, date time.Time) {
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
			ChildID: *childID,
			UserID:  *userID,
			Start:   time.Unix(startDate.Unix()+randomInterval, 0),
			End:     time.Unix(startDate.Unix()+randomInterval+(5400+r.Int63n(1800)), 0),
		}
		sleep.Save(env)
		sleeps = append(sleeps, sleep)
	}
}

func generateRandomFeeding(env *config.Env, childID *string, userID *string, date time.Time) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var numberOfEntries = 6 + r.Intn(4)
	log.Printf("number of feeding entries: %d", numberOfEntries)
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	var feedings []models.Feeding

	var feedingType string
	switch r.Intn(1) {
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
				UserID:    *userID,
				ChildID:   *childID,
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
				UserID:    *userID,
				ChildID:   *childID,
			}
			feeding.Save(env)
			feedings = append(feedings, feeding)
		}
	}
}
