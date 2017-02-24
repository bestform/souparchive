package host

import (
	"github.com/bestform/souparchive/feed"
	"net/http"
	"fmt"
)

type entry struct {
	url string
}

var localFeed feed.Rss

// Host wil host the current archive on localhost via the specified port
func Host(port string) error {
	var err error
	localFeed, err = feed.GetLocalArchiveFeed()
	if err != nil {
		return err
	}

	http.HandleFunc("/", hostList)

	return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func hostList(w http.ResponseWriter, r *http.Request) {
	for _, item := range localFeed.Channel.Items {
		w.Write([]byte(item.Enclosure.Url))
		w.Write([]byte("\n"))
	}
}

