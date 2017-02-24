package feed

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
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
	Enclosure        Enclosure `xml:"enclosure"`
	Link             string    `xml:"link"`
	Guid             string    `xml:"guid"`
	PubDate          pubDate   `xml:"pubDate"`
	AttributesSource string    `xml:"attributes"`
	Attributes       Attributes
}

// Enclosure contains the url and type of the item
type Enclosure struct {
	Url  string `xml:"url,attr"`
	Type string `xml:"type,attr"`
}

type Attributes struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

// pubDate wraps time.Time to implement the needed interface for the xml unmarshaller
type pubDate struct {
	time.Time
}

// UnmarshalXML will parse the time format in the feed
func (c *pubDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.Parse(time.RFC1123, v)
	if err != nil {
		fmt.Println(err)
		return err
	}
	*c = pubDate{parse}

	return nil
}

// NewFeedFromXml produces an Rss struct with the information based on the given xml
func NewFeedFromXml(input []byte) Rss {
	var feed Rss

	xml.Unmarshal(input, &feed)

	for i, item := range feed.Channel.Items {
		source := item.AttributesSource
		var attr Attributes
		json.Unmarshal([]byte(source), &attr)
		feed.Channel.Items[i].Attributes = attr
	}

	return feed
}

type localFileLister interface {
	getLocalFilesInfo() ([]os.FileInfo, error)
}

type defaultLocalFileLister struct{}

func (d defaultLocalFileLister) getLocalFilesInfo() ([]os.FileInfo, error) {
	return ioutil.ReadDir("archive")
}

var fileLister localFileLister = defaultLocalFileLister{}

// GetFeedUrlForUsername produces the rss feed url for a given username
func GetFeedUrlForUsername(user string) string {
	return fmt.Sprintf("http://%s.soup.io/rss", user)
}

// GetLocalArchiveFeed creates a feed including all locally archived data
func GetLocalArchiveFeed() (Rss, error) {
	rss := Rss{}
	fileInfos, err := fileLister.getLocalFilesInfo()
	if err != nil {
		return rss, err
	}

	rss.Channel.Items = make([]Item, len(fileInfos))

	for i, info := range fileInfos {
		item := Item{}
		enc := Enclosure{Url: info.Name()}
		item.Enclosure = enc
		rss.Channel.Items[i] = item
	}

	return rss, nil
}
