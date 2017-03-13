package api

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/sasimpson/goparent/models"
)

type SleepRequest struct {
	SleepData models.Sleep `json:"sleepData"`
}

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
	//how time should be passed "2017-03-09T18:09:31.409Z"
	decoder := json.NewDecoder(r.Body)
	var sleepRequest SleepRequest
	err := decoder.Decode(&sleepRequest)
	if err != nil {
		log.Panicln(err)
	}
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	err = sleepRequest.SleepData.Save()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusConflict)
	}
	json.NewEncoder(w).Encode(sleepRequest.SleepData)
}

func sleepDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprintf(w, "DELETE with id %s", id)
}
