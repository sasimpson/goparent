package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sasimpson/goparent"
)

//ChildrenResponse - response for children lists
type ChildrenResponse struct {
	Children []*goparent.Child `json:"children"`
}

//ChildRequest - incoming request structure
type ChildRequest struct {
	ChildData goparent.Child `json:"childData"`
}

//ChildDeletedResponse - response of deleted child
type ChildDeletedResponse struct {
	Deleted int `json:"deleted"`
}

//ChildSummaryResponse - Summary response
type ChildSummaryResponse struct {
	ChildData goparent.Child `json:"childData"`
	Stats     Summary        `json:"stats"`
}

//Summary - return structure of all summary data
type Summary struct {
	Feeding goparent.FeedingSummary `json:"feeding"`
	Sleep   goparent.SleepSummary   `json:"sleep"`
	Waste   goparent.WasteSummary   `json:"waste"`
}

func (h *Handler) initChildrenHandlers(r *mux.Router) {
	c := r.PathPrefix("/children").Subrouter()
	c.Handle("", h.AuthRequired(h.childrenGetHandler())).Methods("GET").Name("ChildrenGet")
	c.Handle("", h.AuthRequired(h.childNewHandler())).Methods("POST").Name("ChildNew")
	c.Handle("/{id}", h.AuthRequired(h.childViewHandler())).Methods("GET").Name("ChildView")
	c.Handle("/{id}", h.AuthRequired(h.childEditHandler())).Methods("PUT").Name("ChildEdit")
	c.Handle("/{id}", h.AuthRequired(h.childDeleteHandler())).Methods("DELETE").Name("ChildDelete")
	c.Handle("/{id}/summary", h.AuthRequired(h.childSummary())).Methods("GET").Name("ChildSummary")

}

//ChildSummary - handler to assemble and reuturn child summary data
func (h *Handler) childSummary() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := h.Env.DB.GetContext(r)
		childID := mux.Vars(r)["id"]
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

		var summary ChildSummaryResponse
		child, err := h.ChildService.Child(ctx, childID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		summary.ChildData = *child
		if child.FamilyID != family.ID {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		feedingSummary, err := h.FeedingService.Stats(ctx, child)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		summary.Stats.Feeding = *feedingSummary

		sleeps, err := h.SleepService.Stats(ctx, child)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		summary.Stats.Sleep = *sleeps

		wastes, err := h.WasteService.Stats(ctx, child)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		summary.Stats.Waste = *wastes

		w.Header().Set("Content-Type", jsonContentType)
		json.NewEncoder(w).Encode(summary)
	})
}

//childrenGetHandler - GET / - gets all children for user
func (h *Handler) childrenGetHandler() http.Handler {
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

		allChildren, err := h.FamilyService.Children(ctx, family)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		childrenResponse := ChildrenResponse{Children: allChildren}
		w.Header().Set("Content-Type", jsonContentType)
		json.NewEncoder(w).Encode(childrenResponse)
	})
}

//ChildNewHandler - POST / - create a new child for a user
func (h *Handler) childNewHandler() http.Handler {
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
		var childRequest ChildRequest
		err = decoder.Decode(&childRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer r.Body.Close()
		w.Header().Set("Content-Type", jsonContentType)
		childRequest.ChildData.ParentID = user.ID
		childRequest.ChildData.FamilyID = family.ID
		err = h.ChildService.Save(ctx, &childRequest.ChildData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		json.NewEncoder(w).Encode(childRequest.ChildData)
	})
}

//ChildViewHandler - GET /{id} - gets the data for a child for a user
func (h *Handler) childViewHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := h.Env.DB.GetContext(r)

		childID := mux.Vars(r)["id"]
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

		w.Header().Set("Content-Type", jsonContentType)
		json.NewEncoder(w).Encode(child)
	})
}

// ChildEditHandler - PUT /{id} - edit a child for a user
func (h *Handler) childEditHandler() http.Handler {
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
		var childRequest ChildRequest
		err = decoder.Decode(&childRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id := mux.Vars(r)["id"]
		child, err := h.ChildService.Child(ctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		//verify both the child we requested to edit, and that the parent is the user id.
		if (child.ID != childRequest.ChildData.ID) || (childRequest.ChildData.FamilyID != family.ID) {
			http.Error(w, "invalid relationship", http.StatusBadRequest)
			return
		}
		h.ChildService.Save(ctx, &childRequest.ChildData)
		err = json.NewEncoder(w).Encode(childRequest.ChildData)
		return
	})
}

//ChildDeleteHandler - DELETE /{id} - delete a child for a user
//TODO - need to delete or archive child and child's data.
func (h *Handler) childDeleteHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := h.Env.DB.GetContext(r)

		user, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		family, err := h.UserService.GetFamily(ctx, user)
		id := mux.Vars(r)["id"]
		child, err := h.ChildService.Child(ctx, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if child.FamilyID != family.ID {
			http.Error(w, "not found", http.StatusNotFound)
		}

		deleted, err := h.ChildService.Delete(ctx, child)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var deletedResponse ChildDeletedResponse

		deletedResponse.Deleted = deleted

		w.Header().Set("Content-Type", jsonContentType)
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(deletedResponse)
	})
}
