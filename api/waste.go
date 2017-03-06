package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func initWasteHandlers(r *mux.Router) {
	w := r.PathPrefix("/waste").Subrouter()
	w.HandleFunc("", WasteGetHandler).Methods("GET")
	w.HandleFunc("", WasteNewHandler).Methods("POST")
	w.HandleFunc("/{id}", WasteViewHandler).Methods("GET")
	w.HandleFunc("/{id}", WasteEditHandler).Methods("PUT")
	w.HandleFunc("/{id}", WasteDeleteHandler).Methods("DELETE")
}

//WasteGetHandler -
func WasteGetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GET no id")
}

//WasteViewHandler -
func WasteViewHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprintf(w, "GET with id %s", id)
}

// WasteEditHandler -
func WasteEditHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprintf(w, "PUT with id %s", id)
}

//WasteNewHandler -
func WasteNewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "POST with data:")
}

//WasteDeleteHandler -
func WasteDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprintf(w, "DELETE with id %s", id)
}
