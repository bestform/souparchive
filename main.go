package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type rss struct {
	Channel channel `xml:"channel"`
}

type channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []item `xml:"item"`
}

type item struct {
	Enclosure enclosure `xml:"enclosure"`
	Link      string    `xml:"link"`
}

type enclosure struct {
	Url  string `xml:"url,attr"`
	Type string `xml:"type,attr"`
}

type attributes struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

func main() {
	xmlFile, err := ioutil.ReadFile("bestform.xml")
	if err != nil {
		panic(err)
	}

	var feed rss
	xml.Unmarshal(xmlFile, &feed)

	fmt.Printf("%q", feed)

}
