package model

import (
	"strconv"
	"strings"
	"time"
)

// IsActive reports whether a listing can still be purchased or bid on.
// For Yahoo Auction, auctionState may carry end_time (RFC3339 or unix seconds) when status is empty.
func IsActive(market Market, status string, auctionState string) bool {
	switch market {
	case MarketMercari:
		return mercariIsActive(status, auctionState)
	case MarketYahooAuction:
		return yahooAuctionIsActive(status, auctionState)
	case MarketYahooFlea:
		return yahooFleaIsActive(status)
	default:
		return false
	}
}

func mercariIsActive(status, auctionState string) bool {
	if !mercariOnSale(status) {
		return false
	}
	switch auctionState {
	case "", "STATE_ONGOING", "STATE_NO_BID":
		return true
	default:
		return false
	}
}

func mercariOnSale(status string) bool {
	switch status {
	case "on_sale", "ITEM_STATUS_ON_SALE":
		return true
	default:
		return false
	}
}

func yahooAuctionIsActive(status, endTime string) bool {
	switch strings.ToLower(status) {
	case "open":
		return true
	case "closed", "sold", "cancel", "canceled":
		return false
	}
	if status != "" {
		return false
	}
	return endTimeAfterNow(endTime)
}

func endTimeAfterNow(endTime string) bool {
	if endTime == "" {
		return false
	}
	if unix, err := strconv.ParseInt(endTime, 10, 64); err == nil {
		return time.Unix(unix, 0).After(time.Now())
	}
	if t, err := time.Parse(time.RFC3339, endTime); err == nil {
		return t.After(time.Now())
	}
	return false
}

func yahooFleaIsActive(status string) bool {
	return strings.EqualFold(status, "open")
}
