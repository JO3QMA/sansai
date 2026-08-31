package sansai

import (
	"context"
	"errors"
)

var clientFactory = NewClient

// ErrAllMarketsFailed is returned when every requested market search fails.
var ErrAllMarketsFailed = errors.New("all markets failed")

// Search queries one or more marketplaces and returns aggregated results.
// marketSpec is "all", a single market id, or a comma-separated list.
// Per-market failures are recorded in SearchResult.Error; other markets still run.
func Search(ctx context.Context, query string, marketSpec string, opts SearchOptions) (*SearchResponse, error) {
	markets, err := ParseMarkets(marketSpec)
	if err != nil {
		return nil, err
	}

	resp := &SearchResponse{Query: query}
	var hasSuccess bool
	for _, m := range markets {
		result, ok := searchMarket(ctx, m, query, opts)
		resp.Results = append(resp.Results, result)
		if ok {
			hasSuccess = true
		}
	}
	if !hasSuccess {
		return resp, ErrAllMarketsFailed
	}
	return resp, nil
}

func searchMarket(ctx context.Context, m Market, query string, opts SearchOptions) (SearchResult, bool) {
	client, err := clientFactory(m)
	if err != nil {
		return SearchResult{
			Market: m,
			Query:  query,
			Items:  []Item{},
			Error:  err.Error(),
		}, false
	}

	result, err := client.Search(ctx, query, opts)
	if err != nil {
		return SearchResult{
			Market: m,
			Query:  query,
			Items:  []Item{},
			Error:  err.Error(),
		}, false
	}
	return *result, true
}
