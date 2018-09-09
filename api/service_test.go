package api

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sasimpson/goparent"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func getTestHandler() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "thanks")
		return
		// panic("test entered test handler, this should not happen")
	}
	return http.HandlerFunc(fn)
}

//sample/test token creation.
func makeTestToken(user *goparent.User, key interface{}) string {
	// foo := map[string]interface{}{
	// 	"Name":     user.Name,
	// 	"ID":       user.ID,
	// 	"Email":    user.Email,
	// 	"Username": user.Username,
	// 	"exp":      time.Now().Add(time.Hour).Unix(),
	// }
	// claims := &goparent.UserClaims{
	// 	user.ID,
	// 	user.Name,
	// 	user.Email,
	// 	user.Username,
	// 	user.Password,
	// 	jwt.StandardClaims{
	// 		ExpiresAt: time.Now().Add(time.Hour).Unix(),
	// 	},
	// }
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	log.Printf("%#v", token)
	log.Printf("%#v", claims)
	log.Printf("%#v", user)
	claims["Name"] = user.Name
	claims["ID"] = user.ID
	claims["Email"] = user.Email
	claims["Username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	tokenString, err := token.SignedString(key)

	if err != nil {
		panic(err.Error())
	}

	return tokenString
}
func TestBaseHandlers(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(apiHandler)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	req, err = http.NewRequest("GET", "/info", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler = http.HandlerFunc(infoHandler)
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAuthRequiredMiddleware(t *testing.T) {
	t.Skip()
	testCases := []struct {
		desc         string
		env          *goparent.Env
		user         *goparent.User
		responseCode int
	}{
		// {
		// 	desc: "no auth",
		// 	env: &goparent.Env{
		// 		Auth: config.Authentication{
		// 			SigningKey: []byte("testkey"),
		// 		},
		// 	},
		// 	user: &goparent.User{
		// 		ID:       "1",
		// 		Name:     "Test User",
		// 		Email:    "testuser@test.com",
		// 		Username: "testuser",
		// 		Password: "",
		// 	},
		// 	responseCode: http.StatusUnauthorized,
		// },
		{
			desc: "with token",
			env: &goparent.Env{
				Auth: goparent.Authentication{
					SigningKey: []byte("testkey"),
				},
			},
			user: &goparent.User{
				ID:       "1",
				Name:     "Test User",
				Email:    "testuser@test.com",
				Username: "testuser",
				Password: "",
			},
			responseCode: http.StatusUnauthorized,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockHandler := Handler{}

			token := makeTestToken(tC.user, tC.env.Auth.SigningKey)
			req, err := http.NewRequest("GET", "/test", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
			rr := httptest.NewRecorder()
			handler := mockHandler.AuthRequired(getTestHandler())
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tC.responseCode, rr.Code)
		})
	}
}
