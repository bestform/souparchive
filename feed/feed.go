package feed

import (
	"encoding/xml"
	"fmt"
)

type Rss struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Enclosure Enclosure `xml:"enclosure"`
	Link      string    `xml:"link"`
	Guid      string    `xml:"guid"`
}

type Enclosure struct {
	Url  string `xml:"url,attr"`
	Type string `xml:"type,attr"`
}

func NewFeedFromXml(input []byte) (Rss) {
	var feed Rss
	xml.Unmarshal(input, &feed)

	return feed
}

func GetFeedUrlForUsername(user string) string {
	return fmt.Sprintf("http://%s.soup.io/rss", user)
}