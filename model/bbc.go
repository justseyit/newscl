package model

import "encoding/xml"

type BBCRSS struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Dc      string   `xml:"dc,attr"`
	Content string   `xml:"content,attr"`
	Atom    string   `xml:"atom,attr"`
	Version string   `xml:"version,attr"`
	Media   string   `xml:"media,attr"`
	Channel struct {
		Text        string `xml:",chardata"`
		Title       string `xml:"title"`
		Description string `xml:"description"`
		Link        string `xml:"link"`
		Image       struct {
			Text  string `xml:",chardata"`
			URL   string `xml:"url"`
			Title string `xml:"title"`
			Link  string `xml:"link"`
		} `xml:"image"`
		Generator     string `xml:"generator"`
		LastBuildDate string `xml:"lastBuildDate"`
		Copyright     string `xml:"copyright"`
		Language      string `xml:"language"`
		Ttl           string `xml:"ttl"`
		Item          []struct {
			Text        string `xml:",chardata"`
			Title       string `xml:"title"`
			Description string `xml:"description"`
			Link        string `xml:"link"`
			Guid        struct {
				Text        string `xml:",chardata"`
				IsPermaLink string `xml:"isPermaLink,attr"`
			} `xml:"guid"`
			PubDate string `xml:"pubDate"`
		} `xml:"item"`
	} `xml:"channel"`
}

func (b BBCRSS) BBCRSSToNewsClNewsList() NewsClNewsList {
	var newsList NewsClNewsList
	for _, item := range b.Channel.Item {
		newsList.News = append(newsList.News, newsClNews{
			Title:       item.Title,
			Description: item.Description,
			URL:         item.Link,
			Time:        item.PubDate,
			Location:    "UK",
			Language:    "EN",
			Type:        2,
			SourceName:  "BBC",
			SourceURL:   "https://www.bbc.com/news",
			ImageURL:    "",
		})
	}
	return newsList
}
