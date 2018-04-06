package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"encoding/json"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/config"
	"github.com/sasimpson/goparent/rethinkdb"
)

type contextKey string

func (c contextKey) String() string {
	return "api context key " + string(c)
}

const (
	jsonContentType string     = "application/json"
	userContextKey  contextKey = "user"
)

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
	Env                   *config.Env
}

//ServiceInfo - return data about the service
type ServiceInfo struct {
	Version  string `json:"version"`
	Hostname string `json:"hostname"`
}

//ErrService - error message format for service calls
type ErrService struct {
	ErrMessage struct {
		Body string `json:"body"`
		Code int    `json:"code"`
	} `json:"error"`
}

//RunService - Runs service interfaces for app
func RunService(env *config.Env) {
	log.SetOutput(os.Stdout)

	serviceHandler := Handler{
		UserService:           &rethinkdb.UserService{Env: env},
		UserInvitationService: &rethinkdb.UserInviteService{Env: env},
		FamilyService:         &rethinkdb.FamilyService{Env: env},
		ChildService:          &rethinkdb.ChildService{Env: env},
		FeedingService:        &rethinkdb.FeedingService{Env: env},
		SleepService:          &rethinkdb.SleepService{Env: env},
		WasteService:          &rethinkdb.WasteService{Env: env},
		Env:                   env,
	}

	r := mux.NewRouter()
	a := r.PathPrefix("/api").Subrouter()
	a.HandleFunc("/", apiHandler)
	a.HandleFunc("/info", infoHandler)

	serviceHandler.initUsersHandlers(a)
	serviceHandler.initChildrenHandlers(a)
	serviceHandler.initFeedingHandlers(a)
	serviceHandler.initSleepHandlers(a)
	serviceHandler.initWasteHandlers(a)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Accept", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	log.Printf("starting service on 8000")
	http.Handle("/", r)
	http.ListenAndServe(":8000", handlers.CORS(originsOk, headersOk, methodsOk)(r))
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
		token, err := request.ParseFromRequestWithClaims(r, request.AuthorizationHeaderExtractor, &goparent.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return sh.Env.Auth.SigningKey, nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(*goparent.UserClaims); ok && token.Valid {
			user, err := sh.UserService.User(claims.ID)
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
