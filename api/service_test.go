package api

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestBaseHandlers(t *testing.T) {
// 	req, err := http.NewRequest("GET", "/", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	handler := http.HandlerFunc(apiHandler)
// 	rr := httptest.NewRecorder()
// 	handler.ServeHTTP(rr, req)
// 	assert.Equal(t, http.StatusOK, rr.Code)

// 	req, err = http.NewRequest("GET", "/info", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	handler = http.HandlerFunc(infoHandler)
// 	handler.ServeHTTP(rr, req)
// 	assert.Equal(t, http.StatusOK, rr.Code)
// }

// func TestAuthRequiredHandler(t *testing.T) {

// }
