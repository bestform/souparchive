package feed

import (
	"os"
	"testing"
	"time"
)

func TestUnmarshallingXml(t *testing.T) {
	input := `<?xml version="1.0" encoding="UTF-8"?>
<rss xmlns:soup="http://www.soup.io/rss" version="2.0">
  <channel>
    <title>Testtitle</title>
    <link>Testlink</link>
    <description>Testdescription</description>
    <item>
      <enclosure url="enc1Url" type="enc1Type" />
       <soup:attributes>
	 {"type":"attrType1","url":"attrUrl1"}
       </soup:attributes>
       <pubDate>Thu, 23 Feb 2017 14:14:29 GMT</pubDate>
       <link>Item1Link</link>
       <guid>Item1GUID</guid>
    </item>
    <item>
      <enclosure url="enc2Url" type="enc2Type" />
       <soup:attributes>
	 {"type":"attrType2","url":"attrUrl2"}
       </soup:attributes>
       <pubDate>Fri, 24 Feb 2017 14:14:29 GMT</pubDate>
       <link>Item2Link</link>
       <guid>Item2GUID</guid>
    </item>
  </channel>
</rss>
`
	result := NewFeedFromXml([]byte(input))

	check(result.Channel.Title, "Testtitle", t)
	check(result.Channel.Link, "Testlink", t)
	check(result.Channel.Description, "Testdescription", t)

	if len(result.Channel.Items) != 2 {
		t.Fatalf("Expected 2 items, but got %d", len(result.Channel.Items))
	}

	loc, err := time.LoadLocation("GMT")
	if err != nil {
		t.Fatal("Could not load location for time comparison. Error in test!", err)
	}
	check(result.Channel.Items[0].Guid, "Item1GUID", t)
	check(result.Channel.Items[0].Link, "Item1Link", t)
	checkTime(result.Channel.Items[0].PubDate, time.Date(2017, time.February, 23, 14, 14, 29, 0, loc), t)
	check(result.Channel.Items[0].Enclosure.Url, "enc1Url", t)
	check(result.Channel.Items[0].Enclosure.Type, "enc1Type", t)
	check(result.Channel.Items[0].Attributes.Type, "attrType1", t)
	check(result.Channel.Items[0].Attributes.Url, "attrUrl1", t)

	check(result.Channel.Items[1].Guid, "Item2GUID", t)
	check(result.Channel.Items[1].Link, "Item2Link", t)
	checkTime(result.Channel.Items[1].PubDate, time.Date(2017, time.February, 24, 14, 14, 29, 0, loc), t)
	check(result.Channel.Items[1].Enclosure.Url, "enc2Url", t)
	check(result.Channel.Items[1].Enclosure.Type, "enc2Type", t)
	check(result.Channel.Items[1].Attributes.Type, "attrType2", t)
	check(result.Channel.Items[1].Attributes.Url, "attrUrl2", t)
}

func check(actual, expected string, t *testing.T) {
	if actual != expected {
		t.Fatalf("expected: '%s' but got '%s'", expected, actual)
	}
}

func checkTime(actual pubDate, expected time.Time, t *testing.T) {
	if !actual.Equal(expected) {
		t.Fatalf("Dates not equivalent. Expected: %v, got: %v", expected, actual)
	}
}

func TestUrlCreation(t *testing.T) {
	url := GetFeedUrlForUsername("foo")
	if url != "http://foo.soup.io/rss" {
		t.Fatalf("Wrong feed URL: %s", url)
	}
}

type testFileLister struct {
	filenames []string
	err       error
}

type testFile struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
	sys     interface{}
}

func (t testFile) Name() string       { return t.name }
func (t testFile) Size() int64        { return t.size }
func (t testFile) Mode() os.FileMode  { return t.mode }
func (t testFile) ModTime() time.Time { return t.modTime }
func (t testFile) IsDir() bool        { return t.isDir }
func (t testFile) Sys() interface{}   { return t.sys }

func (t testFileLister) getLocalFilesInfo() ([]os.FileInfo, error) {
	if t.err != nil {
		return nil, t.err
	}

	fileInfos := make([]os.FileInfo, len(t.filenames))
	for i, name := range t.filenames {
		fileInfos[i] = testFile{name: name}
	}

	return fileInfos, nil
}

func TestGetLocalArchiveFeed(t *testing.T) {
	fileLister = testFileLister{filenames: []string{"foo.jpg", "bar.gif"}}

	feed, err := GetLocalArchiveFeed()
	if err != nil {
		t.Fatal(err)
	}

	if len(feed.Channel.Items) != 2 {
		t.Fatalf("Expected 2 elements in feed, got %d", len(feed.Channel.Items))
	}

	expectedNames := []string{"foo.jpg", "bar.gif"}

	for i, item := range feed.Channel.Items {
		if item.Enclosure.Url != expectedNames[i] {
			t.Fatalf("Expected url %s on position %d, but got %s", expectedNames[i], i, item.Enclosure.Url)
		}
	}

}
