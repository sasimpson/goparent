package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sasimpson/goparent/config"
	"github.com/sasimpson/goparent/models"
)

//ChildrenResponse - response for children lists
type ChildrenResponse struct {
	Children []models.Child `json:"children"`
}

//ChildRequest - incoming request structure
type ChildRequest struct {
	ChildData models.Child `json:"childData"`
}

//ChildDeletedResponse - response of deleted child
type ChildDeletedResponse struct {
	Deleted int `json:"deleted"`
}

//ChildSummaryResponse - Summary response
type ChildSummaryResponse struct {
	ChildData models.Child `json:"childData"`
	Stats     Summary      `json:"stats"`
}

//Summary - return structure of all summary data
type Summary struct {
	Feeding models.FeedingSummary `json:"feeding"`
	Sleep   models.SleepSummary   `json:"sleep"`
	Waste   models.WasteSummary   `json:"waste"`
}

func initChildrenHandlers(env *config.Env, r *mux.Router) {
	c := r.PathPrefix("/children").Subrouter()
	c.Handle("", AuthRequired(ChildrenGetHandler(env), env)).Methods("GET").Name("ChildrenGet")
	c.Handle("", AuthRequired(ChildNewHandler(env), env)).Methods("POST").Name("ChildNew")
	c.Handle("/{id}", AuthRequired(ChildViewHandler(env), env)).Methods("GET").Name("ChildView")
	c.Handle("/{id}", AuthRequired(ChildEditHandler(env), env)).Methods("PUT").Name("ChildEdit")
	c.Handle("/{id}", AuthRequired(ChildDeleteHandler(env), env)).Methods("DELETE").Name("ChildDelete")
	c.Handle("/{id}/summary", AuthRequired(ChildSummary(env), env)).Methods("GET").Name("ChildSummary")

}

//ChildSummary - handler to assemble and reuturn child summary data
func ChildSummary(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		childID := mux.Vars(r)["id"]
		log.Println("Child Summary: ", childID)
		user, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		var summary ChildSummaryResponse
		var child models.Child
		err = child.GetChild(env, &user, childID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		summary.ChildData = child

		feedingSummary, err := models.FeedingGetStats(env, &user, &child)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		summary.Stats.Feeding = feedingSummary

		sleeps, err := models.SleepGetStats(env, &user, &child)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		summary.Stats.Sleep = sleeps

		wastes, err := models.WasteGetStats(env, &user, &child)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		summary.Stats.Waste = wastes

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(summary)
	})
}

//ChildrenGetHandler - GET / - gets all children for user
func ChildrenGetHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Children GET ")
		user, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		allChildren, err := models.GetAllChildren(env, &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		childrenResponse := ChildrenResponse{Children: allChildren}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(childrenResponse)
	})
}

//ChildNewHandler - POST / - create a new child for a user
func ChildNewHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("POST Child")
		user, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var childRequest ChildRequest
		err = decoder.Decode(&childRequest)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		childRequest.ChildData.ParentID = user.ID
		err = childRequest.ChildData.Save(env)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		json.NewEncoder(w).Encode(childRequest.ChildData)
	})
}

//ChildViewHandler - GET /{id} - gets the data for a child for a user
func ChildViewHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		childID := mux.Vars(r)["id"]
		log.Println("Child View: ", childID)
		user, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var child models.Child
		err = child.GetChild(env, &user, childID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(child)
	})
}

// ChildEditHandler - PUT /{id} - edit a child for a user
func ChildEditHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("PUT Child")
		user, err := UserFromContext(r.Context())
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
		var child models.Child
		err = child.GetChild(env, &user, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//verify both the child we requested to edit, and that the parent is the user id.
		if (child.ID != childRequest.ChildData.ID) || (childRequest.ChildData.ParentID != user.ID) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		childRequest.ChildData.Save(env)
		err = json.NewEncoder(w).Encode(childRequest.ChildData)
		return
	})
}

//ChildDeleteHandler - DELETE /{id} - delete a child for a user
//TODO - need to delete or archive child and child's data.
func ChildDeleteHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		id := mux.Vars(r)["id"]
		var child models.Child
		err = child.GetChild(env, &user, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		deleted, err := child.DeleteChild(env, &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var deletedResponse ChildDeletedResponse

		deletedResponse.Deleted = deleted

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(deletedResponse)
	})
}

// func RandomData(env *config.Env) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		user, err := UserFromContext(r.Context())
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		id := mux.Vars(r)["id"]
// 		var child models.Child
// 		err = child.GetChild(env, &user, id)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusNotFound)
// 			return
// 		}

// 	})
// }
