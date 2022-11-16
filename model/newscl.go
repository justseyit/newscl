package model

import "net/http"

type NewsClNewsList struct {
	News []NewsClNews `json:"news" xml:"news"`
}

type NewsClNews struct {
	ID          string `json:"guid" xml:"guid" bson:"_id"`
	Title       string `json:"title" xml:"title"`
	URL         string `json:"url" xml:"url"`
	Description string `json:"description" xml:"description"`
	SourceName  string `json:"source_name" xml:"source_name"`
	SourceURL   string `json:"source_url" xml:"source_url"`
	ImageURL    string `json:"image_url" xml:"image_url"`
	Language    string `json:"language" xml:"language"`
	Location    string `json:"location" xml:"location"`
	Time        int    `json:"time" xml:"time"`
	Tags        string `json:"tags" xml:"tags"`
	Type        int    `json:"type" xml:"type"`
}

type NewsClResponse struct {
	http.Response
	Status       string         `json:"status" xml:"status"`
	TotalResults int            `json:"totalResults" xml:"totalResults"`
	NewsList     NewsClNewsList `json:"newsList" xml:"newsList"`
}
