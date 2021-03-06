package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sasimpson/goparent"
)

//FeedingRequest - request structure for feedings
type FeedingRequest struct {
	FeedingData goparent.Feeding `json:"feedingData"`
}

//FeedingResponse - response structure for feedings
type FeedingResponse struct {
	FeedingData []*goparent.Feeding `json:"feedingData"`
}

func (h *Handler) initFeedingHandlers(r *mux.Router) {
	f := r.PathPrefix("/feeding").Subrouter()
	f.Handle("", h.AuthRequired(h.feedingGetHandler())).Methods("GET").Name("FeedingGet")
	f.Handle("", h.AuthRequired(h.feedingNewHandler())).Methods("POST").Name("FeedingNew")
	f.Handle("/graph/{id}", h.AuthRequired(h.feedingGraphDataHandler())).Methods("GET").Name("FeedingGraphData")
	f.Handle("/{id}", h.AuthRequired(h.feedingViewHandler())).Methods("GET").Name("FeedingView")
	f.Handle("/{id}", h.AuthRequired(h.feedingEditHandler())).Methods("PUT").Name("FeedingEdit")
	f.Handle("/{id}", h.AuthRequired(h.feedingDeleteHandler())).Methods("DELETE").Name("FeedingDelete")
}

func (h *Handler) feedingGetHandler() http.Handler {
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

		feedingData, err := h.FeedingService.Feeding(ctx, family, pagination.Days)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		feedingResponse := FeedingResponse{FeedingData: feedingData}
		// log.Printf("feedingResponse %#v", feedingResponse)
		w.Header().Set("Content-Type", jsonContentType)
		json.NewEncoder(w).Encode(feedingResponse)
	})
}

func (h *Handler) feedingViewHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		http.Error(w, "not implemented", http.StatusNotImplemented)
	})
}

func (h *Handler) feedingEditHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		http.Error(w, "not implemented", http.StatusNotImplemented)
	})
}

func (h *Handler) feedingNewHandler() http.Handler {
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
		var feedingRequest FeedingRequest
		err = decoder.Decode(&feedingRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer r.Body.Close()

		w.Header().Set("Content-Type", jsonContentType)
		feedingRequest.FeedingData.UserID = user.ID
		feedingRequest.FeedingData.FamilyID = family.ID
		err = h.FeedingService.Save(ctx, &feedingRequest.FeedingData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
		}
		json.NewEncoder(w).Encode(feedingRequest)
	})
}

func (h *Handler) feedingDeleteHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, "not implemented", http.StatusNotImplemented)
	})
}

func (h *Handler) feedingGraphDataHandler() http.Handler {
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

		feedingGraphData, err := h.FeedingService.GraphData(ctx, child)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", jsonContentType)
		json.NewEncoder(w).Encode(feedingGraphData)
	})
}
