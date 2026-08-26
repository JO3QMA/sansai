package yahooauction

import (
	"os"
	"testing"
)

func TestParseSearchHTML(t *testing.T) {
	html, err := os.ReadFile("/tmp/yahoo_auction.html")
	if err != nil {
		t.Skip("missing /tmp/yahoo_auction.html")
	}
	items := parseSearchHTML(string(html), 0, 0)
	if len(items) == 0 {
		t.Fatal("expected items")
	}
	if items[0].ID == "" || items[0].Title == "" {
		t.Fatalf("unexpected first item: %+v", items[0])
	}
}
