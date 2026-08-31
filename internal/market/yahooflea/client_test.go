package yahooflea

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jo3qma/sansai/internal/model"
)

func TestSearchNotFoundReturnsEmpty(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	t.Cleanup(server.Close)

	oldBase := searchBase
	searchBase = server.URL + "/"
	t.Cleanup(func() { searchBase = oldBase })

	c := &Client{}
	res, err := c.Search(context.Background(), "6028U", model.SearchOptions{Limit: 5})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Items) != 0 {
		t.Fatalf("items: %#v", res.Items)
	}
	if res.Market != model.MarketYahooFlea || res.Query != "6028U" {
		t.Fatalf("result: %#v", res)
	}
}

func TestItemFromDetailFixture(t *testing.T) {
	raw, err := os.ReadFile("testdata/item_detail.json")
	if err != nil {
		t.Fatal(err)
	}
	var detail fleaItemDetail
	if err := json.Unmarshal(raw, &detail); err != nil {
		t.Fatal(err)
	}

	item := itemFromDetail(detail)
	if item.SaleType != model.SaleTypeFixedPrice {
		t.Fatalf("sale_type: got %q", item.SaleType)
	}
	if item.Description != "フリマ説明文" {
		t.Fatalf("description: got %q", item.Description)
	}
	if len(item.ImageURLs) != 2 {
		t.Fatalf("image_urls: got %#v", item.ImageURLs)
	}
	if extra, ok := item.Extra.(map[string]any); ok {
		if _, exists := extra["description"]; exists {
			t.Fatalf("description should not be in extra")
		}
	}
	if !item.IsActive {
		t.Fatal("expected is_active=true for open listing")
	}
}

func TestItemFromDetailFixtureInactive(t *testing.T) {
	raw, err := os.ReadFile("testdata/item_sold.json")
	if err != nil {
		t.Fatal(err)
	}
	var detail fleaItemDetail
	if err := json.Unmarshal(raw, &detail); err != nil {
		t.Fatal(err)
	}

	item := itemFromDetail(detail)
	if item.IsActive {
		t.Fatal("expected is_active=false for sold listing")
	}
}
