package fetch

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"errors"

	"github.com/bestform/souparchive/db"
	"github.com/bestform/souparchive/feed"
)

// response is a thin wrapper around http.Response. It is used to mock actual responses in tests
type response struct {
	StatusCode int
	Body       io.ReadCloser
}

// httpClient abstracts the needed interface from the http package to be able to mock it in tests
type httpClient interface {
	Get(string) (*response, error)
}

// defaultHttpClient wraps the corresponding methods from the http package
type defaultHttpClient struct{}

// Get wraps http.Get and produces a response as defined privately in this package
func (d *defaultHttpClient) Get(url string) (*response, error) {
	resp, err := http.Get(url)

	return &response{StatusCode: resp.StatusCode, Body: resp.Body}, err
}

// osLayer abstracts the needed interface from the io and os packages to be able to mock them in tests
type osLayer interface {
	Create(string) (io.ReadWriteCloser, error)
	Copy(io.Writer, io.Reader) (int64, error)
}

// defaultOsLayer wraps the corresponding methods from the io and os packages
type defaultOsLayer struct{}

// Create wraps os.Create
func (d *defaultOsLayer) Create(filename string) (io.ReadWriteCloser, error) {
	return os.Create(filename)
}

// Copy wraps io.Copy
func (d *defaultOsLayer) Copy(w io.Writer, r io.Reader) (int64, error) {
	return io.Copy(w, r)
}

// default setup for live code. Tests will substitute those vars with mocks
var osl osLayer = &defaultOsLayer{}
var httpc httpClient = &defaultHttpClient{}

// Fetch tries to download the item contained in the given feed.Items, if it isn't already in the archive
func Fetch(i feed.Item, a db.Archive) (string, error) {
	if a.Contains(i.Guid) {
		// already in archive
		return "", errors.New(i.Attributes.Url + " already in archive")
	}

	response, err := httpc.Get(i.Attributes.Url)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error fetching %s: %s", i.Attributes.Url, err))
	}
	if response.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("Error fetching %s: Status %d", i.Attributes.Url, response.StatusCode))
	}

	filepath := "archive/" + path.Base(i.Enclosure.Url)
	file, err := osl.Create(filepath)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error opening file %s: %s", filepath, err))
	}

	_, err = osl.Copy(file, response.Body)
	if err != nil {
		response.Body.Close()
		file.Close()
		return "", errors.New(fmt.Sprintf("Error writing file %s: %s", filepath, err))
	}
	response.Body.Close()
	file.Close()

	return i.Guid, nil
}
