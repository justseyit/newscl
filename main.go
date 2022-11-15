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

const (
	port = 9999
)

func main() {

	repository.InitMongoDB()
	model.InitServiceInfo()

	mux := mux.NewRouter()

	log.Println("Starting the scheduled task")

	executor := sch.NewTimedExecutor(time.Second*3, time.Minute)

	executor.Start(func() {
		bbcNews, _ := repository.GetNewsByProvider(model.BBC)
		errBBC := api.PostToApi(bbcNews)
		if errBBC != nil {
			log.Fatalf("Error posting BBC news to API: %v", errBBC)
		}
		
	}, true)

	executor.Start(func() {
		reutersNews, _ := repository.GetNewsByProvider(model.REUTERS)
		errReuters := api.PostToApi(reutersNews)
		if errReuters != nil {
			log.Fatalf("Error posting Reuters news to API: %v", errReuters)
		}
	}, true)


	mux.HandleFunc("/info", api.GetInfoHandler).Methods("GET")

	fmt.Printf("Server running on port %d\n", port)

	log.Fatal(http.ListenAndServe(":"+fmt.Sprint(port), mux))
}
