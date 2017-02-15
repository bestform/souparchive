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
	"sync"
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

	var wg sync.WaitGroup

	a, err := readArchive()
	if err != nil {
		// no archive yet. it will be created
		a = archive{}
	}

	if _, err := os.Stat("archive"); os.IsNotExist(err) {
		err = os.Mkdir("archive", 0755)
		if err != nil {
			fmt.Printf("Error creating archive: %s\n", err)
			os.Exit(1)
		}
	}

	c := make(chan string)

	for _, i := range feed.Channel.Items {
		wg.Add(1)
		go func(i item, c chan string) {
			defer wg.Done()
			if inArchive(i.Guid, a) {
				fmt.Printf("Skipping %s. Already in archive.\n", i.Enclosure.Url)
				return
			}
			fmt.Printf("Saving %s...\n", i.Enclosure.Url)

			filepath := "archive/" + path.Base(i.Enclosure.Url)
			file, err := os.Create(filepath)
			if err != nil {
				fmt.Printf("Error opening file %s: %s\n", filepath, err)
				return
			}

			response, err := http.Get(i.Enclosure.Url)
			if err != nil {
				fmt.Printf("Error fetching %s: %s\n", i.Enclosure.Url, err)
				return
			}
			_, err = io.Copy(file, response.Body)
			if err != nil {
				fmt.Printf("Error writing file %s: %s\n", filepath, err)
				response.Body.Close()
				return
			}
			response.Body.Close()
			file.Close()

			c <- i.Guid
		}(i, c)
	}

	waitForArchive := make(chan bool)
	go func(c chan string) {
		for guid := range c {
			addToArchive(guid)
		}
		waitForArchive <- true
	}(c)

	wg.Wait()
	close(c)
	<-waitForArchive
}

func addToArchive(guid string) {
	a, err := readArchive()
	if err != nil {
		a = archive{}
	}
	a.Guid = append(a.Guid, guid)
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
