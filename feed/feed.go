package feed

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"time"
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
	Enclosure  Enclosure  `xml:"enclosure"`
	Link       string     `xml:"link"`
	Guid       string     `xml:"guid"`
	PubDate    PubDate    `xml:"pubDate"`
	Attributes Attributes `xml:"attributes"`
}

// Enclosure contains the url and type of the item
type Enclosure struct {
	Url  string `xml:"url,attr"`
	Type string `xml:"type,attr"`
}

// Attributes represents the json structure inside the attributes node
type Attributes struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

// UnmarshalXML will parse the enclosed json and produce an Attributes element
func (c *Attributes) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)

	var attr Attributes
	err := json.Unmarshal([]byte(v), &attr)
	if err != nil {
		return err
	}

	*c = attr

	return nil
}

// PubDate wraps time.Time to implement the needed interface for the xml unmarshaller
type PubDate struct {
	time.Time
}

// UnmarshalXML will parse the time format in the feed
func (c *PubDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.Parse(time.RFC1123, v)
	if err != nil {
		return err
	}
	*c = PubDate{parse}

	return nil
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
