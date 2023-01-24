package service

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"newscl/model"
)

func GetRSS(url string) (model.RSS, error) {
	var rss model.RSS
	resp, err := http.Get(url)
	if err != nil {
		return rss, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return rss, err
	}
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		return rss, err
	}
	return rss, nil
}

func GetNews(link string) (model.NewsClNewsList, error) {
	var newsList model.NewsClNewsList
	rss, err := GetRSS(link)
	if err != nil {
		return newsList, err
	}
	newsList = rss.RSSToNewsClNewsList()
	return newsList, nil
}
