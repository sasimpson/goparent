package api

import (
	"encoding/json"
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
	Pagination
	WasteData []*goparent.Waste `json:"wasteData"`
}

func (h *Handler) initWasteHandlers(r *mux.Router) {
	w := r.PathPrefix("/waste").Subrouter()
	w.Handle("", h.AuthRequired(h.wasteGetHandler())).Methods("GET").Name("WasteGet")
	w.Handle("", h.AuthRequired(h.wasteNewHandler())).Methods("POST").Name("WasteNew")
	w.Handle("/graph/{id}", h.AuthRequired(h.wasteGraphDataHandler())).Methods("GET").Name("WasteGraphData")
	w.Handle("/{id}", h.AuthRequired(h.wasteViewHandler())).Methods("GET").Name("WasteView")
	w.Handle("/{id}", h.AuthRequired(h.wasteEditHandler())).Methods("PUT").Name("WasteEdit")
	w.Handle("/{id}", h.AuthRequired(h.wasteDeleteHandler())).Methods("DELETE").Name("WasteDelete")
}

func (h *Handler) wasteGetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := h.Env.DB.GetContext(r)

		user, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		pagination := getPagination(r)

		family, err := h.UserService.GetFamily(ctx, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		wasteData, err := h.WasteService.Waste(ctx, family, pagination.Days)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		wasteResponse := WasteResponse{WasteData: wasteData, Pagination: *pagination}
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
		ctx := h.Env.DB.GetContext(r)

		user, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		family, err := h.UserService.GetFamily(ctx, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var wasteRequest WasteRequest
		err = decoder.Decode(&wasteRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		w.Header().Set("Content-Type", jsonContentType)
		wasteRequest.WasteData.UserID = user.ID
		wasteRequest.WasteData.FamilyID = family.ID
		err = h.WasteService.Save(ctx, &wasteRequest.WasteData)
		if err != nil {
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

func (h *Handler) wasteGraphDataHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := h.Env.DB.GetContext(r)

		user, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		childID := mux.Vars(r)["id"]

		family, err := h.UserService.GetFamily(ctx, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		child, err := h.ChildService.Child(ctx, childID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//child needs to belong to the user's family.
		if child.FamilyID != family.ID {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		wasteGraphData, err := h.WasteService.GraphData(ctx, child)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", jsonContentType)
		json.NewEncoder(w).Encode(wasteGraphData)
	})
}
