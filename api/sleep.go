package api

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/sasimpson/goparent/config"
	"github.com/sasimpson/goparent/models"
)

type SleepRequest struct {
	SleepData models.Sleep `json:"sleepData"`
}

type SleepResponse struct {
	SleepData []models.Sleep `json:"sleepData"`
}

func initSleepHandlers(env *config.Env, r *mux.Router) {
	s := r.PathPrefix("/sleep").Subrouter()
	s.Handle("", sleepGetHandler(env)).Methods("GET")
	s.Handle("", sleepNewHandler(env)).Methods("POST")
	s.Handle("/status", sleepToggleStatus(env)).Methods("GET")
	s.Handle("/start", sleepStartHandler(env)).Methods("POST")
	s.Handle("/end", sleepEndHandler(env)).Methods("POST")
	s.Handle("/{id}", sleepViewHandler(env)).Methods("GET")
	s.Handle("/{id}", sleepEditHandler(env)).Methods("PUT")
	s.Handle("/{id}", sleepDeleteHandler(env)).Methods("DELETE")
}

func sleepGetHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET sleep")
		user, err := validateAuthToken(env, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		var sleep models.Sleep
		sleepData, err := sleep.GetAll(env, &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sleepResponse := SleepResponse{SleepData: sleepData}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sleepResponse.SleepData)
	})
}

func sleepViewHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		fmt.Fprintf(w, "GET with id %s", id)
	})
}

func sleepEditHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		fmt.Fprintf(w, "PUT with id %s", id)
	})
}

func sleepNewHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//how time should be passed "2017-03-09T18:09:31.409Z"
		log.Println("PUT sleep")
		user, err := validateAuthToken(env, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		decoder := json.NewDecoder(r.Body)
		var sleepRequest SleepRequest
		err = decoder.Decode(&sleepRequest)
		if err != nil {
			log.Panicln(err)
		}
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		sleepRequest.SleepData.UserID = user.ID
		err = sleepRequest.SleepData.Save(env)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		json.NewEncoder(w).Encode(sleepRequest.SleepData)
	})
}

func sleepDeleteHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		fmt.Fprintf(w, "DELETE with id %s", id)
	})
}

func sleepStartHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET sleep start")
		user, err := validateAuthToken(env, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		var sleep models.Sleep
		err = sleep.Start(env, &user)
		if err != nil {
			log.Println(err.Error())
			if err == models.ExistingStartErr {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sleep.UserID = user.ID
		sleep.Save(env)
		fmt.Fprintf(w, "started Sleep")
		return
	})
}

func sleepEndHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET sleep end")
		user, err := validateAuthToken(env, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		var sleep models.Sleep
		err = sleep.End(env, &user)
		if err != nil {
			if err == models.NoExistingSessionErr {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sleep.UserID = user.ID
		sleep.Save(env)
		fmt.Fprintf(w, "ended Sleep")
	})
}

func sleepToggleStatus(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET sleep toggle status")
		user, err := validateAuthToken(env, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		var sleep models.Sleep
		ok, err := sleep.Status(env, &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if ok {
			fmt.Fprintf(w, "sleep session active")
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
		return
	})
}
