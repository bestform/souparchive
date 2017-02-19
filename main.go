package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime/trace"
	"sync"

	"github.com/bestform/souparchive/db"
	"github.com/bestform/souparchive/feed"
	"github.com/bestform/souparchive/fetch"
)

// DEBUG will write a trace if set to true. The only way to set this to true is to manipulate this very code
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
		go func(i feed.Item, a db.Archive, c chan string) {
			defer wg.Done()
			fmt.Printf("Saving %s...\n", i.Enclosure.Url)
			guid, err := fetch.Fetch(i, a)
			if err != nil {
				fmt.Println(err)
				return
			}
			c <- guid
		}(i, a, c)
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
