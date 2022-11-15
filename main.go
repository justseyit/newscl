package main

import (
	"fmt"
	"log"
	"net/http"
	"newscl/api"
	"newscl/model"
	"newscl/repository"
	"time"

	sch "github.com/gbenroscience/scheduled-executor/utils"
	"github.com/gorilla/mux"
)

const(
	port = 9999
)

func main() {

	repository.InitMongoDB()
	model.InitServiceInfo()

	mux := mux.NewRouter()

	log.Println("Starting the scheduled task")

	executor := sch.NewTimedExecutor(time.Second * 3, time.Minute)

	bbcNews, _ := repository.GetNewsByProvider(model.BBC)
	reutersNews, _ := repository.GetNewsByProvider(model.REUTERS)

	executor.Start(func() {
		err := api.PostToApi(bbcNews)
		if err != nil {
			log.Fatalf("Error posting news to API: %v", err)
		}
	}, true)

	executor.Start(func() {
		err := api.PostToApi(reutersNews)
		if err != nil {
			log.Fatalf("Error posting news to API: %v", err)
		}
	}, true)

	mux.HandleFunc("/info", api.GetInfoHandler).Methods("GET")

	fmt.Printf("Server running on port %d\n", port)

	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(port), mux))
}