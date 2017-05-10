package api

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sasimpson/goparent/models"
)

type ServiceInfo struct {
	Version string `json:"version"`
}

var mySigningKey = []byte("supersecretsquirrl")

//RunService - Runs service interfaces for app
func RunService() {
	r := mux.NewRouter()
	a := r.PathPrefix("/api").Subrouter()
	a.HandleFunc("/", apiHandler)
	a.HandleFunc("/info", infoHandler)

	initUsersHandlers(a)
	initFeedingHandlers(a)
	initSleepHandlers(a)
	initWasteHandlers(a)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Accept", "Content-Type", "x-auth-token"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Println("starting service on port 8000")
	http.Handle("/", r)
	http.ListenAndServe(":8000", handlers.CORS(originsOk, headersOk, methodsOk)(r))
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "please check docs at: https://github.com/sasimpson/goparent")
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	si := ServiceInfo{Version: "v0.1"}
	json.NewEncoder(w).Encode(si)
	return
}

func validateAuthToken(r *http.Request) (models.User, error) {
	tokenString := r.Header.Get("x-auth-token")
	var user models.User
	_, err := user.ValidateToken(tokenString, mySigningKey)
	return user, err
}
