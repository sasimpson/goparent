package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/rethinkdb"
)

var (
	env        *goparent.Env
	createFlag bool
	genFlag    bool
	childID    string
	userID     string
	date       string
	startDate  string
	endDate    string
)

func main() {
	env, dbenv := rethinkdb.InitRethinkDBConfig()

	flag.BoolVar(&createFlag, "createTables", false, "create all the needed tables")
	flag.BoolVar(&genFlag, "generate", false, "generate test data")
	flag.StringVar(&childID, "child", "", "child id for test data")
	flag.StringVar(&userID, "user", "", "user id for test data")
	flag.StringVar(&date, "date", time.Now().Format("2006-01-02"), "day for test data")
	flag.StringVar(&startDate, "startDate", "", "date to start filling data")
	flag.StringVar(&endDate, "endDate", "", "date to end filling data")
	flag.Parse()

	//create tables in the database
	if createFlag {
		log.Println("creating tables")
		rethinkdb.CreateTables(dbenv)
		os.Exit(0)
	}

	//if generate, just run it and exit
	if genFlag {
		log.Println("generating some data")

		if startDate != "" {
			if endDate == "" {
				endDate = time.Now().Format("2006-01-02")
			}
			sd, err := time.Parse("2006-01-02", startDate)
			if err != nil {
				panic(err)
			}

			ed, err := time.Parse("2006-01-02", endDate)
			if err != nil {
				panic(err)
			}

			for i := sd; i.Before(ed); i = i.AddDate(0, 0, 1) {
				log.Printf("\tgenerating for %s", i)
				err := generateRandomData(env, dbenv, childID, userID, i.Format("2006-01-02"))
				if err != nil {
					panic(err)
				}
			}
			os.Exit(0)
		}
		//if no start date is passed, then we assume to use the the date flag.
		err := generateRandomData(env, dbenv, childID, userID, date)
		if err != nil {
			panic(err)
		}

		os.Exit(0)
	}
}

func generateRandomData(env *goparent.Env, dbenv *rethinkdb.DBEnv, childID string, userID string, dateString string) error {
	var children []*goparent.Child
	var child *goparent.Child
	var user *goparent.User
	var family *goparent.Family

	if userID == "" {
		return errors.New("must have a user id")
	}
	userService := rethinkdb.UserService{Env: env, DB: dbenv}
	familyService := rethinkdb.FamilyService{Env: env, DB: dbenv}
	childService := rethinkdb.ChildService{Env: env, DB: dbenv}
	ctx := context.Background()

	user, err := userService.User(ctx, userID)
	if err != nil {
		return err
	}

	family, err = userService.GetFamily(ctx, user)
	if err != nil {
		return err
	}

	switch childID {
	case "":
		children, err = familyService.Children(ctx, family)
		if err != nil {
			return err
		}
	default:
		child, err = childService.Child(ctx, childID)
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
		log.Printf("\t\tfor child: %s", child.ID)
		generateRandomDiaper(ctx, env, dbenv, child, user, family, date)
		generateRandomSleep(ctx, env, dbenv, child, user, family, date)
		generateRandomFeeding(ctx, env, dbenv, child, user, family, date)
	}

	return nil
}

func generateRandomDiaper(ctx context.Context, env *goparent.Env, dbenv *rethinkdb.DBEnv, child *goparent.Child, user *goparent.User, family *goparent.Family, date time.Time) {
	wasteService := rethinkdb.WasteService{Env: env, DB: dbenv}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var numberOfEntries = r.Intn(7) + 7
	log.Printf("\t\t\tnumber of diaper entries: %d", numberOfEntries)
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	for x := 0; x < numberOfEntries; x++ {
		randoTime := time.Unix(date.Unix()+r.Int63n(86400), 0)
		diaper := &goparent.Waste{
			TimeStamp: randoTime,
			ChildID:   child.ID,
			UserID:    user.ID,
			FamilyID:  family.ID,
			Type:      r.Intn(3) + 1,
		}
		wasteService.Save(ctx, diaper)
	}
}

func generateRandomSleep(ctx context.Context, env *goparent.Env, dbenv *rethinkdb.DBEnv, child *goparent.Child, user *goparent.User, family *goparent.Family, date time.Time) {
	sleepService := rethinkdb.SleepService{Env: env, DB: dbenv}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var numberOfEntries = 8
	log.Printf("\t\t\tnumber of sleep entries: %d", numberOfEntries)
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	var sleeps []*goparent.Sleep
	for x := 0; x < numberOfEntries; x++ {
		if len(sleeps) > 0 {
			startDate = sleeps[len(sleeps)-1].End
		}
		randomInterval := (r.Int63n(60) + 60) * 60
		sleep := &goparent.Sleep{
			ChildID:  child.ID,
			UserID:   user.ID,
			FamilyID: family.ID,
			Start:    time.Unix(startDate.Unix()+randomInterval, 0),
			End:      time.Unix(startDate.Unix()+randomInterval+(5400+r.Int63n(1800)), 0),
		}
		sleepService.Save(sleep)
		sleeps = append(sleeps, sleep)
	}
}

func generateRandomFeeding(ctx context.Context, env *goparent.Env, dbenv *rethinkdb.DBEnv, child *goparent.Child, user *goparent.User, family *goparent.Family, date time.Time) {
	feedingService := rethinkdb.FeedingService{Env: env, DB: dbenv}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var numberOfEntries = 6 + r.Intn(4)
	log.Printf("\t\t\tnumber of feeding entries: %d", numberOfEntries)
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	var feedings []*goparent.Feeding

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
			feeding := &goparent.Feeding{
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
			feedingService.Save(feeding)
			feedingService.Save(feeding2)
			feedings = append(feedings, feeding)
			feedings = append(feedings, feeding2)
		} else {
			feeding := &goparent.Feeding{
				TimeStamp: time.Unix(startDate.Unix()+randomInterval, 0),
				Type:      feedingType,
				Amount:    float32(r.Intn(7) + 1),
				UserID:    user.ID,
				FamilyID:  family.ID,
				ChildID:   child.ID,
			}
			feedingService.Save(feeding)
			feedings = append(feedings, feeding)
		}
	}
}
