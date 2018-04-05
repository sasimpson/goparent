package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sasimpson/goparent"
)

//WasteRequest - request structure for waste
type WasteRequest struct {
	WasteData goparent.Waste `json:"wasteData"`
}

//WasteResponse - response structure for waste
type WasteResponse struct {
	WasteData []*goparent.Waste `json:"wasteData"`
}

func (h *Handler) initWasteHandlers(r *mux.Router) {
	w := r.PathPrefix("/waste").Subrouter()
	w.Handle("", h.AuthRequired(h.wasteGetHandler())).Methods("GET").Name("WasteGet")
	w.Handle("", h.AuthRequired(h.wasteNewHandler())).Methods("POST").Name("WasteNew")
	w.Handle("/{id}", h.AuthRequired(h.wasteViewHandler())).Methods("GET").Name("WasteView")
	w.Handle("/{id}", h.AuthRequired(h.wasteEditHandler())).Methods("PUT").Name("WasteEdit")
	w.Handle("/{id}", h.AuthRequired(h.wasteDeleteHandler())).Methods("DELETE").Name("WasteDelete")
}

func (h *Handler) wasteGetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		family, err := h.UserService.GetFamily(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		wasteData, err := h.WasteService.Waste(family)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		wasteResponse := WasteResponse{WasteData: wasteData}
		w.Header().Set("Content-Type", jsonContentType)
		json.NewEncoder(w).Encode(wasteResponse)
	})
}

func (h *Handler) wasteViewHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		http.Error(w, "not implemented", http.StatusNotImplemented)

	})
}

func (h *Handler) wasteEditHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		http.Error(w, "not implemented", http.StatusNotImplemented)
	})
}

func (h *Handler) wasteNewHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		family, err := h.UserService.GetFamily(user)
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
		err = h.WasteService.Save(&wasteRequest.WasteData)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		json.NewEncoder(w).Encode(wasteRequest)
	})
}

func (h *Handler) wasteDeleteHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		http.Error(w, "not implemented", http.StatusNotImplemented)
	})
}
