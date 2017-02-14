package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
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
	Guid      string    `xml:"guid"`
}

type enclosure struct {
	Url  string `xml:"url,attr"`
	Type string `xml:"type,attr"`
}

type archive struct {
	Guid []string `json:"guid"`
}

func main() {

	accountPtr := flag.String("user", "", "soup.io username")
	flag.Parse()

	if *accountPtr == "" {
		flag.Usage()
		os.Exit(0)
	}

	url := fmt.Sprintf("http://%s.soup.io/rss", *accountPtr)
	feedResponse, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching %s: %s", url, err)
		os.Exit(1)
	}
	defer feedResponse.Body.Close()
	feedBody, err := ioutil.ReadAll(feedResponse.Body)
	if err != nil {
		fmt.Printf("Error reading feed from %s: %s", url, err)
		os.Exit(1)
	}

	var feed rss
	xml.Unmarshal(feedBody, &feed)

	a, err := readArchive()
	if err != nil {
		fmt.Println("No archive data found. Will create a fresh one")
	}

	for _, i := range feed.Channel.Items {
		if inArchive(i.Guid, a) {
			fmt.Printf("Skipping %s. Already in archive.\n", i.Enclosure.Url)
			continue
		}
		response, err := http.Get(i.Enclosure.Url)
		if err != nil {
			fmt.Printf("Error fetching %s: %s\n", i.Enclosure.Url, err)
			continue
		}
		filepath := "archive/" + path.Base(i.Enclosure.Url)
		file, err := os.Create(filepath)
		if err != nil {
			fmt.Printf("Error opening file %s: %s\n", filepath, err)
			response.Body.Close()
			continue
		}
		fmt.Printf("Saving %s...\n", i.Enclosure.Url)
		_, err = io.Copy(file, response.Body)
		if err != nil {
			fmt.Printf("Error writing file %s: %s\n", filepath, err)
			response.Body.Close()
			continue
		}
		response.Body.Close()
		file.Close()

		a.Guid = append(a.Guid, i.Guid)
		saveArchive(a)
	}
}

func saveArchive(a archive) {
	data, err := json.Marshal(a)
	if err != nil {
		fmt.Println("Error marshalling archive data: ", err)
		return
	}
	ioutil.WriteFile("archive/guids.json", data, 0600)
}

func inArchive(guid string, a archive) bool {
	for _, s := range a.Guid {
		if guid == s {
			return true
		}
	}

	return false
}

func readArchive() (archive, error) {
	var a archive

	archiveData, err := ioutil.ReadFile("archive/guids.json")
	if err != nil {
		return a, err
	}

	err = json.Unmarshal(archiveData, &a)
	if err != nil {
		return a, err
	}

	return a, nil
}
