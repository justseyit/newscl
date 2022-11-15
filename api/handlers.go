package api

import (
	"encoding/json"
	"net/http"
	"newscl/model"
	"newscl/repository"

	"github.com/gorilla/mux"
)

func sendResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, err error, statusCode int) {
	sendResponse(w, map[string]string{"error": err.Error()}, statusCode)
}

func GetNewsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := model.ToProvider(vars["provider"])
	newsList, err := repository.GetNewsByProvider(provider)
	if err != nil {
		sendError(w, err, http.StatusInternalServerError)
		return
	}
	sendResponse(w, newsList, http.StatusOK)

	///TODO: Add a database to store the news
}

func GetInfoHandler(w http.ResponseWriter, r *http.Request) {
	
}
