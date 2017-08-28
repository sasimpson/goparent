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

type WasteRequest struct {
	WasteData models.Waste `json:"wasteData"`
}
type WasteResponse struct {
	WasteData []models.Waste `json:"wasteData"`
}

func initWasteHandlers(env *config.Env, r *mux.Router) {
	w := r.PathPrefix("/waste").Subrouter()
	w.Handle("", WasteGetHandler(env)).Methods("GET")
	w.Handle("", WasteNewHandler(env)).Methods("POST")
	w.Handle("/{id}", WasteViewHandler(env)).Methods("GET")
	w.Handle("/{id}", WasteEditHandler(env)).Methods("PUT")
	w.Handle("/{id}", WasteDeleteHandler(env)).Methods("DELETE")
}

//WasteGetHandler -
func WasteGetHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET /api/waste")
		user, err := validateAuthToken(env, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		var waste models.Waste
		wasteData, err := waste.GetAll(env, &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		wasteResponse := WasteResponse{WasteData: wasteData}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(wasteResponse.WasteData)
	})
}

//WasteViewHandler -
func WasteViewHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		var waste models.Waste
		waste.GetByID(env, id)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(waste)
	})
}

// WasteEditHandler -
func WasteEditHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		fmt.Fprintf(w, "PUT with id %s", id)
	})
}

//WasteNewHandler -
func WasteNewHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("POST Waste")
		user, err := validateAuthToken(env, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		decoder := json.NewDecoder(r.Body)
		var wasteRequest WasteRequest
		err = decoder.Decode(&wasteRequest)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		wasteRequest.WasteData.UserID = user.ID
		err = wasteRequest.WasteData.Save(env)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusConflict)
		}
		json.NewEncoder(w).Encode(wasteRequest.WasteData)
	})
}

//WasteDeleteHandler -
func WasteDeleteHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		fmt.Fprintf(w, "DELETE with id %s", id)
	})
}
