package sansai

import (
	"context"
	"testing"
)

func TestParseMarkets(t *testing.T) {
	t.Parallel()

	got, err := ParseMarkets("mercari,yahoo_auction")
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0] != MarketMercari || got[1] != MarketYahooAuction {
		t.Fatalf("got %#v", got)
	}

	all, err := ParseMarkets("all")
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != len(AllMarkets()) {
		t.Fatalf("all: got %d markets", len(all))
	}
}

func TestSearchEmptyMarketsFails(t *testing.T) {
	t.Parallel()

	_, err := ParseMarkets("nope")
	if err == nil {
		t.Fatal("expected error for unknown market")
	}
}

func TestGetUnknownMarket(t *testing.T) {
	t.Parallel()

	_, err := Get(context.Background(), Market("unknown"), "id")
	if err == nil {
		t.Fatal("expected error")
	}
}
