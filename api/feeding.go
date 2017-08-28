package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sasimpson/goparent/config"
	"github.com/sasimpson/goparent/models"
)

type FeedingRequest struct {
	FeedingData models.Feeding `json:"feedingData"`
}

type FeedingResponse struct {
	FeedingData []models.Feeding `json:"feedingData"`
}

func initFeedingHandlers(env *config.Env, r *mux.Router) {
	f := r.PathPrefix("/feeding").Subrouter()
	f.Handle("", FeedingGetHandler(env)).Methods("GET")
	f.Handle("", FeedingNewHandler(env)).Methods("POST")
	f.Handle("/{id}", FeedingViewHandler(env)).Methods("GET")
	f.Handle("/{id}", FeedingEditHandler(env)).Methods("PUT")
	f.Handle("/{id}", FeedingDeleteHandler(env)).Methods("DELETE")
}

//-------------------

//FeedingGetHandler -
func FeedingGetHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET feeding")
		user, err := validateAuthToken(env, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		var feeding models.Feeding
		feedingData, err := feeding.GetAll(env, &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		feedingResponse := FeedingResponse{FeedingData: feedingData}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(feedingResponse.FeedingData)
	})
}

//FeedingViewHandler -
func FeedingViewHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		fmt.Fprintf(w, "GET with id %s", id)
	})
}

// FeedingEditHandler -
func FeedingEditHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		fmt.Fprintf(w, "PUT with id %s", id)
	})
}

//FeedingNewHandler -
func FeedingNewHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("POST Feeding")
		user, err := validateAuthToken(env, r)
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
		err = feedingRequest.FeedingData.Save(env)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusConflict)
		}
		json.NewEncoder(w).Encode(feedingRequest.FeedingData)
	})
}

//FeedingDeleteHandler -
func FeedingDeleteHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		fmt.Fprintf(w, "DELETE with id %s", id)
	})
}
