package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"encoding/json"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sasimpson/goparent"
	"github.com/sasimpson/goparent/config"
	"github.com/sasimpson/goparent/datastore"
	"github.com/sasimpson/goparent/rethinkdb"
	"google.golang.org/appengine"
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

	r := buildAPIRouting()

	log.Printf("starting service on 8000")
	http.Handle("/", r)
	http.ListenAndServe(":8000", handlers.CORS(originsOk, headersOk, methodsOk)(r))
}

//RunAppEngineService - runs service in appengine
func RunAppEngineService(env *config.Env) {
	serviceHandler := Handler{
		Env:                   env,
		UserService:           &datastore.UserService{Env: env},
		UserInvitationService: &datastore.UserInviteService{Env: env},
		FamilyService:         &datastore.FamilyService{Env: env},
		ChildService:          &datastore.ChildService{Env: env},
		FeedingService:        &datastore.FeedingService{Env: env},
		SleepService:          &datastore.SleepService{Env: env},
		WasteService:          &datastore.WasteService{Env: env},
	}

	r := buildAPIRouting()
	log.Printf("starting appengine service...")
	http.Handle("/", r)
	appengine.Main()
}

//buildAPIRouting - common api routing here
func buildAPIRouting() *mux.Router {
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

func getPagination(r *http.Request) *Pagination {
	q := r.URL.Query()

	days, err := strconv.ParseUint(q.Get("days"), 10, 64)
	if err != nil {
		days = 7
	}
	return &Pagination{Days: days}
}
