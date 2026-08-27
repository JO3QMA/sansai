package model

import "strings"

// IsActive reports whether a listing can still be purchased or bid on.
func IsActive(market Market, status string, auctionState string) bool {
	switch market {
	case MarketMercari:
		return mercariIsActive(status, auctionState)
	case MarketYahooAuction:
		return yahooAuctionIsActive(status)
	case MarketYahooFlea:
		return yahooFleaIsActive(status)
	default:
		return false
	}
}

func mercariIsActive(status, auctionState string) bool {
	if status != "on_sale" {
		return false
	}
	switch auctionState {
	case "", "STATE_ONGOING", "STATE_NO_BID":
		return true
	default:
		return false
	}
}

func yahooAuctionIsActive(status string) bool {
	return strings.EqualFold(status, "open")
}

func yahooFleaIsActive(status string) bool {
	return strings.EqualFold(status, "open")
}
