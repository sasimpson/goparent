package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func initFeedingHandlers(r *mux.Router) {
	f := r.PathPrefix("/food").Subrouter()
	f.HandleFunc("/", FeedingGetHandler).Methods("GET")
	f.HandleFunc("/", FeedingNewHandler).Methods("POST")
	f.HandleFunc("/{id}", FeedingViewHandler).Methods("GET")
	f.HandleFunc("/{id}", FeedingEditHandler).Methods("PUT")
	f.HandleFunc("/{id}", FeedingDeleteHandler).Methods("DELETE")
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
