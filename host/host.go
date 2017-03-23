package host

import (
	"fmt"
	"net/http"

	"github.com/bestform/souparchive/db"
)

type entry struct {
	url string
}

var localFeed []db.Item

// Host will host the current archive on localhost via the specified port
func Host(port string) error {
	archive := db.NewArchive("archive/archive.json")
	archive.Read()

	localFeed = archive.Data.Items

	http.HandleFunc("/", hostList)

	return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func hostList(w http.ResponseWriter, r *http.Request) {

	for _, item := range localFeed {
		w.Write([]byte(item.Filename))
		w.Write([]byte("\n"))
	}
}
