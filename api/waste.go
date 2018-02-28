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

//WasteRequest - request structure for waste
type WasteRequest struct {
	WasteData models.Waste `json:"wasteData"`
}

//WasteResponse - response structure for waste
type WasteResponse struct {
	WasteData []models.Waste `json:"wasteData"`
}

func initWasteHandlers(env *config.Env, r *mux.Router) {
	w := r.PathPrefix("/waste").Subrouter()
	w.Handle("", AuthRequired(wasteGetHandler(env), env)).Methods("GET").Name("WasteGet")
	w.Handle("", AuthRequired(wasteNewHandler(env), env)).Methods("POST").Name("WasteNew")
	w.Handle("/{id}", AuthRequired(wasteViewHandler(env), env)).Methods("GET").Name("WasteView")
	w.Handle("/{id}", AuthRequired(wasteEditHandler(env), env)).Methods("PUT").Name("WasteEdit")
	w.Handle("/{id}", AuthRequired(wasteDeleteHandler(env), env)).Methods("DELETE").Name("WasteDelete")
}

func wasteGetHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		var waste models.Waste
		wasteData, err := waste.GetAll(env, &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		wasteResponse := WasteResponse{WasteData: wasteData}
		w.Header().Set("Content-Type", jsonContentType)
		json.NewEncoder(w).Encode(wasteResponse)
	})
}

func wasteViewHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		id := mux.Vars(r)["id"]
		var waste models.Waste
		waste.GetByID(env, id)
		w.Header().Set("Content-Type", jsonContentType)
		json.NewEncoder(w).Encode(waste)
	})
}

func wasteEditHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		id := mux.Vars(r)["id"]
		fmt.Fprintf(w, "PUT with id %s", id)
	})
}

func wasteNewHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		family, err := user.GetFamily(env)
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
			return
		}
		defer r.Body.Close()

		w.Header().Set("Content-Type", jsonContentType)
		wasteRequest.WasteData.UserID = user.ID
		wasteRequest.WasteData.FamilyID = family.ID
		err = wasteRequest.WasteData.Save(env)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		json.NewEncoder(w).Encode(wasteRequest.WasteData)
	})
}

func wasteDeleteHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		id := mux.Vars(r)["id"]
		fmt.Fprintf(w, "DELETE with id %s", id)
	})
}
