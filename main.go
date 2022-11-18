package main

import (
	"log"
	"newscl/api"
	"newscl/model"
	"newscl/repository"
)

/*
const (
	port = 9999
)*/

func main() {

	repository.InitMongoDB()

	//mux := mux.NewRouter()

	//log.Println("Starting the scheduled task")
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

	/*executor := sch.NewTimedExecutor(time.Second*3, time.Minute)

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
	}, true)*/
}
