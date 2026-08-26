package market

import (
	"context"
	"fmt"

	"github.com/jo3qma/yafuoku/internal/market/mercari"
	"github.com/jo3qma/yafuoku/internal/market/yahooauction"
	"github.com/jo3qma/yafuoku/internal/market/yahooflea"
	"github.com/jo3qma/yafuoku/internal/model"
)

type Client interface {
	Search(ctx context.Context, query string, opts model.SearchOptions) (*model.SearchResult, error)
	Get(ctx context.Context, id string) (*model.Item, error)
}

func New(m model.Market) (Client, error) {
	switch m {
	case model.MarketYahooAuction:
		return &yahooauction.Client{}, nil
	case model.MarketYahooFlea:
		return &yahooflea.Client{}, nil
	case model.MarketMercari:
		return &mercari.Client{}, nil
	default:
		return nil, fmt.Errorf("unknown market: %s", m)
	}
}

func ParseMarkets(raw string) ([]model.Market, error) {
	if raw == "" || raw == "all" {
		return model.AllMarkets(), nil
	}

	seen := map[model.Market]bool{}
	var out []model.Market
	for _, part := range splitCSV(raw) {
		m, ok := model.ParseMarket(part)
		if !ok {
			return nil, fmt.Errorf("unknown market: %s", part)
		}
		if !seen[m] {
			seen[m] = true
			out = append(out, m)
		}
	}
	return out, nil
}

func splitCSV(s string) []string {
	var parts []string
	for _, p := range splitOn(s, ',') {
		if p != "" {
			parts = append(parts, p)
		}
	}
	return parts
}

func splitOn(s string, sep byte) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == sep {
			out = append(out, trim(s[start:i]))
			start = i + 1
		}
	}
	out = append(out, trim(s[start:]))
	return out
}

func trim(s string) string {
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\t') {
		s = s[:len(s)-1]
	}
	return s
}
