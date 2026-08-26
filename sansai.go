package sansai

import (
	"github.com/jo3qma/sansai/internal/market"
	"github.com/jo3qma/sansai/internal/model"
)

// Market identifiers.
type Market = model.Market

const (
	MarketYahooAuction = model.MarketYahooAuction
	MarketYahooFlea    = model.MarketYahooFlea
	MarketMercari      = model.MarketMercari
)

type Item = model.Item
type SearchOptions = model.SearchOptions
type SearchResult = model.SearchResult
type SearchResponse = model.SearchResponse

// Client queries a single marketplace.
type Client = market.Client

// NewClient returns a client for the given market.
func NewClient(m Market) (Client, error) {
	return market.New(m)
}

// AllMarkets returns every supported market.
func AllMarkets() []Market {
	return model.AllMarkets()
}

// ParseMarket resolves a market name or alias.
func ParseMarket(s string) (Market, bool) {
	return model.ParseMarket(s)
}

// ParseMarkets parses a comma-separated market list ("all" = every market).
func ParseMarkets(raw string) ([]Market, error) {
	return market.ParseMarkets(raw)
}
