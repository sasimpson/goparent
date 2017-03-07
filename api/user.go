package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sasimpson/goparent/models"
)

type UserGetRequest struct {
	UserData models.User `json:"userData"`
}

func initUsersHandlers(r *mux.Router) {
	u := r.PathPrefix("/user").Subrouter()
	u.HandleFunc("/", userGetHandler).Methods("GET")
	u.HandleFunc("/", userNewHandler).Methods("POST")
}

func userGetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GET no id")
}

func userNewHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var userRequest UserGetRequest
	err := decoder.Decode(&userRequest)
	if err != nil {
		log.Panicln(err)
	}
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	err = userRequest.UserData.Save()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusConflict)
	}
	// user := models.GetUser(userRequest.UserData.ID)
	json.NewEncoder(w).Encode(userRequest.UserData)
}
