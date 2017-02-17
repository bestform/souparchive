package feed

import (
	"testing"
)

func TestUnmarshallingXml(t *testing.T) {
	input := `<?xml version="1.0" encoding="UTF-8"?>
<rss>
  <channel>
    <title>Testtitle</title>
    <link>Testlink</link>
    <description>Testdescription</description>
    <item>
      <enclosure url="enc1Url" type="enc1Type" />
       <link>Item1Link</link>
       <guid>Item1GUID</guid>
    </item>
    <item>
      <enclosure url="enc2Url" type="enc2Type" />
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

	check(result.Channel.Items[0].Guid, "Item1GUID", t)
	check(result.Channel.Items[0].Link, "Item1Link", t)
	check(result.Channel.Items[0].Enclosure.Url, "enc1Url", t)
	check(result.Channel.Items[0].Enclosure.Type, "enc1Type", t)

	check(result.Channel.Items[1].Guid, "Item2GUID", t)
	check(result.Channel.Items[1].Link, "Item2Link", t)
	check(result.Channel.Items[1].Enclosure.Url, "enc2Url", t)
	check(result.Channel.Items[1].Enclosure.Type, "enc2Type", t)
}

func check(actual, expected string, t *testing.T) {
	if actual != expected {
		t.Fatalf("expected: %s but got %s", expected, actual)
	}
}

func TestUrlCreation(t *testing.T) {
	url := GetFeedUrlForUsername("foo")
	if url != "http://foo.soup.io/rss" {
		t.Fatalf("Wrong feed URL: %s", url)
	}
}
