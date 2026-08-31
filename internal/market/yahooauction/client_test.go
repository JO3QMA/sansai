package yahooauction

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jo3qma/sansai/internal/model"
)

func TestYahooSearchStatus(t *testing.T) {
	t.Parallel()
	future := strconv.FormatInt(time.Now().Add(24*time.Hour).Unix(), 10)
	past := strconv.FormatInt(time.Now().Add(-24*time.Hour).Unix(), 10)

	if got := yahooSearchStatus(future); got != "open" {
		t.Fatalf("future: got %q", got)
	}
	if got := yahooSearchStatus(past); got != "closed" {
		t.Fatalf("past: got %q", got)
	}
}

func TestJoinDescription(t *testing.T) {
	got := joinDescription([]string{"", "説明1", "", "説明2"})
	if got != "説明1\n説明2" {
		t.Fatalf("got %q", got)
	}
}

func TestItemFromDetailFixture(t *testing.T) {
	raw, err := os.ReadFile("testdata/item_detail.json")
	if err != nil {
		t.Fatal(err)
	}
	var detail auctionItemDetail
	if err := json.Unmarshal(raw, &detail); err != nil {
		t.Fatal(err)
	}

	item := itemFromDetail(detail)
	if item.SaleType != model.SaleTypeAuction {
		t.Fatalf("sale_type: got %q", item.SaleType)
	}
	if item.Description != "説明1\n説明2" {
		t.Fatalf("description: got %q", item.Description)
	}
	if len(item.ImageURLs) != 2 || item.ImageURLs[0] != "https://example.com/img1.jpg" {
		t.Fatalf("image_urls: got %#v", item.ImageURLs)
	}
	if item.EndTime != "2026-08-28T22:25:00+09:00" {
		t.Fatalf("end_time: got %q", item.EndTime)
	}
	if extra, ok := item.Extra.(map[string]any); !ok || extra["end_time"] != nil {
		t.Fatalf("end_time should not be in extra: %#v", item.Extra)
	}
	if !item.IsActive {
		t.Fatal("expected is_active=true for open auction")
	}
}

func TestItemFromDetailFixtureInactive(t *testing.T) {
	raw, err := os.ReadFile("testdata/item_closed.json")
	if err != nil {
		t.Fatal(err)
	}
	var detail auctionItemDetail
	if err := json.Unmarshal(raw, &detail); err != nil {
		t.Fatal(err)
	}

	item := itemFromDetail(detail)
	if item.IsActive {
		t.Fatal("expected is_active=false for closed auction")
	}
}

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
	if items[0].SaleType != model.SaleTypeAuction {
		t.Fatalf("sale_type: got %q", items[0].SaleType)
	}
	if !items[0].IsActive {
		t.Fatal("expected is_active=true for search results")
	}
}
