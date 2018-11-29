package api

import (
	"fmt"
	"net/http"

	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/sasimpson/goparent"
)

//SleepRequest - request structure for sleep
type SleepRequest struct {
	SleepData goparent.Sleep `json:"sleepData"`
}

//SleepResponse - response structure for sleep
type SleepResponse struct {
	SleepData []*goparent.Sleep `json:"sleepData"`
}

func (h *Handler) initSleepHandlers(r *mux.Router) {
	s := r.PathPrefix("/sleep").Subrouter()
	s.Handle("", h.AuthRequired(h.sleepGetHandler())).Methods("GET").Name("SleepGet")
	s.Handle("", h.AuthRequired(h.sleepNewHandler())).Methods("POST").Name("SleepNew")
	s.Handle("/{id}", h.AuthRequired(h.sleepViewHandler())).Methods("GET").Name("SleepView")
	s.Handle("/{id}", h.AuthRequired(h.sleepEditHandler())).Methods("PUT").Name("SleepEdit")
	s.Handle("/{id}", h.AuthRequired(h.sleepDeleteHandler())).Methods("DELETE").Name("SleepDelete")
	s.Handle("/status/{childID}", h.AuthRequired(h.sleepToggleStatus())).Methods("GET").Name("SleepStatus")
	s.Handle("/start/{childID}", h.AuthRequired(h.sleepStartHandler())).Methods("POST").Name("SleepStart")
	s.Handle("/end/{childID}", h.AuthRequired(h.sleepEndHandler())).Methods("POST").Name("SleepEnd")
	s.Handle("/graph/{childID}", h.AuthRequired(h.sleepGraphDataHandler())).Methods("GET").Name("SleepGraphData")

}

func (h *Handler) sleepGetHandler() http.Handler {
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

		sleepData, err := h.SleepService.Sleep(ctx, family, 7)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sleepResponse := SleepResponse{SleepData: sleepData}
		w.Header().Set("Content-Type", jsonContentType)
		json.NewEncoder(w).Encode(sleepResponse)
	})
}

func (h *Handler) sleepViewHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		http.Error(w, "not implemented", http.StatusNotImplemented)
	})
}

func (h *Handler) sleepEditHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		http.Error(w, "not implemented", http.StatusNotImplemented)
	})
}

func (h *Handler) sleepNewHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := h.Env.DB.GetContext(r)

		//how time should be passed "2017-03-09T18:09:31.409Z"
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
		var sleepRequest SleepRequest
		err = decoder.Decode(&sleepRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		w.Header().Set("Content-Type", jsonContentType)
		sleepRequest.SleepData.UserID = user.ID
		sleepRequest.SleepData.FamilyID = family.ID
		err = h.SleepService.Save(ctx, &sleepRequest.SleepData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		json.NewEncoder(w).Encode(sleepRequest)
	})
}

func (h *Handler) sleepDeleteHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		http.Error(w, "not implemented", http.StatusNotImplemented)
	})
}

func (h *Handler) sleepStartHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := h.Env.DB.GetContext(r)
		user, err := UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		childID := mux.Vars(r)["childID"]
		if childID == "" {
			http.Error(w, "invalid Child ID", http.StatusBadRequest)
			return
		}

		family, err := h.FamilyService.Family(ctx, user.CurrentFamily)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		child, err := h.ChildService.Child(ctx, childID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = h.SleepService.Start(ctx, family, child)
		if err != nil {
			if err == goparent.ErrExistingStart {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "started Sleep")
		return
	})
}

func (h *Handler) sleepEndHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := h.Env.DB.GetContext(r)
		user, err := UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		childID := mux.Vars(r)["childID"]
		if childID == "" {
			http.Error(w, "invalid Child ID", http.StatusBadRequest)
			return
		}

		family, err := h.FamilyService.Family(ctx, user.CurrentFamily)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		child, err := h.ChildService.Child(ctx, childID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = h.SleepService.End(ctx, family, child)
		if err != nil {
			if err == goparent.ErrNoExistingSession {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		fmt.Fprintf(w, "ended Sleep")
		return
	})
}

func (h *Handler) sleepToggleStatus() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := h.Env.DB.GetContext(r)
		user, err := UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		childID := mux.Vars(r)["childID"]
		if childID == "" {
			http.Error(w, "invalid Child ID", http.StatusBadRequest)
			return
		}

		family, err := h.FamilyService.Family(ctx, user.CurrentFamily)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		child, err := h.ChildService.Child(ctx, childID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, ok, err := h.SleepService.Status(ctx, family, child)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if ok {
			fmt.Fprintf(w, "sleep session active")
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
		return
	})
}

func (h *Handler) sleepGraphDataHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := h.Env.DB.GetContext(r)

		user, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		childID := mux.Vars(r)["childID"]

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

		sleepGraphData, err := h.SleepService.GraphData(ctx, child)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", jsonContentType)
		json.NewEncoder(w).Encode(sleepGraphData)
	})
}
