package feed

import (
	"encoding/xml"
	"fmt"
)

// Rss it the root node of the rss feed containing just one Channel node
type Rss struct {
	Channel Channel `xml:"channel"`
}

// Channel contains all the meta information for the channel and a list of items
type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

// Item is one entry in the feed
type Item struct {
	Enclosure Enclosure `xml:"enclosure"`
	Link      string    `xml:"link"`
	Guid      string    `xml:"guid"`
}

// Enclosure contains the url and type of the item
type Enclosure struct {
	Url  string `xml:"url,attr"`
	Type string `xml:"type,attr"`
}

// NewFeedFromXml produces an Rss struct with the information based on the given xml
func NewFeedFromXml(input []byte) Rss {
	var feed Rss
	xml.Unmarshal(input, &feed)

	return feed
}

// GetFeedUrlForUsername produces the rss feed url for a given username
func GetFeedUrlForUsername(user string) string {
	return fmt.Sprintf("http://%s.soup.io/rss", user)
}
