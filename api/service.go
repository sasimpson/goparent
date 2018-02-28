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
	"github.com/sasimpson/goparent/config"
	"github.com/sasimpson/goparent/models"
)

type contextKey string

func (c contextKey) String() string {
	return "api context key " + string(c)
}

const (
	jsonContentType string     = "application/json"
	userContextKey  contextKey = "user"
)

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
	r := mux.NewRouter()
	a := r.PathPrefix("/api").Subrouter()
	a.HandleFunc("/", apiHandler)
	a.HandleFunc("/info", infoHandler)

	initUsersHandlers(env, a)
	initFeedingHandlers(env, a)
	initSleepHandlers(env, a)
	initWasteHandlers(env, a)
	initChildrenHandlers(env, a)

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
func AuthRequired(h http.Handler, env *config.Env) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := request.ParseFromRequestWithClaims(r, request.AuthorizationHeaderExtractor, &models.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return env.Auth.SigningKey, nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		var user models.User
		if claims, ok := token.Claims.(*models.UserClaims); ok && token.Valid {
			user.GetUser(env, claims.ID)
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
func UserFromContext(ctx context.Context) (models.User, error) {
	user, ok := ctx.Value(userContextKey).(models.User)
	if !ok {
		return models.User{}, errors.New("no user found in context")
	}
	return user, nil
}
