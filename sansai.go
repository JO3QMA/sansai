package sansai

import (
	"context"

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
type Client interface {
	Search(ctx context.Context, query string, opts SearchOptions) (*SearchResult, error)
	Get(ctx context.Context, id string) (*Item, error)
}

type client struct {
	inner market.Client
}

func (c *client) Search(ctx context.Context, query string, opts SearchOptions) (*SearchResult, error) {
	return c.inner.Search(ctx, query, opts)
}

func (c *client) Get(ctx context.Context, id string) (*Item, error) {
	return c.inner.Get(ctx, id)
}

// NewClient returns a client for the given market.
func NewClient(m Market) (Client, error) {
	inner, err := market.New(m)
	if err != nil {
		return nil, err
	}
	return &client{inner: inner}, nil
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
