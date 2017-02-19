package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime/trace"
	"sync"

	"github.com/bestform/souparchive/db"
	"github.com/bestform/souparchive/feed"
)

var DEBUG = false

func main() {
	if DEBUG {
		file, err := os.Create("trace.out")
		if err != nil {
			panic("Error creating trace file")
		}
		defer file.Close()
		trace.Start(file)
		defer trace.Stop()
	}

	accountPtr := flag.String("user", "", "soup.io username")
	flag.Parse()

	if *accountPtr == "" {
		flag.Usage()
		os.Exit(0)
	}

	url := feed.GetFeedUrlForUsername(*accountPtr)
	feedResponse, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching %s: %s", url, err)
		os.Exit(1)
	}
	defer feedResponse.Body.Close()
	feedBody, err := ioutil.ReadAll(feedResponse.Body)
	if err != nil {
		fmt.Printf("Error reading rssFeed from %s: %s", url, err)
		os.Exit(1)
	}

	rssFeed := feed.NewFeedFromXml(feedBody)

	var wg sync.WaitGroup

	a := db.NewArchive("archive/guids.json")
	a.Read()

	c := make(chan string)

	for _, i := range rssFeed.Channel.Items {
		wg.Add(1)
		go func(i feed.Item, c chan string) {
			defer wg.Done()
			if a.Contains(i.Guid) {
				// already in archive
				return
			}
			fmt.Printf("Saving %s...\n", i.Enclosure.Url)

			response, err := http.Get(i.Enclosure.Url)
			if err != nil {
				fmt.Printf("Error fetching %s: %s\n", i.Enclosure.Url, err)
				return
			}
			if response.StatusCode != http.StatusOK {
				fmt.Printf("Error fetching %s: Status %d\n", i.Enclosure.Url, response.StatusCode)
				return
			}

			filepath := "archive/" + path.Base(i.Enclosure.Url)
			file, err := os.Create(filepath)
			if err != nil {
				fmt.Printf("Error opening file %s: %s\n", filepath, err)
				return
			}

			_, err = io.Copy(file, response.Body)
			if err != nil {
				fmt.Printf("Error writing file %s: %s\n", filepath, err)
				response.Body.Close()
				file.Close()
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
			a.Read()
			a.Add(guid)
			err := a.Persist()
			if err != nil {
				fmt.Println("error persisting database", err)
			}
		}
		waitForArchive <- true
	}(c)

	wg.Wait()
	close(c)
	<-waitForArchive
}
