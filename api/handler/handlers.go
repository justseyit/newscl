package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"newscl/api"
	"newscl/model"
	"newscl/repository"
	"newscl/service"
	"sync"
	"time"
	"html/template"

	"github.com/gbenroscience/scheduled-executor/utils"
)

var executors = make(map[string]*utils.ScheduledExecutor)

// HTTP Handler for getting source list
func GetSourceList(w http.ResponseWriter, r *http.Request) {
	sources, err := repository.Mongo.GetSourcesFromDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sources)
}

// HTTP Handler for deleting source
func DeleteSource(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		var deleteSourceRequest struct {
			IDList []string `json:"idList"`
		}
		err := json.NewDecoder(r.Body).Decode(&deleteSourceRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(deleteSourceRequest.IDList) == 0 {
			http.Error(w, "idList is empty", http.StatusBadRequest)
			return
		}
		err = repository.Mongo.RemoveSourcesFromDB(deleteSourceRequest.IDList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HTTP Handler for adding source
func PostSource(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var sources []model.Source
	err := json.NewDecoder(r.Body).Decode(&sources)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(sources) == 0 {
		http.Error(w, "sources is empty", http.StatusBadRequest)
		return
	}
	err = repository.Mongo.SaveSourcesToDB(sources)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HTTP Handler for starting the service
func StartService(w http.ResponseWriter, r *http.Request) {
	if repository.ServiceInfo.ServiceState != "running" {
		var parameters struct {
			Interval uint `json:"interval"`
		}
		err := json.NewDecoder(r.Body).Decode(&parameters)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if parameters.Interval == 0 {
			http.Error(w, "interval is empty", http.StatusBadRequest)
			return
		}
		repository.ServiceInfo.SourceCheckFrequency = parameters.Interval
		repository.ServiceInfo.ServiceState = "running"
		repository.ServiceInfo.RunningTime = 0
		startTicker()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Service has been started."))
	} else {
		http.Error(w, "Service is already running.", http.StatusBadRequest)
	}
}

func startTicker() {

	executorServiceRunningTimer := utils.NewTimedExecutor(1*time.Second, 1*time.Second)
	executorServiceRunningTimer.Start(func() {
		repository.ServiceInfo.RunningTime++
	}, true)
	executors["executorServiceRunningTimer"] = &executorServiceRunningTimer

	executorSourceCheckTimer := utils.NewTimedExecutor(time.Duration(repository.ServiceInfo.SourceCheckFrequency*uint(time.Second)), time.Duration(repository.ServiceInfo.SourceCheckFrequency*uint(time.Second)))
	executorSourceCheckTimer.Start(func() {
		_, err := repository.Mongo.GetSourcesFromDB()
		if err != nil {
			log.Println(err)
		}
	}, true)
	executors["executorSourceCheckTimer"] = &executorSourceCheckTimer

	sources, err := repository.Mongo.GetSourcesFromDB()
	if err != nil {
		log.Println(err)
	}

	mutex := sync.Mutex{}
	for _, source := range sources {
		executorSourceFetchTimer := utils.NewTimedExecutor(time.Duration(source.FetchFrequency*uint(time.Second)), time.Duration(source.FetchFrequency*uint(time.Second)))
		executorSourceFetchTimer.Start(func() {
			mutex.Lock()
			rss, err := service.GetRSS(source.URL)
			if err != nil {
				log.Println(err)
			}
			nl := rss.RSSToNewsClNewsList()
			repository.NewsClNewsList.News = append(repository.NewsClNewsList.News, nl.News...)
			api.PostToApi()
			err = repository.Mongo.SendNews(nl)
			if err != nil {
				log.Println(err)
			}
			repository.NewsClNewsList.News = []model.NewsClNews{}
			mutex.Unlock()
		}, true)
		executors[source.ID] = &executorSourceFetchTimer
	}
}

// HTTP Handler for getting service info
func GetServiceInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(repository.ServiceInfo)
}

// HTTP Handler for stopping the service
func StopService(w http.ResponseWriter, r *http.Request) {
	for name, executor := range executors {
		log.Println("Stopping executor: " + name)
		executor.Close()
	}
	repository.ServiceInfo.ServiceState = "stopped"
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Service has been stopped."))
}

func LanguagesHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/html/addLanguage.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, nil)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
