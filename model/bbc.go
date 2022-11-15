package model

import (
	"encoding/xml"
	"strings"
	"time"
)

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
		newsList.News = append(newsList.News, NewsClNews{
			ID:          item.Guid.Text,
			Title:       item.Title,
			Description: item.Description,
			URL:         item.Link,
			Time:        b.toTimestamp(item.PubDate),
			Location:    "UK",
			Language:    "EN",
			Type:        2,
			SourceName:  "BBC",
			SourceURL:   "https://www.bbc.com/news",
			Tags:        "",
			ImageURL:    "https://www.bbc.com/news/special/2015/newsspec_10857/bbc_news_logo.png?cb=1",
		})
	}
	return newsList
}

func (b BBCRSS) toTimestamp(pubDate string) int {
	s := strings.Split(pubDate, " ")
	month := s[2]
	day := s[1]
	year := s[3]
	tim := s[4]
	timezone := s[5]

	loc := time.FixedZone(timezone, 0)
	t, _ := time.ParseInLocation("02 Jan 2006 15:04:05", day+" "+month+" "+year+" "+tim, loc)

	//fmt.Println(pubDate, " -> ", t.Unix())
	return int(t.Unix())

}
