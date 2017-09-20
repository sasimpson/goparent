package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"encoding/json"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sasimpson/goparent/config"
	"github.com/sasimpson/goparent/models"
)

type ServiceInfo struct {
	Version string `json:"version"`
}

//RunService - Runs service interfaces for app
func RunService(env *config.Env) {
	r := mux.NewRouter()
	a := r.PathPrefix("/api").Subrouter()
	a.HandleFunc("/", apiHandler)
	a.HandleFunc("/info", infoHandler)

	initUsersHandlers(env, a)
	initFeedingHandlers(env, a)
	initSleepHandlers(env, a)
	initWasteHandlers(env, a)

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

//AuthRequired - handler to handle authentication of users tokens.
func AuthRequired(h http.Handler, env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//switch this to auth bearer header
		tokenString := r.Header.Get("x-auth-token")
		var user models.User
		_, err := user.ValidateToken(env, tokenString)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}
