package sansai

import "context"

// Get fetches a single item by market and id.
func Get(ctx context.Context, m Market, id string) (*Item, error) {
	client, err := NewClient(m)
	if err != nil {
		return nil, err
	}
	return client.Get(ctx, id)
}
