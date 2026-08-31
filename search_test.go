package sansai

import (
	"context"
	"errors"
	"testing"
)

type stubClient struct {
	search func(ctx context.Context, query string, opts SearchOptions) (*SearchResult, error)
	get    func(ctx context.Context, id string) (*Item, error)
}

func (s *stubClient) Search(ctx context.Context, query string, opts SearchOptions) (*SearchResult, error) {
	return s.search(ctx, query, opts)
}

func (s *stubClient) Get(ctx context.Context, id string) (*Item, error) {
	if s.get != nil {
		return s.get(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func TestSearchPartialFailure(t *testing.T) {
	t.Parallel()

	old := clientFactory
	t.Cleanup(func() { clientFactory = old })

	clientFactory = func(m Market) (Client, error) {
		switch m {
		case MarketMercari:
			return &stubClient{
				search: func(ctx context.Context, query string, opts SearchOptions) (*SearchResult, error) {
					return &SearchResult{
						Market: MarketMercari,
						Query:  query,
						Items:  []Item{{Market: MarketMercari, ID: "m1", Title: "ok"}},
					}, nil
				},
			}, nil
		case MarketYahooAuction:
			return &stubClient{
				search: func(ctx context.Context, query string, opts SearchOptions) (*SearchResult, error) {
					return nil, errors.New("HTTP 404 from auctions.yahoo.co.jp")
				},
			}, nil
		default:
			return nil, errors.New("unexpected market")
		}
	}

	resp, err := Search(context.Background(), "test", "mercari,yahoo_auction", SearchOptions{Limit: 5})
	if err != nil {
		t.Fatalf("expected partial success, got err %v", err)
	}
	if len(resp.Results) != 2 {
		t.Fatalf("got %d results", len(resp.Results))
	}

	var mercari, auction SearchResult
	for _, r := range resp.Results {
		switch r.Market {
		case MarketMercari:
			mercari = r
		case MarketYahooAuction:
			auction = r
		}
	}
	if len(mercari.Items) != 1 || mercari.Error != "" {
		t.Fatalf("mercari: %#v", mercari)
	}
	if auction.Error == "" || len(auction.Items) != 0 {
		t.Fatalf("yahoo_auction: %#v", auction)
	}
}

func TestSearchAllMarketsFailed(t *testing.T) {
	t.Parallel()

	old := clientFactory
	t.Cleanup(func() { clientFactory = old })

	clientFactory = func(m Market) (Client, error) {
		return &stubClient{
			search: func(ctx context.Context, query string, opts SearchOptions) (*SearchResult, error) {
				return nil, errors.New("down")
			},
		}, nil
	}

	resp, err := Search(context.Background(), "test", "mercari", SearchOptions{})
	if !errors.Is(err, ErrAllMarketsFailed) {
		t.Fatalf("expected ErrAllMarketsFailed, got %v", err)
	}
	if len(resp.Results) != 1 || resp.Results[0].Error != "down" {
		t.Fatalf("resp: %#v", resp)
	}
}
