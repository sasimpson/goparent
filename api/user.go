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

//UserResponse - structure for responding to user info requests
type UserResponse struct {
	UserData   *models.User   `json:"userData"`
	FamilyData *models.Family `json:"familyData"`
}

//NewUserRequest - this is for submitting password in new user request
type NewUserRequest struct {
	UserData struct {
		ID            string `json:"id,omitempty"`
		Name          string `json:"name"`
		Email         string `json:"email"`
		Username      string `json:"username"`
		Password      string `json:"password"`
		CurrentFamily string `json:"currentFamily"`
	} `json:"userData"`
}

//UserAuthResponse - auth response structure
type UserAuthResponse struct {
	UserData models.User `json:"userData"`
	Token    string      `json:"token"`
}

//InvitesResponse - response structure for invites
type InvitesResponse struct {
	SentInviteData    []models.UserInvitation `json:"sentInviteData"`
	PendingInviteData []models.UserInvitation `json:"pendingInviteData`
}

func initUsersHandlers(env *config.Env, r *mux.Router) {
	u := r.PathPrefix("/user").Subrouter()
	// u.Handle("/{id}", AuthRequired(userGetHandler(env), env)).Methods("GET").Name("UserView")
	u.Handle("/", userNewHandler(env)).Methods("POST").Name("UserNew")
	u.Handle("/", AuthRequired(userGetHandler(env), env)).Methods("GET").Name("UserGetData")
	u.Handle("/login", loginHandler(env)).Methods("POST").Name("UserLogin")
	u.Handle("/invite", AuthRequired(userListInviteHandler(env), env)).Methods("GET").Name("UserGetSentInvites")
	u.Handle("/invite", AuthRequired(userNewInviteHandler(env), env)).Methods("POST").Name("UserNewInvite")
	u.Handle("/invite/{id}", AuthRequired(userDeleteInviteHandler(env), env)).Methods("DELETE").Name("UserDeleteInvite")
	u.Handle("/invite/accept/{id}", AuthRequired(userAcceptInviteHandler(env), env)).Methods("POST").Name("UserAcceptInvite")
}

func loginHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		username := r.FormValue("username")
		password := r.FormValue("password")

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
		w.Header().Set("Content-Type", jsonContentType)
		w.Header().Set("x-auth-token", token)
		json.NewEncoder(w).Encode(userResp)
	})
}

func userGetHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		family, _ := user.GetFamily(env)
		userInfo := UserResponse{
			UserData:   &user,
			FamilyData: &family}
		w.Header().Set("Content-Type", jsonContentType)
		json.NewEncoder(w).Encode(userInfo)
		return
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

		w.Header().Set("Content-Type", jsonContentType)
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

func userNewInviteHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		err = r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		invitedUserEmail := r.PostFormValue("email")
		if invitedUserEmail == "" || len(invitedUserEmail) <= 0 {
			http.Error(w, "no invite email submitted", http.StatusBadRequest)
			return
		}

		err = user.InviteParent(env, invitedUserEmail)
		if err != nil {
			if err.Error() == models.ErrExistingInvitation {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})
}

func userListInviteHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		sentInvites, err := user.GetSentInvites(env)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pendingInvites, err := user.GetInvites(env)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		invitesResponse := InvitesResponse{SentInviteData: sentInvites, PendingInviteData: pendingInvites}
		w.Header().Set("Content-Type", jsonContentType)
		json.NewEncoder(w).Encode(invitesResponse)
	})
}

func userAcceptInviteHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		id := mux.Vars(r)["id"]
		err = user.AcceptInvite(env, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		w.Header().Set("Content-Type", jsonContentType)
		return
	})
}

func userDeleteInviteHandler(env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := UserFromContext(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		id := mux.Vars(r)["id"]
		err = user.DeleteInvite(env, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", jsonContentType)
		w.WriteHeader(http.StatusNoContent)
	})
}
