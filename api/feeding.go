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

//FeedingRequest - request structure for feedings
type FeedingRequest struct {
	FeedingData models.Feeding `json:"feedingData"`
}

//FeedingResponse - response structure for feedings
type FeedingResponse struct {
	FeedingData []models.Feeding `json:"feedingData"`
}

func initFeedingHandlers(env *config.Env, r *mux.Router) {
	f := r.PathPrefix("/feeding").Subrouter()
	f.Handle("", AuthRequired(FeedingGetHandler(env), env)).Methods("GET").Name("FeedingGet")
	f.Handle("", AuthRequired(FeedingNewHandler(env), env)).Methods("POST").Name("FeedingNew")
	f.Handle("/{id}", AuthRequired(FeedingViewHandler(env), env)).Methods("GET").Name("FeedingView")
	f.Handle("/{id}", AuthRequired(FeedingEditHandler(env), env)).Methods("PUT").Name("FeedingEdit")
	f.Handle("/{id}", AuthRequired(FeedingDeleteHandler(env), env)).Methods("DELETE").Name("FeedingDelete")
}

//FeedingGetHandler -
func FeedingGetHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET feeding")
		user, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var feeding models.Feeding
		feedingData, err := feeding.GetAll(env, &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		feedingResponse := FeedingResponse{FeedingData: feedingData}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(feedingResponse)
	})
}

//FeedingViewHandler -
func FeedingViewHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id := mux.Vars(r)["id"]
		fmt.Fprintf(w, "GET feeding with id %s", id)
	})
}

// FeedingEditHandler -
func FeedingEditHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id := mux.Vars(r)["id"]
		fmt.Fprintf(w, "PUT with id %s", id)
	})
}

//FeedingNewHandler -
func FeedingNewHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("POST Feeding")
		user, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		family, err := user.GetFamily(env)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
		feedingRequest.FeedingData.FamilyID = family.ID
		err = feedingRequest.FeedingData.Save(env)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusConflict)
		}
		json.NewEncoder(w).Encode(feedingRequest)
	})
}

//FeedingDeleteHandler -
func FeedingDeleteHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id := mux.Vars(r)["id"]
		fmt.Fprintf(w, "DELETE with id %s", id)
	})
}
