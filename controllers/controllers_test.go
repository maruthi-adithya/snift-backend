package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"runtime"
	"snift-backend/models"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestHomePage(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HomePage)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, rr.Body.String(), "Welcome to Snift!")

}
func TestInvalidURL(t *testing.T) {

	tokenreq, _ := http.NewRequest("GET", "/token", nil)
	tokenrr := httptest.NewRecorder()
	tokenhandler := http.HandlerFunc(GetAuthToken)
	tokenhandler.ServeHTTP(tokenrr, tokenreq)

	assert.Equal(t, tokenrr.Code, http.StatusOK)

	var token models.Token
	json.NewDecoder(tokenrr.Body).Decode(&token)

	var urlJSON = `{"url":"example"}`

	req, _ := http.NewRequest("POST", "/scores", strings.NewReader(urlJSON))
	req.Header.Set("X-Auth-Token", token.Token)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetScore)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusBadRequest)
	assert.Equal(t, rr.Body.String(), "{\"error\":\"Invalid URL\"}")

}
func TestValidURL(t *testing.T) {

	tokenreq, _ := http.NewRequest("GET", "/token", nil)
	tokenrr := httptest.NewRecorder()
	tokenhandler := http.HandlerFunc(GetAuthToken)
	tokenhandler.ServeHTTP(tokenrr, tokenreq)

	assert.Equal(t, tokenrr.Code, http.StatusOK)

	var token models.Token
	json.NewDecoder(tokenrr.Body).Decode(&token)

	var urlJSON = `{"url":"https://www.example.com"}`

	req, _ := http.NewRequest("POST", "/scores", strings.NewReader(urlJSON))
	req.Header.Set("X-Auth-Token", token.Token)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetScore)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)

}

func TestUnauthenticatedRequest(t *testing.T) {
	var urlJSON = `{"url":"https://www.example.com"}`
	req, _ := http.NewRequest("POST", "/scores", strings.NewReader(urlJSON))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetScore)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusUnauthorized)
	assert.Equal(t, rr.Body.String(), "{\"error\":\"Invalid Token\"}")
}
