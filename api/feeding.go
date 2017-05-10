package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sasimpson/goparent/models"
)

type FeedingRequest struct {
	FeedingData models.Feeding `json:"feedingData"`
}

type FeedingResponse struct {
	FeedingData []models.Feeding `json:"feedingData"`
}

func initFeedingHandlers(r *mux.Router) {
	f := r.PathPrefix("/feeding").Subrouter()
	f.HandleFunc("", FeedingGetHandler).Methods("GET")
	f.HandleFunc("", FeedingNewHandler).Methods("POST")
	f.HandleFunc("/{id}", FeedingViewHandler).Methods("GET")
	f.HandleFunc("/{id}", FeedingEditHandler).Methods("PUT")
	f.HandleFunc("/{id}", FeedingDeleteHandler).Methods("DELETE")
}

//-------------------

//FeedingGetHandler -
func FeedingGetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GET feeding")
	user, err := validateAuthToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	var feeding models.Feeding
	feedingData, err := feeding.GetAll(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	feedingResponse := FeedingResponse{FeedingData: feedingData}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feedingResponse.FeedingData)
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
	log.Println("POST Feeding")
	user, err := validateAuthToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var feedingRequest FeedingRequest
	err = decoder.Decode(&feedingRequest)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	feedingRequest.FeedingData.UserID = user.ID
	err = feedingRequest.FeedingData.Save()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusConflict)
	}
	json.NewEncoder(w).Encode(feedingRequest.FeedingData)
}

//FeedingDeleteHandler -
func FeedingDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprintf(w, "DELETE with id %s", id)
}
