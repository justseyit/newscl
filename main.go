package main

import (
	"fmt"
	"log"
	"net/http"
	"newscl/api"
	"newscl/repository"

	"github.com/gorilla/mux"
)

func main() {

	repository.InitMongoDB()

	mux := mux.NewRouter()

	mux.HandleFunc("/news/{provider}", api.GetNewsHandler).Methods("GET")

	fmt.Println("Server running on port 9999")

	log.Fatal(http.ListenAndServe(":9999", mux))
}