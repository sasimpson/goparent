package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sasimpson/goparent/models"
)

//UserRequest - structure for incoming user request
type UserRequest struct {
	UserData models.User `json:"userData"`
}

func initUsersHandlers(r *mux.Router) {
	u := r.PathPrefix("/user").Subrouter()
	u.HandleFunc("/{id}", userGetHandler).Methods("GET")
	u.HandleFunc("/", userNewHandler).Methods("POST")
	u.HandleFunc("/login", loginHandler).Methods("POST")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")
	log.Println("POST /api/user/login", username, password)
	user, err := models.GetUserByLogin(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func userGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println("GET  /api/user/", vars["id"])
	user, err := models.GetUser(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func userNewHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /api/user")
	decoder := json.NewDecoder(r.Body)
	var userRequest UserRequest
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
		return
	}
	json.NewEncoder(w).Encode(userRequest.UserData)
}
