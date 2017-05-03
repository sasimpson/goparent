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

type SleepResponse struct {
	SleepData []models.Sleep `json:"sleepData"`
}

func initSleepHandlers(r *mux.Router) {
	s := r.PathPrefix("/sleep").Subrouter()
	s.HandleFunc("", sleepGetHandler).Methods("GET")
	s.HandleFunc("", sleepNewHandler).Methods("POST")
	s.HandleFunc("/status", sleepToggleStatus).Methods("GET")
	s.HandleFunc("/start", sleepStartHandler).Methods("GET")
	s.HandleFunc("/end", sleepEndHandler).Methods("GET")
	s.HandleFunc("/{id}", sleepViewHandler).Methods("GET")
	s.HandleFunc("/{id}", sleepEditHandler).Methods("PUT")
	s.HandleFunc("/{id}", sleepDeleteHandler).Methods("DELETE")
}

func sleepGetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GET sleep")
	var sleep models.Sleep
	sleepData, err := sleep.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sleepResponse := SleepResponse{SleepData: sleepData}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sleepResponse.SleepData)
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
		return
	}
	json.NewEncoder(w).Encode(sleepRequest.SleepData)
}

func sleepDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprintf(w, "DELETE with id %s", id)
}

func sleepStartHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GET sleep start")
	var sleep models.Sleep
	err := sleep.Start()
	if err != nil {
		log.Println(err.Error())
		if err == models.ExistingStartErr {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sleep.Save()
	fmt.Fprintf(w, "started Sleep")
	return
}

func sleepEndHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GET sleep end")
	var sleep models.Sleep
	err := sleep.End()
	if err != nil {
		if err == models.NoExistingSessionErr {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sleep.Save()
	fmt.Fprintf(w, "ended Sleep")
}

func sleepToggleStatus(w http.ResponseWriter, r *http.Request) {
	log.Println("GET sleep toggle status")
	var sleep models.Sleep
	ok, err := sleep.Status()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if ok {
		fmt.Fprintf(w, "sleep session active")
		return
	}
	http.Error(w, "not found", http.StatusNotFound)
	return
}
