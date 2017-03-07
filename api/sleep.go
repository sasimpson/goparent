package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func initSleepHandlers(r *mux.Router) {
	s := r.PathPrefix("/sleep").Subrouter()
	s.HandleFunc("/", sleepGetHandler).Methods("GET")
	s.HandleFunc("/", sleepNewHandler).Methods("POST")
	s.HandleFunc("/{id}", sleepViewHandler).Methods("GET")
	s.HandleFunc("/{id}", sleepEditHandler).Methods("PUT")
	s.HandleFunc("/{id}", sleepDeleteHandler).Methods("DELETE")
}

func sleepGetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GET no id")
}

func sleepViewHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprintf(w, "GET with id %s", id)
}

func sleepEditHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprintf(w, "PUT with id %s", id)
}

func sleepNewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "POST with data:")
}

func sleepDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprintf(w, "DELETE with id %s", id)
}
