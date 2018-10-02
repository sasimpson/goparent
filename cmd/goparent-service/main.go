package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/api"
	"github.com/sasimpson/goparent/rethinkdb"
)

func main() {
	env, dbenv := rethinkdb.InitRethinkDBConfig()
	runService(env, dbenv)
}

//RunService - Runs service interfaces for app
func runService(env *goparent.Env, dbenv *rethinkdb.DBEnv) {
	log.SetOutput(os.Stdout)
	serviceHandler := api.Handler{
		UserService:           &rethinkdb.UserService{Env: env, DB: dbenv},
		UserInvitationService: &rethinkdb.UserInviteService{Env: env, DB: dbenv},
		FamilyService:         &rethinkdb.FamilyService{Env: env, DB: dbenv},
		ChildService:          &rethinkdb.ChildService{Env: env, DB: dbenv},
		FeedingService:        &rethinkdb.FeedingService{Env: env, DB: dbenv},
		SleepService:          &rethinkdb.SleepService{Env: env, DB: dbenv},
		WasteService:          &rethinkdb.WasteService{Env: env, DB: dbenv},
		Env:                   env,
	}

	r := api.BuildAPIRouting(&serviceHandler)
	// setup cors, this should end up in Env.Service config, which this should receive.
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Accept", "Content-Type", "Authorization", "Origin"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	r.Use(simpleRequestLog)

	log.Printf("starting service on 8000")
	http.ListenAndServe(":8000", handlers.CORS(originsOk, headersOk, methodsOk)(r))
}

func simpleRequestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
