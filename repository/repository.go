package repository

import (
	"newscl/model"
	"newscl/service"
)

var (
	Mongo          service.MongoDB
	NewsClNewsList model.NewsClNewsList
	ServiceInfo    model.ServiceInfo
)

func Init() {
	Mongo = *service.NewMongoDB()
	ServiceInfo = model.ServiceInfo{
		ServiceState:         "waiting",
		SourceCheckFrequency: 0,
		SourceInfo: make([]struct {
			SourceID             string "json:\"sourceID\" bson:\"sourceID\""
			SourceName           string "json:\"sourceName\" bson:\"sourceName\""
			SourceURL            string "json:\"sourceURL\" bson:\"sourceURL\""
			SourceFetchFrequency uint   "json:\"sourceFetchFrequency\" bson:\"sourceFetchFrequency\""
			LastChecked          int    "json:\"lastChecked\" bson:\"lastChecked\""
		}, 0),
		RunningTime: 0,
	}
}

func GetNewsByProviderLink(link string) error {
	news, err := service.GetRSS(link)
	if err != nil {
		return err
	}
	NewsClNewsList.News = append(NewsClNewsList.News, news.RSSToNewsClNewsList().News...)
	Mongo.SendNews(NewsClNewsList)
	return nil
}
