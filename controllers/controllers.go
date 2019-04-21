package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"snift-backend/models"
	"snift-backend/services"
	"snift-backend/utils"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// HandleRequests - Handler for all API Requests
func HandleRequests() {
	port := os.Getenv("PORT")
	log.Print("Server starting at PORT ", port)
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", HomePage).Methods("GET")
	myRouter.HandleFunc("/scores", GetScore).Methods("POST")
	myRouter.HandleFunc("/token", GetAuthToken).Methods("GET")
	log.Fatal(http.ListenAndServe(port, myRouter))
}

// HomePage - the default root endpoint of Snift Backend
func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /")
	fmt.Fprintf(w, "Welcome to Snift!")
}

// GetScore - POST /scores handler
func GetScore(w http.ResponseWriter, r *http.Request) {
	if !utils.ValidateToken(r) {
		utils.Unauthorized(w, true, "Invalid Token")
		return
	}
	start := time.Now()
	var scoresRequest models.ScoresRequest
	err := json.NewDecoder(r.Body).Decode(&scoresRequest)
	if err != nil {
		fmt.Println(err)
		utils.BadRequest(w, true, "Unexpected Error Occured")
		return
	}
	log.Print("POST /scores")

	err = utils.IsValidURL(scoresRequest.URL)
	if err != nil {
		utils.BadRequest(w, true, "Invalid URL")
		return
	}
	response, scoresError := services.CalculateOverallScore(scoresRequest.URL)
	if scoresError != nil {
		if strings.Contains(scoresError.Error(), "no such host") {
			utils.BadRequest(w, true, "Invalid Domain")
			return
		}
		utils.InternalServerError(w, true, "Unexpected Error Occured")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Printf("Score for %s obtained in %v seconds \n", scoresRequest.URL, time.Since(start).Seconds())
	utils.Writer(w.Write(response))
}

// GetAuthToken - GET /scores handler
func GetAuthToken(w http.ResponseWriter, r *http.Request) {
	response, err := utils.GetToken(r)
	if err != nil {
		log.Println("Unexpected Error Occured", err)
	}

	responseBody, jsonError := json.Marshal(response)
	if jsonError != nil {
		utils.InternalServerError(w, true, "Unexpected Error Occured")
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	utils.Writer(w.Write([]byte(responseBody)))
}
