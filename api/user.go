package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sasimpson/goparent"
)

//UserRequest - structure for incoming user request
type UserRequest struct {
	UserData *goparent.User `json:"userData"`
}

//UserResponse - structure for responding to user info requests
type UserResponse struct {
	UserData   *goparent.User   `json:"userData"`
	FamilyData *goparent.Family `json:"familyData"`
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
	UserData *goparent.User `json:"userData"`
	Token    string         `json:"token"`
}

//InvitesResponse - response structure for invites
type InvitesResponse struct {
	SentInviteData    []*goparent.UserInvitation `json:"sentInviteData"`
	PendingInviteData []*goparent.UserInvitation `json:"pendingInviteData"`
}

func (h *Handler) initUsersHandlers(r *mux.Router) {
	u := r.PathPrefix("/user").Subrouter()
	u.Handle("/", h.userNewHandler()).Methods("POST").Name("UserNew")
	u.Handle("/", h.AuthRequired(h.userGetHandler())).Methods("GET").Name("UserGetData")
	u.Handle("/{id}", h.AuthRequired(h.userGetHandler())).Methods("GET").Name("UserView")
	u.Handle("/login", h.loginHandler()).Methods("POST").Name("UserLogin")
	u.Handle("/invite", h.AuthRequired(h.userListInviteHandler())).Methods("GET").Name("UserGetSentInvites")
	u.Handle("/invite", h.AuthRequired(h.userNewInviteHandler())).Methods("POST").Name("UserNewInvite")
	u.Handle("/invite/{id}", h.AuthRequired(h.userDeleteInviteHandler())).Methods("DELETE").Name("UserDeleteInvite")
	u.Handle("/invite/accept/{id}", h.AuthRequired(h.userAcceptInviteHandler())).Methods("POST").Name("UserAcceptInvite")
}

func (h *Handler) loginHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		password := r.FormValue("password")
		ctx := h.Env.DB.GetContext(r)

		user, err := h.UserService.UserByLogin(ctx, username, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		token, err := h.UserService.GetToken(user)
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

func (h *Handler) userGetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := h.Env.DB.GetContext(r)
		user, err := UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		family, err := h.UserService.GetFamily(ctx, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		userInfo := UserResponse{
			UserData:   user,
			FamilyData: family}
		w.Header().Set("Content-Type", jsonContentType)
		json.NewEncoder(w).Encode(userInfo)
		return
	})
}

func (h *Handler) userNewHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := h.Env.DB.GetContext(r)
		decoder := json.NewDecoder(r.Body)
		var newUserRequest NewUserRequest
		err := decoder.Decode(&newUserRequest)
		userData := goparent.User(newUserRequest.UserData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		w.Header().Set("Content-Type", jsonContentType)
		err = h.UserService.Save(ctx, &userData)
		if err != nil {
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

func (h *Handler) userNewInviteHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := h.Env.DB.GetContext(r)
		user, err := UserFromContext(ctx)
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

		err = h.UserInvitationService.InviteParent(ctx, user, invitedUserEmail, time.Now())
		if err != nil {
			if err == goparent.ErrExistingInvitation {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})
}

func (h *Handler) userListInviteHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := h.Env.DB.GetContext(r)
		user, err := UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		sentInvites, err := h.UserInvitationService.SentInvites(ctx, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pendingInvites, err := h.UserInvitationService.Invites(ctx, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		invitesResponse := InvitesResponse{SentInviteData: sentInvites, PendingInviteData: pendingInvites}
		w.Header().Set("Content-Type", jsonContentType)
		json.NewEncoder(w).Encode(invitesResponse)
	})
}

func (h *Handler) userAcceptInviteHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := h.Env.DB.GetContext(r)
		user, err := UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		id := mux.Vars(r)["id"]
		err = h.UserInvitationService.Accept(ctx, user, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		w.Header().Set("Content-Type", jsonContentType)
		return
	})
}

func (h *Handler) userDeleteInviteHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := h.Env.DB.GetContext(r)
		_, err := UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		id := mux.Vars(r)["id"]
		invite, err := h.UserInvitationService.Invite(ctx, id)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		err = h.UserInvitationService.Delete(ctx, invite)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", jsonContentType)
		w.WriteHeader(http.StatusNoContent)
	})
}
