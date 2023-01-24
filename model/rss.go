package model

import (
	"encoding/xml"
	"strings"
	"time"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr,omitempty"`
	Content string   `xml:"content,attr,omitempty"`
	Wfw     string   `xml:"wfw,attr,omitempty"`
	Dc      string   `xml:"dc,attr,omitempty"`
	Atom    string   `xml:"atom,attr,omitempty"`
	Sy      string   `xml:"sy,attr,omitempty"`
	Slash   string   `xml:"slash,attr,omitempty"`
	Channel struct {
		Text  string `xml:",chardata"`
		Title string `xml:"title"`
		Link  struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr,omitempty"`
			Rel  string `xml:"rel,attr,omitempty"`
			Type string `xml:"type,attr,omitempty"`
		} `xml:"link"`
		Description     string `xml:"description"`
		LastBuildDate   string `xml:"lastBuildDate"`
		Language        string `xml:"language"`
		UpdatePeriod    string `xml:"updatePeriod"`
		UpdateFrequency string `xml:"updateFrequency"`
		Generator       string `xml:"generator"`
		Image           struct {
			Text   string `xml:",chardata"`
			URL    string `xml:"url"`
			Title  string `xml:"title"`
			Link   string `xml:"link"`
			Width  string `xml:"width"`
			Height string `xml:"height"`
		} `xml:"image"`
		Item []struct {
			Text     string   `xml:",chardata"`
			Title    string   `xml:"title,omitempty"`
			Link     string   `xml:"link,omitempty"`
			Comments string   `xml:"comments,omitempty"`
			Creator  string   `xml:"creator,omitempty"`
			PubDate  string   `xml:"pubDate,omitempty"`
			Category []string `xml:"category,omitempty"`
			Guid     struct {
				Text        string `xml:",chardata"`
				IsPermaLink string `xml:"isPermaLink,attr,omitempty"`
			} `xml:"guid,omitempty"`
			Description string `xml:"description,omitempty"`
			Encoded     string `xml:"encoded,omitempty"`
			CommentRss  string `xml:"commentRss,omitempty"`
		} `xml:"item"`
	} `xml:"channel"`
}

func (r RSS) RSSToNewsClNewsList() NewsClNewsList {
	var newsList NewsClNewsList
	for _, item := range r.Channel.Item {
		newsList.News = append(newsList.News, NewsClNews{
			ID:          item.Guid.Text,
			Title:       item.Title,
			Description: item.Description,
			URL:         item.Link,
			Time:        r.toTimestamp(item.PubDate),
			Location:    "UK",
			Language:    "EN",
			Type:        2,
			SourceName:  "Reuters",
			SourceURL:   "https://www.reuters.com",
			Tags:        "",
			ImageURL:    "https://www.reuters.com/resources_v2/images/reuters125.png",
		})
	}
	return newsList
}

//Convert pub date to unix timestamp
//format: 12 Nov 2022 23:18:17 +0000

func (r RSS) toTimestamp(pubDate string) int {
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
