package sansai

import (
	"context"
	"fmt"
)

// Search queries one or more marketplaces and returns aggregated results.
// marketSpec is "all", a single market id, or a comma-separated list.
func Search(ctx context.Context, query string, marketSpec string, opts SearchOptions) (*SearchResponse, error) {
	markets, err := ParseMarkets(marketSpec)
	if err != nil {
		return nil, err
	}

	resp := &SearchResponse{Query: query}
	for _, m := range markets {
		client, err := NewClient(m)
		if err != nil {
			return nil, err
		}
		result, err := client.Search(ctx, query, opts)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", m, err)
		}
		resp.Results = append(resp.Results, *result)
	}
	return resp, nil
}
