package service

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"newscl/model"
)

func GetBBCNews() (model.NewsClNewsList, error) {
	var newsList model.NewsClNewsList
	bbcNews, err := GetBBCNewsRSS()
	if err != nil {
		return newsList, err
	}
	newsList = bbcNews.BBCRSSToNewsClNewsList()
	return newsList, nil
}

func GetBBCNewsRSS() (model.BBCRSS, error) {
	var bbcNews model.BBCRSS
	err := GetRSS("https://feeds.bbci.co.uk/news/world/rss.xml", &bbcNews)
	if err != nil {
		return model.BBCRSS{}, err
	}
	return bbcNews, nil
}

func GetRSS(url string, v interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = xml.Unmarshal(body, v)
	if err != nil {
		return err
	}
	return nil
}

func GetReutersNews() (model.NewsClNewsList, error) {
	var newsList model.NewsClNewsList
	reutersNews, err := GetReutersNewsRSS()
	if err != nil {
		return newsList, err
	}
	newsList = reutersNews.ReutersRSSToNewsClNewsList()
	return newsList, nil
}

func GetReutersNewsRSS() (model.ReutersRSS, error) {
	var reutersNews model.ReutersRSS
	err := GetRSS("https://www.reutersagency.com/feed", &reutersNews)
	if err != nil {
		return model.ReutersRSS{}, err
	}
	return reutersNews, nil
}