package host

import (
	"fmt"
	"net/http"

	"sort"

	"html/template"

	"github.com/bestform/souparchive/db"
)

type entry struct {
	url string
}

var localFeed []db.Item

type ByTime []db.Item

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].Timestamp > a[j].Timestamp }

// Host will host the current archive on localhost via the specified port
func Host(port string) error {
	archive := db.NewArchive("archive/archive.json")
	archive.Read()

	localFeed = archive.Data.Items
	sort.Sort(ByTime(localFeed))

	fs := http.FileServer(http.Dir("archive"))
	http.Handle("/images/", http.StripPrefix("/images/", fs))
	http.HandleFunc("/", hostList)

	return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func hostList(w http.ResponseWriter, r *http.Request) {

	t, err := template.New("index.html").ParseFiles("host/templates/index.html")
	if err != nil {
		panic(err)
	}
	err = t.Execute(w, localFeed)
	if err != nil {
		panic(err)
	}

}
