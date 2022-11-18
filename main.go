package main

import (
	"log"
	"newscl/api"
	"newscl/model"
	"newscl/repository"
)

const (
	port = 9999
)

func main() {

	repository.InitMongoDB()

	log.Println("Starting the task")

	bbcNews, _ := repository.GetNewsByProvider(model.BBC)
	errBBC := api.PostToApi(bbcNews, model.BBC)
	if errBBC != nil {
		log.Fatalf("Error posting BBC news to API: %v", errBBC)
	}

	reutersNews, _ := repository.GetNewsByProvider(model.REUTERS)
	errReuters := api.PostToApi(reutersNews, model.REUTERS)
	if errReuters != nil {
		log.Fatalf("Error posting Reuters news to API: %v", errReuters)
	}
}
