package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sasimpson/goparent/config"
	"github.com/sasimpson/goparent/models"
)

//UserRequest - structure for incoming user request
type UserRequest struct {
	UserData models.User `json:"userData"`
}

//NewUserRequest - this is for submitting password in new user request
type NewUserRequest struct {
	UserData struct {
		ID       string `json:"id,omitempty"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"userData"`
}

//UserAuthResponse - auth response structure
type UserAuthResponse struct {
	UserData models.User `json:"userData"`
	Token    string      `json:"token"`
}

func initUsersHandlers(env *config.Env, r *mux.Router) {
	u := r.PathPrefix("/user").Subrouter()
	u.Handle("/{id}", AuthRequired(userGetHandler(env), env)).Methods("GET").Name("UserView")
	u.Handle("/", userNewHandler(env)).Methods("POST").Name("UserNew")
	u.Handle("/login", loginHandler(env)).Methods("POST").Name("UserLogin")
}

func loginHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		username := r.FormValue("username")
		password := r.FormValue("password")

		log.Println("POST /api/user/login,", username, password)
		var user models.User
		err := user.GetUserByLogin(env, username, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		token, err := user.GetToken(env)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		var userResp UserAuthResponse
		userResp.UserData = user
		userResp.Token = token
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("x-auth-token", token)
		json.NewEncoder(w).Encode(userResp)
	})
}

func userGetHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// user, err := UserFromContext(r.Context())
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		vars := mux.Vars(r)
		log.Printf("mux vars: %v", vars)
		log.Println("GET  /api/user/", vars["id"])
		var lookupUser models.User
		err := lookupUser.GetUser(env, vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		// if user.ID == lookupUser.ID {
		json.NewEncoder(w).Encode(lookupUser)
		return
		// }
		// http.Error(w, "not authorized", http.StatusUnauthorized)
	})
}

func userNewHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("POST /api/user")
		decoder := json.NewDecoder(r.Body)
		var newUserRequest NewUserRequest
		err := decoder.Decode(&newUserRequest)
		userData := models.User(newUserRequest.UserData)
		if err != nil {
			log.Panicln(err)
		}
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		err = userData.Save(env)
		if err != nil {
			log.Println(err)
			var errMsg ErrService
			errMsg.ErrMessage.Body = err.Error()
			errMsg.ErrMessage.Code = http.StatusConflict
			js, _ := json.Marshal(errMsg)
			http.Error(w, string(js), http.StatusConflict)
			return
		}
		json.NewEncoder(w).Encode(userData)
	})
}
