package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//RunService - Runs service interfaces for app
func RunService() {
	r := mux.NewRouter()
	a := r.PathPrefix("/api").Subrouter()
	a.HandleFunc("/", apiHandler)

	initUsersHandlers(a)
	initFeedingHandlers(a)
	initSleepHandlers(a)
	initWasteHandlers(a)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Accept", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Println("starting service on port 8000")
	http.Handle("/", r)
	http.ListenAndServe(":8000", handlers.CORS(originsOk, headersOk, methodsOk)(r))
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "please check docs at: https://github.com/sasimpson/goparent")
}
