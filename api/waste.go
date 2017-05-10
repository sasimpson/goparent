package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sasimpson/goparent/models"
)

type WasteRequest struct {
	WasteData models.Waste `json:"wasteData"`
}
type WasteResponse struct {
	WasteData []models.Waste `json:"wasteData"`
}

func initWasteHandlers(r *mux.Router) {
	w := r.PathPrefix("/waste").Subrouter()
	w.HandleFunc("", WasteGetHandler).Methods("GET")
	w.HandleFunc("", WasteNewHandler).Methods("POST")
	w.HandleFunc("/{id}", WasteViewHandler).Methods("GET")
	w.HandleFunc("/{id}", WasteEditHandler).Methods("PUT")
	w.HandleFunc("/{id}", WasteDeleteHandler).Methods("DELETE")
}

//WasteGetHandler -
func WasteGetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /api/waste")
	user, err := validateAuthToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var waste models.Waste
	wasteData, err := waste.GetAll(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	wasteResponse := WasteResponse{WasteData: wasteData}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wasteResponse.WasteData)
}

//WasteViewHandler -
func WasteViewHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var waste models.Waste
	waste.GetByID(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(waste)
}

// WasteEditHandler -
func WasteEditHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprintf(w, "PUT with id %s", id)
}

//WasteNewHandler -
func WasteNewHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("POST Waste")
	user, err := validateAuthToken(r)
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
	err = wasteRequest.WasteData.Save()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusConflict)
	}
	json.NewEncoder(w).Encode(wasteRequest.WasteData)
}

//WasteDeleteHandler -
func WasteDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprintf(w, "DELETE with id %s", id)
}
