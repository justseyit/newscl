package repository

import (
	"newscl/model"
	"newscl/service"
)

func GetNewsByProvider(provider model.Provider) (model.NewsClNewsList, error) {
	var newsList model.NewsClNewsList
	switch provider {
	case model.BBC:
		bbcNews, err := service.GetBBCNewsRSS()
		if err != nil {
			return newsList, err
		}
		newsList = bbcNews.BBCRSSToNewsClNewsList()
	case model.REUTERS:
		reutersNews, err := service.GetReutersNewsRSS()
		if err != nil {
			return newsList, err
		}
		newsList = reutersNews.ReutersRSSToNewsClNewsList()
	}
	return newsList, nil
}