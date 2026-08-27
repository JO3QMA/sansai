package yahooflea

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/jo3qma/sansai/internal/model"
)

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
