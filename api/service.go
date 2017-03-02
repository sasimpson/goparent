package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

//RunService - Runs service interfaces for app
func RunService() {
	r := mux.NewRouter()
	f := r.PathPrefix("/food").Subrouter()
	w := r.PathPrefix("/waste").Subrouter()
	s := r.PathPrefix("/sleep").Subrouter()
	//Feeding
	f.HandleFunc("/", FeedingGetHandler).Methods("GET")
	f.HandleFunc("/", FeedingNewHandler).Methods("POST")
	f.HandleFunc("/{id}", FeedingViewHandler).Methods("GET")
	f.HandleFunc("/{id}", FeedingEditHandler).Methods("PUT")
	f.HandleFunc("/{id}", FeedingDeleteHandler).Methods("DELETE")
	//Waste
	w.HandleFunc("/", WasteGetHandler).Methods("GET")
	w.HandleFunc("/", WasteNewHandler).Methods("POST")
	w.HandleFunc("/{id}", WasteViewHandler).Methods("GET")
	w.HandleFunc("/{id}", WasteEditHandler).Methods("PUT")
	w.HandleFunc("/{id}", WasteDeleteHandler).Methods("DELETE")
	//Sleep
	s.HandleFunc("/", SleepGetHandler).Methods("GET")
	s.HandleFunc("/", SleepNewHandler).Methods("POST")
	s.HandleFunc("/{id}", SleepViewHandler).Methods("GET")
	s.HandleFunc("/{id}", SleepEditHandler).Methods("PUT")
	s.HandleFunc("/{id}", SleepDeleteHandler).Methods("DELETE")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

//-------------------

//FeedingGetHandler -
func FeedingGetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GET no id")
}

//FeedingViewHandler -
func FeedingViewHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprintf(w, "GET with id %s", id)
}

// FeedingEditHandler -
func FeedingEditHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprintf(w, "PUT with id %s", id)
}

//FeedingNewHandler -
func FeedingNewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "POST with data:")
}

//FeedingDeleteHandler -
func FeedingDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprintf(w, "DELETE with id %s", id)
}

//-------------------

//WasteGetHandler -
func WasteGetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GET no id")
}

//WasteViewHandler -
func WasteViewHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprintf(w, "GET with id %s", id)
}

// WasteEditHandler -
func WasteEditHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprintf(w, "PUT with id %s", id)
}

//WasteNewHandler -
func WasteNewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "POST with data:")
}

//WasteDeleteHandler -
func WasteDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprintf(w, "DELETE with id %s", id)
}

//-------------------

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
