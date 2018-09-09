package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"encoding/json"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gorilla/mux"
	"github.com/sasimpson/goparent"
)

type contextKey string

func (c contextKey) String() string {
	return "api context key " + string(c)
}

const (
	jsonContentType string     = "application/json"
	userContextKey  contextKey = "user"
)

type ServiceHandler interface {
	GetContext(*http.Request) *context.Context
}

//Handler - this is the handler struct that contains all of the interfaces for
// the api.  the implementation can be changed by inserting different implementations
// of the interface
type Handler struct {
	UserService           goparent.UserService
	UserInvitationService goparent.UserInvitationService
	FamilyService         goparent.FamilyService
	ChildService          goparent.ChildService
	FeedingService        goparent.FeedingService
	SleepService          goparent.SleepService
	WasteService          goparent.WasteService
	Env                   *goparent.Env
}

//ServiceInfo - return data about the service
type ServiceInfo struct {
	Version  string `json:"version"`
	Hostname string `json:"hostname"`
}

//Pagination - structure to hold the pagination data for service calls and responses
type Pagination struct {
	Skip  uint64
	Take  uint64
	Total uint64
	Days  uint64
}

//ErrService - error message format for service calls
type ErrService struct {
	ErrMessage struct {
		Body string `json:"body"`
		Code int    `json:"code"`
	} `json:"error"`
}

//BuildAPIRouting - common api routing here if passed a handler
func BuildAPIRouting(serviceHandler *Handler) *mux.Router {
	r := mux.NewRouter()
	a := r.PathPrefix("/api").Subrouter()
	a.HandleFunc("/", apiHandler)
	a.HandleFunc("/info", infoHandler)

	serviceHandler.initUsersHandlers(a)
	serviceHandler.initChildrenHandlers(a)
	serviceHandler.initFeedingHandlers(a)
	serviceHandler.initSleepHandlers(a)
	serviceHandler.initWasteHandlers(a)
	return r
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "please check docs at: https://github.com/sasimpson/goparent")
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	name, _ := os.Hostname()

	si := ServiceInfo{Version: "v0.1", Hostname: name}
	json.NewEncoder(w).Encode(si)
	return
}

//AuthRequired - handler to handle authentication of users tokens.
func (sh *Handler) AuthRequired(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := sh.Env.DB.GetContext(r)
		token, err := request.ParseFromRequestWithClaims(r, request.AuthorizationHeaderExtractor, &goparent.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return sh.Env.Auth.SigningKey, nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(*goparent.UserClaims); ok && token.Valid {
			user, err := sh.UserService.User(ctx, claims.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), userContextKey, user)
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
			return
		}
		http.Error(w, "failed", http.StatusInternalServerError)
		return
	})
}

//UserFromContext - helper to get the user from the request context
func UserFromContext(ctx context.Context) (*goparent.User, error) {
	user, ok := ctx.Value(userContextKey).(*goparent.User)
	if !ok {
		return nil, errors.New("no user found in context")
	}
	return user, nil
}

func getPagination(r *http.Request) *Pagination {
	q := r.URL.Query()

	days, err := strconv.ParseUint(q.Get("days"), 10, 64)
	if err != nil {
		days = 7
	}
	return &Pagination{Days: days}
}
