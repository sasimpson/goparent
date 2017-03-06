package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func initSleepHandlers(r *mux.Router) {
	s := r.PathPrefix("/sleep").Subrouter()
	s.HandleFunc("/", SleepGetHandler).Methods("GET")
	s.HandleFunc("/", SleepNewHandler).Methods("POST")
	s.HandleFunc("/{id}", SleepViewHandler).Methods("GET")
	s.HandleFunc("/{id}", SleepEditHandler).Methods("PUT")
	s.HandleFunc("/{id}", SleepDeleteHandler).Methods("DELETE")
}

//SleepGetHandler -
func SleepGetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GET no id")
}

//SleepViewHandler -
func SleepViewHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprintf(w, "GET with id %s", id)
}

// SleepEditHandler -
func SleepEditHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprintf(w, "PUT with id %s", id)
}

//SleepNewHandler -
func SleepNewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "POST with data:")
}

//SleepDeleteHandler -
func SleepDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprintf(w, "DELETE with id %s", id)
}
