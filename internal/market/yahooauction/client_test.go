package yahooauction

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jo3qma/sansai/internal/model"
)

func TestHtmlToText(t *testing.T) {
	got := htmlToText("line1<BR><br/>line2<b>bold</b>")
	want := "line1\n\nline2bold"
	if got != want {
		t.Fatalf("got %q", got)
	}
}

func TestAuctionDescriptionPrefersHTML(t *testing.T) {
	raw := auctionItemDetail{
		Description:     []string{"短い要約"},
		DescriptionHTML: "全文<BR>続き",
	}
	if got := auctionDescription(raw); got != "全文\n続き" {
		t.Fatalf("got %q", got)
	}
}

func TestDescriptionTruncatedFlag(t *testing.T) {
	raw := auctionItemDetail{
		Description: []string{"短い"},
		DescriptionUlt: &struct {
			RawDescriptionLength int `json:"rawDescriptionLength"`
		}{RawDescriptionLength: 500},
	}
	extra := auctionExtra(raw)
	if extra["description_truncated"] != true {
		t.Fatalf("extra: %#v", extra)
	}
}

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

func TestYahooAuctionPageSize(t *testing.T) {
	t.Parallel()
	tests := []struct {
		limit int
		want  int
	}{
		{limit: 1, want: 20},
		{limit: 8, want: 20},
		{limit: 20, want: 20},
		{limit: 21, want: 50},
		{limit: 50, want: 50},
		{limit: 51, want: 100},
		{limit: 100, want: 100},
	}
	for _, tt := range tests {
		if got := yahooAuctionPageSize(tt.limit); got != tt.want {
			t.Fatalf("yahooAuctionPageSize(%d) = %d, want %d", tt.limit, got, tt.want)
		}
	}
}

func TestSearchOffsetUsesPageSize(t *testing.T) {
	t.Parallel()
	pageSize := yahooAuctionPageSize(8)
	offset := (2-1)*pageSize + 1
	if offset != 21 {
		t.Fatalf("page 2 with limit 8: offset = %d, want 21 (pageSize=%d)", offset, pageSize)
	}
}

func TestBuildSearchParams(t *testing.T) {
	t.Parallel()
	opts := model.SearchOptions{MinPrice: 10000, MaxPrice: 80000}
	full := buildSearchParams("test query", opts, 21, 20, searchParamFull)
	if full.Get("va") != "test query" || full.Get("min") != "10000" {
		t.Fatalf("full: %#v", full)
	}
	noVA := buildSearchParams("test query", opts, 21, 20, searchParamNoVA)
	if noVA.Get("va") != "" || noVA.Get("min") != "10000" {
		t.Fatalf("noVA: %#v", noVA)
	}
	noPrice := buildSearchParams("test query", opts, 21, 20, searchParamNoPrice)
	if noPrice.Get("min") != "" || noPrice.Get("max") != "" || noPrice.Get("p") != "test query" {
		t.Fatalf("noPrice: %#v", noPrice)
	}
}

func TestSearchNotFoundReturnsEmpty(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	t.Cleanup(server.Close)

	oldBase := searchBase
	searchBase = server.URL
	t.Cleanup(func() { searchBase = oldBase })

	c := &Client{}
	res, err := c.Search(context.Background(), "ストレージサーバー 8ベイ", model.SearchOptions{
		Limit:    8,
		MinPrice: 10000,
		MaxPrice: 80000,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Items) != 0 {
		t.Fatalf("items: %#v", res.Items)
	}
	if res.Market != model.MarketYahooAuction {
		t.Fatalf("market: %q", res.Market)
	}
}

func TestAuctionDescriptionPreservesFullHTMLText(t *testing.T) {
	raw := auctionItemDetail{
		Description:     []string{"短い要約"},
		DescriptionHTML:   "型番: PowerEdge R740<BR>詳細な説明文",
		DescriptionUlt: &struct {
			RawDescriptionLength int `json:"rawDescriptionLength"`
		}{RawDescriptionLength: 500},
	}
	got := auctionDescription(raw)
	if !strings.Contains(got, "PowerEdge R740") || !strings.Contains(got, "詳細な説明文") {
		t.Fatalf("got %q", got)
	}
	extra := auctionExtra(raw)
	if _, ok := extra["description_truncated"]; ok {
		t.Fatalf("description_truncated should not be set when HTML is present: %#v", extra)
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
