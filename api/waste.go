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
	w.Handle("", AuthRequired(WasteGetHandler(env), env)).Methods("GET").Name("WasteGet")
	w.Handle("", AuthRequired(WasteNewHandler(env), env)).Methods("POST").Name("WasteNew")
	w.Handle("/{id}", AuthRequired(WasteViewHandler(env), env)).Methods("GET").Name("WasteView")
	w.Handle("/{id}", AuthRequired(WasteEditHandler(env), env)).Methods("PUT").Name("WasteEdit")
	w.Handle("/{id}", AuthRequired(WasteDeleteHandler(env), env)).Methods("DELETE").Name("WasteDelete")
}

//WasteGetHandler -
func WasteGetHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET /api/waste")
		user, err := models.UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
		_, err := models.UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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
		_, err := models.UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id := mux.Vars(r)["id"]
		fmt.Fprintf(w, "PUT with id %s", id)
	})
}

//WasteNewHandler -
func WasteNewHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("POST Waste")
		user, err := models.UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
		_, err := models.UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id := mux.Vars(r)["id"]
		fmt.Fprintf(w, "DELETE with id %s", id)
	})
}
