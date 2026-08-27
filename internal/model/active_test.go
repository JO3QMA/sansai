package model

import "testing"

func TestIsActive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		market       Market
		status       string
		auctionState string
		want         bool
	}{
		// mercari
		{name: "mercari on_sale", market: MarketMercari, status: "on_sale", want: true},
		{name: "mercari auction ongoing", market: MarketMercari, status: "on_sale", auctionState: "STATE_ONGOING", want: true},
		{name: "mercari auction no bid", market: MarketMercari, status: "on_sale", auctionState: "STATE_NO_BID", want: true},
		{name: "mercari sold_out", market: MarketMercari, status: "sold_out", want: false},
		{name: "mercari trading", market: MarketMercari, status: "trading", want: false},
		{name: "mercari cancel", market: MarketMercari, status: "cancel", want: false},
		{name: "mercari auction ended", market: MarketMercari, status: "on_sale", auctionState: "STATE_ENDED", want: false},

		// yahoo auction
		{name: "yahoo auction open", market: MarketYahooAuction, status: "open", want: true},
		{name: "yahoo auction closed", market: MarketYahooAuction, status: "closed", want: false},
		{name: "yahoo auction sold", market: MarketYahooAuction, status: "sold", want: false},
		{name: "yahoo auction cancel", market: MarketYahooAuction, status: "cancel", want: false},

		// yahoo flea
		{name: "yahoo flea open", market: MarketYahooFlea, status: "open", want: true},
		{name: "yahoo flea sold", market: MarketYahooFlea, status: "sold", want: false},
		{name: "yahoo flea cancel", market: MarketYahooFlea, status: "cancel", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := IsActive(tt.market, tt.status, tt.auctionState); got != tt.want {
				t.Fatalf("IsActive(%q, %q, %q) = %v, want %v", tt.market, tt.status, tt.auctionState, got, tt.want)
			}
		})
	}
}
