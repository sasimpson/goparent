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
	// w := r.PathPrefix("/waste").Subrouter()
	// s := r.PathPrefix("/sleep").Subrouter()

	f.HandleFunc("/", FeedingGet).Methods("GET")
	f.HandleFunc("/", FeedingNew).Methods("POST")
	f.HandleFunc("/{id}", FeedingView).Methods("GET")
	f.HandleFunc("/{id}", FeedingEdit).Methods("PUT")
	f.HandleFunc("/{id}", FeedingDelete).Methods("DELETE")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

//FeedingGet -
func FeedingGet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GET no id")
}

//FeedingView -
func FeedingView(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprintf(w, "GET with id %s", id)
}

// FeedingEdit -
func FeedingEdit(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprintf(w, "PUT with id %s", id)
}

//FeedingNew -
func FeedingNew(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "POST with data:")
}

//FeedingDelete -
func FeedingDelete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprintf(w, "DELETE with id %s", id)
}

//WasteGet
//WasteView
//WasteEdit
//WasteNew
//WasteDelete

//SleepGet
//SleepView
//SleepEdit
//SleepNew
//SleepDelete
