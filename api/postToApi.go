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
	"time"
)

const (
	serviceURL = "https://source.trademarketdata.com/api/paper-news/store-bulk"
)

// PostToApi sends the news to the API
func PostToApi(newsList model.NewsClNewsList) error {

	var newsListToSend model.NewsClNewsList
	newsListToSend.News = make([]model.NewsClNews, 0)

	for _, news := range newsList.News {
		succeeded, _ := repository.Mongo.IsSucceeded(news.ID)
		if !succeeded {
			newsListToSend.News = append(newsListToSend.News, news)
		} else {
			log.Printf("%s has already been sent to the API", news.ID)
		}
	}

	if len(newsListToSend.News) == 0 {
		log.Println("No news to send to the API.")
	} else {
		jsonBody, err := json.Marshal(newsListToSend)
		if err != nil {
			return err
		}
		bodyReader := bytes.NewReader(jsonBody)

		req, err := http.NewRequest(http.MethodPost, serviceURL, bodyReader)

		if err != nil {
			return err
			//fmt.Printf("client: could not create request: %s\n", err)
			//os.Exit(1)
		}
		req.Header.Set("Content-Type", "application/json")

		client := http.Client{
			Timeout: 30 * time.Second,
		}

		res, err1 := client.Do(req)
		if err1 != nil {
			return err1
			//fmt.Printf("client: error making http request: %s\n", err)
			//os.Exit(1)
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(res.Body)
			fmt.Println(string(bodyBytes))
			return fmt.Errorf("error sending news to the API. Status code: %d", res.StatusCode)
		}

		for _, news := range newsList.News {
			err = repository.Mongo.FlagSucceeded(news)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
