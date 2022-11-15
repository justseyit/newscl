package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"newscl/model"
	"newscl/repository"
	"os"
	"time"
)

const (
	serviceURL = "https://source.trademarketdata.com/api/paper-news/store-bulk"
)

// PostToApi sends the news to the API
func PostToApi(newsList model.NewsClNewsList) error {

	var newsListToSend model.NewsClNewsList

	for _, news := range newsList.News {
		succeeded, _ := repository.Mongo.IsSucceeded(news.ID)
		if !succeeded {
			newsListToSend.News = append(newsListToSend.News, news)
		}else{
			log.Printf("%s has already been sent to the API", news.ID)
		}
	}

	jsonBody, err := json.Marshal(newsList)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(http.MethodPost, serviceURL, bodyReader)

	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	defer res.Body.Close()


	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		fmt.Println(string(bodyBytes))
		return fmt.Errorf("error sending news to API. Status code: %d", res.StatusCode)
	}

	for _, news := range newsList.News {
		err = repository.Mongo.FlagSucceeded(news)
		if err != nil {
			return err
		}
	}

	return nil
}
