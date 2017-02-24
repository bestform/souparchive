package fetch

import (
	"io"
	"testing"

	"net/http"

	"time"

	"github.com/bestform/souparchive/db"
	"github.com/bestform/souparchive/feed"
	"github.com/pkg/errors"
)

// testFile implements an io.ReadWriteCloser for use this the testOsLayer
type testFile struct{}

func (t testFile) Read(p []byte) (n int, err error) {
	return 0, nil
}
func (t testFile) Write(p []byte) (n int, err error) {
	return 0, nil
}
func (t testFile) Close() error {
	return nil
}

type testOsLayer struct {
	created         string
	copyCalledTimes int
}

func (d *testOsLayer) Create(filename string) (io.ReadWriteCloser, error) {
	d.created = filename
	return testFile{}, nil
}
func (d *testOsLayer) Copy(w io.Writer, r io.Reader) (int64, error) {
	d.copyCalledTimes++
	return 0, nil
}

type testHttpClient struct {
	askedForUrl string
	getError    error
	response    response
}

func (d *testHttpClient) Get(url string) (*response, error) {
	d.askedForUrl = url
	return &d.response, d.getError
}

type testBody struct{}

func (t *testBody) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (t *testBody) Close() error {
	return nil
}

func TestReturnOnItemAlreadyInArchive(t *testing.T) {
	a := db.Archive{}
	a.Data.Items = append(a.Data.Items, db.Item{"foo", 0})
	i := feed.Item{}
	i.Guid = "foo"

	_, _, err := Fetch(i, a)
	if err == nil {
		t.Fatal("Expected error on item already in archive, but got none")
	}
}

func TestFetchOfCorrectUrl(t *testing.T) {
	osl = &testOsLayer{}
	mockHttpClient := &testHttpClient{}
	httpc = mockHttpClient
	a := db.Archive{}
	i := feed.Item{}
	i.Guid = "foo"
	i.Attributes.Url = "testURL"

	Fetch(i, a)
	if mockHttpClient.askedForUrl != "testURL" {
		t.Fatalf("Expected http get on %s but got %s", "testURL", mockHttpClient.askedForUrl)
	}
}

func TestErrorOnGetError(t *testing.T) {
	osl = &testOsLayer{}
	mockHttpClient := &testHttpClient{}
	mockHttpClient.getError = errors.New("error on get")
	httpc = mockHttpClient
	a := db.Archive{}
	i := feed.Item{}

	_, _, err := Fetch(i, a)
	if err == nil {
		t.Fatal("Expected error on http get error, but got nil")
	}

}

func TestErrorOnBadHttpStatus(t *testing.T) {
	osl = &testOsLayer{}
	mockHttpClient := &testHttpClient{}
	mockHttpClient.response.StatusCode = http.StatusBadGateway
	httpc = mockHttpClient
	a := db.Archive{}
	i := feed.Item{}

	_, _, err := Fetch(i, a)
	if err == nil {
		t.Fatal("Expected error on bad http status, but got nil")
	}
}

func TestCreateCorrectFile(t *testing.T) {
	mockOsLayer := testOsLayer{}
	osl = &mockOsLayer
	mockHttpClient := &testHttpClient{}
	mockHttpClient.response.StatusCode = http.StatusOK
	mockHttpClient.response.Body = &testBody{}
	httpc = mockHttpClient
	a := db.Archive{}
	i := feed.Item{}
	i.Enclosure.Url = "foo/bar/baz"

	Fetch(i, a)
	if mockOsLayer.created != "archive/baz" {
		t.Fatalf("Expected file created to be %s, but got %s", "archive/baz", mockOsLayer.created)
	}
}

func TestCopyFromResponse(t *testing.T) {
	mockOsLayer := testOsLayer{}
	osl = &mockOsLayer
	mockHttpClient := &testHttpClient{}
	mockHttpClient.response.StatusCode = http.StatusOK
	mockHttpClient.response.Body = &testBody{}
	httpc = mockHttpClient
	a := db.Archive{}
	i := feed.Item{}
	i.Enclosure.Url = "foo/bar/baz"

	Fetch(i, a)
	if mockOsLayer.copyCalledTimes != 1 {
		t.Fatal("Expected data to be copied one time but got", mockOsLayer.copyCalledTimes)
	}
}

func TestReturnOfGuidOnSuccessfulCopy(t *testing.T) {
	mockOsLayer := testOsLayer{}
	osl = &mockOsLayer
	mockHttpClient := &testHttpClient{}
	mockHttpClient.response.StatusCode = http.StatusOK
	mockHttpClient.response.Body = &testBody{}
	httpc = mockHttpClient
	a := db.Archive{}
	i := feed.Item{}
	i.Guid = "foo"
	i.PubDate = feed.PubDate{time.Unix(100, 0)}

	guid, timestamp, err := Fetch(i, a)
	if err != nil {
		t.Fatal("Expected return of guid without error but got", err)
	}

	if guid != "foo" {
		t.Fatal("Expected returned guid to be 'foo', but got", guid)
	}

	if timestamp != 100 {
		t.Fatal("Expected timestamp to be 100, got", timestamp)
	}
}
