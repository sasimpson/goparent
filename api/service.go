package api

import (
	"fmt"
	"log"
	"net/http"

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

	log.Println("starting service on port 8000")
	http.Handle("/", r)
	http.ListenAndServe(":8000", nil)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "please check docs at: https://github.com/sasimpson/goparent")
}
