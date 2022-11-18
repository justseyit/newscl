package main

import (
	"fmt"
	"log"
	"newscl/api"
	"newscl/model"
	"newscl/repository"
	"time"

	sch "github.com/gbenroscience/scheduled-executor/utils"
)

func main() {

	repository.InitMongoDB()

	log.Println("Starting the scheduled task")

	executor := sch.NewTimedExecutor(time.Second*3, time.Hour)

	executor.Start(func() {
		bbcNews, _ := repository.GetNewsByProvider(model.BBC)
		errBBC := api.PostToApi(bbcNews, model.BBC)
		if errBBC != nil {
			log.Fatalf("Error posting BBC news to API: %v", errBBC)
		}

	}, true)

	executor.Start(func() {
		reutersNews, _ := repository.GetNewsByProvider(model.REUTERS)
		errReuters := api.PostToApi(reutersNews, model.REUTERS)
		if errReuters != nil {
			log.Fatalf("Error posting Reuters news to API: %v", errReuters)
		}
	}, true)

	fmt.Scanln()
}
