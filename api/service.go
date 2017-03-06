package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

//RunService - Runs service interfaces for app
func RunService() {
	r := mux.NewRouter()
	a := r.PathPrefix("/api").Subrouter()
	a.HandleFunc("/", apiHandler)

	initFeedingHandlers(a)
	initSleepHandlers(a)
	initWasteHandlers(a)

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "please check docs at: http://doesntexist.yet")
}
