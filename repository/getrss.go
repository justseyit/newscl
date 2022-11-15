package repository

import (
	"newscl/model"
	"newscl/service"
)

var Mongo service.MongoDB

func InitMongoDB() {
	Mongo = *service.NewMongoDB()
}

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
	Mongo.SendNews(newsList)
	return newsList, nil
}

func PostNews(newsList model.NewsClNewsList) error {
	return Mongo.SendNews(newsList)
}