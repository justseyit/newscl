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
	limit      = 500
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
		if len(newsListToSend.News) > limit {
			log.Printf("Sending %d news to the API", limit)
			newsLists := splitList(newsListToSend)
			for _, newsList := range newsLists {
				err := postToApi(newsList)
				if err != nil {
					return err
				}
			}
		} else {
			log.Printf("Sending %d news to the API", len(newsListToSend.News))
			err := postToApi(newsListToSend)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func splitList(newsList model.NewsClNewsList) []model.NewsClNewsList {
	var newsLists []model.NewsClNewsList

	for i := 0; i < len(newsList.News); i += limit {
		end := i + limit

		if end > len(newsList.News) {
			end = len(newsList.News)
		}

		newsLists = append(newsLists, model.NewsClNewsList{News: newsList.News[i:end]})
	}

	return newsLists
}

func postToApi(newsList model.NewsClNewsList) error {
	jsonBody, err := json.Marshal(newsList)
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

	return nil
}
