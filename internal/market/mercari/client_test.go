package mercari

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/jo3qma/sansai/internal/model"
)

func TestIsMercariPersonalID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		id   string
		want bool
	}{
		{id: "m29820723477", want: true},
		{id: "m1", want: true},
		{id: "2JV9ANthhYNHiJGoq39omJ", want: false},
		{id: "uC2mpJ9BJhcPaSNxqEW3wN", want: false},
		{id: "mabc", want: false},
		{id: "", want: false},
	}
	for _, tt := range tests {
		if got := isMercariPersonalID(tt.id); got != tt.want {
			t.Fatalf("isMercariPersonalID(%q) = %v, want %v", tt.id, got, tt.want)
		}
	}
}

func TestGetUnsupportedShopsID(t *testing.T) {
	t.Parallel()
	c := &Client{}
	_, err := c.Get(context.Background(), "2JV9ANthhYNHiJGoq39omJ")
	if err == nil {
		t.Fatal("expected error")
	}
	if got := err.Error(); got != "unsupported id: mercari shops listing (2JV9ANthhYNHiJGoq39omJ)" {
		t.Fatalf("got %q", got)
	}
}

func TestItemFromDetailFixture(t *testing.T) {
	t.Run("fixed_price", func(t *testing.T) {
		raw := readFixture(t, "testdata/item_fixed_price.json")
		var detail mercariItemDetail
		json.Unmarshal(raw, &detail)
		item := itemFromDetail(detail)

		if item.SaleType != model.SaleTypeFixedPrice {
			t.Fatalf("sale_type: got %q", item.SaleType)
		}
		if item.Description != "固定価格の説明" {
			t.Fatalf("description: got %q", item.Description)
		}
		if len(item.ImageURLs) != 2 {
			t.Fatalf("image_urls: got %#v", item.ImageURLs)
		}
		if item.EndTime != "" {
			t.Fatalf("end_time should be empty: %q", item.EndTime)
		}
		if !item.IsActive {
			t.Fatal("expected is_active=true for on_sale")
		}
	})

	t.Run("auction_with_bids", func(t *testing.T) {
		raw := readFixture(t, "testdata/item_auction_bids.json")
		var detail mercariItemDetail
		json.Unmarshal(raw, &detail)
		item := itemFromDetail(detail)

		if item.SaleType != model.SaleTypeAuction {
			t.Fatalf("sale_type: got %q", item.SaleType)
		}
		if item.EndTime == "" {
			t.Fatal("expected end_time")
		}
	})

	t.Run("auction_no_bid", func(t *testing.T) {
		raw := readFixture(t, "testdata/item_auction_no_bid.json")
		var detail mercariItemDetail
		json.Unmarshal(raw, &detail)
		item := itemFromDetail(detail)

		if item.SaleType != model.SaleTypeAuction {
			t.Fatalf("sale_type: got %q", item.SaleType)
		}
		if item.EndTime != "" {
			t.Fatalf("end_time should be omitted before first bid: %q", item.EndTime)
		}
	})

	t.Run("sold_out", func(t *testing.T) {
		raw := readFixture(t, "testdata/item_sold_out.json")
		var detail mercariItemDetail
		json.Unmarshal(raw, &detail)
		item := itemFromDetail(detail)

		if item.IsActive {
			t.Fatal("expected is_active=false for sold_out")
		}
	})
}

func TestUnixRFC3339(t *testing.T) {
	got := unixRFC3339(1787915990)
	if got == "" {
		t.Fatal("expected RFC3339 time")
	}
	if unixRFC3339(0) != "" {
		t.Fatal("zero should be empty")
	}
}

func readFixture(t *testing.T, path string) []byte {
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func TestSearchIntegration(t *testing.T) {
	if os.Getenv("SANSAI_INTEGRATION") == "" {
		t.Skip("set SANSAI_INTEGRATION=1 to run")
	}
	c := &Client{}
	res, err := c.Search(context.Background(), "nintendo switch", model.SearchOptions{Limit: 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Items) == 0 {
		t.Fatal("expected items")
	}
	b, _ := json.MarshalIndent(res, "", "  ")
	t.Log(string(b))
}
