package yahooflea

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/jo3qma/sansai/internal/httpclient"
	"github.com/jo3qma/sansai/internal/model"
	"github.com/jo3qma/sansai/internal/nextdata"
)

const (
	searchBase = "https://paypayfleamarket.yahoo.co.jp/search/"
	itemBase   = "https://paypayfleamarket.yahoo.co.jp/item/"
)

type Client struct{}

func (c *Client) Search(_ context.Context, query string, opts model.SearchOptions) (*model.SearchResult, error) {
	if opts.Limit <= 0 {
		opts.Limit = 20
	}
	if opts.Page <= 0 {
		opts.Page = 1
	}

	searchURL := searchBase + url.PathEscape(query)
	if opts.Page > 1 {
		searchURL += fmt.Sprintf("?page=%d", opts.Page)
	}

	body, err := httpclient.Get(searchURL)
	if err != nil {
		return nil, err
	}

	var data struct {
		Props struct {
			InitialState struct {
				SearchState struct {
					Search struct {
						Result struct {
							TotalResultsAvailable int `json:"totalResultsAvailable"`
							Items                 []struct {
								ID                string `json:"id"`
								Title             string `json:"title"`
								Price             int    `json:"price"`
								ItemStatus        string `json:"itemStatus"`
								ThumbnailImageURL string `json:"thumbnailImageUrl"`
								Condition         string `json:"condition"`
								SellerID          string `json:"sellerId"`
							} `json:"items"`
						} `json:"result"`
					} `json:"search"`
				} `json:"searchState"`
			} `json:"initialState"`
		} `json:"props"`
	}
	if err := nextdata.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	raw := data.Props.InitialState.SearchState.Search.Result
	items := make([]model.Item, 0, len(raw.Items))
	for _, it := range raw.Items {
		if opts.MinPrice > 0 && it.Price < opts.MinPrice {
			continue
		}
		if opts.MaxPrice > 0 && it.Price > opts.MaxPrice {
			continue
		}
		items = append(items, model.Item{
			Market:    model.MarketYahooFlea,
			ID:        it.ID,
			Title:     it.Title,
			Price:     it.Price,
			Currency:  "JPY",
			URL:       itemBase + it.ID,
			ImageURL:  it.ThumbnailImageURL,
			Status:    it.ItemStatus,
			Condition: it.Condition,
			Seller:    it.SellerID,
		})
		if len(items) >= opts.Limit {
			break
		}
	}

	return &model.SearchResult{
		Market: model.MarketYahooFlea,
		Query:  query,
		Items:  items,
		Total:  raw.TotalResultsAvailable,
		Page:   opts.Page,
	}, nil
}

func (c *Client) Get(_ context.Context, id string) (*model.Item, error) {
	body, err := httpclient.Get(itemBase + url.PathEscape(id))
	if err != nil {
		return nil, err
	}

	var data struct {
		Props struct {
			InitialState struct {
				ItemsState struct {
					Items struct {
						Item struct {
							ID          string `json:"id"`
							Title       string `json:"title"`
							Description string `json:"description"`
							Price       int    `json:"price"`
							Status      string `json:"status"`
							LikeCount   int    `json:"likeCount"`
							Images      []struct {
								URL string `json:"url"`
							} `json:"images"`
							Condition struct {
								Text string `json:"text"`
							} `json:"condition"`
							Brand struct {
								Name string `json:"name"`
							} `json:"brand"`
						} `json:"item"`
					} `json:"items"`
				} `json:"itemsState"`
			} `json:"initialState"`
		} `json:"props"`
	}
	if err := nextdata.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	raw := data.Props.InitialState.ItemsState.Items.Item
	if raw.ID == "" {
		return nil, fmt.Errorf("item %s not found", id)
	}

	imageURL := ""
	if len(raw.Images) > 0 {
		imageURL = raw.Images[0].URL
	}

	return &model.Item{
		Market:    model.MarketYahooFlea,
		ID:        raw.ID,
		Title:     raw.Title,
		Price:     raw.Price,
		Currency:  "JPY",
		URL:       itemBase + raw.ID,
		ImageURL:  imageURL,
		Status:    raw.Status,
		Condition: raw.Condition.Text,
		Extra: map[string]any{
			"description": strings.TrimSpace(raw.Description),
			"like_count":  raw.LikeCount,
			"brand":       raw.Brand.Name,
		},
	}, nil
}
