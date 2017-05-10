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
type UserAuthResponse struct {
	UserData models.User `json:"userData"`
	Token    string      `json:"token"`
}

func initUsersHandlers(r *mux.Router) {
	u := r.PathPrefix("/user").Subrouter()
	u.HandleFunc("/{id}", userGetHandler).Methods("GET")
	u.HandleFunc("/", userNewHandler).Methods("POST")
	u.HandleFunc("/login", loginHandler).Methods("POST")
	u.HandleFunc("/validate", validateUserTokenHandler).Methods("POST")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")
	log.Println("POST /api/user/login", username, password)
	var user models.User
	err := user.GetUserByLogin(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	token, err := user.GetToken(mySigningKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var userResp UserAuthResponse
	userResp.UserData = user
	userResp.Token = token
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("x-auth-token", token)
	json.NewEncoder(w).Encode(userResp)
}

func userGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Println("GET  /api/user/", vars["id"])
	var user models.User
	err := user.GetUser(vars["id"])
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

func validateUserTokenHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /api/user/validate")
	tokenString := r.Header.Get("x-auth-token")
	var user models.User
	token, err := user.ValidateToken(tokenString, mySigningKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if token {
		w.Header().Set("x-auth-token", tokenString)
		w.WriteHeader(http.StatusAccepted)
		return
	}
	http.Error(w, err.Error(), http.StatusUnauthorized)
	return
}
