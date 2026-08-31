package yahooflea

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/jo3qma/sansai/internal/httpclient"
	"github.com/jo3qma/sansai/internal/model"
	"github.com/jo3qma/sansai/internal/nextdata"
)

var searchBase = "https://paypayfleamarket.yahoo.co.jp/search/"

const itemBase = "https://paypayfleamarket.yahoo.co.jp/item/"

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

	body, status, err := httpclient.GetStatus(searchURL)
	if err != nil {
		return nil, err
	}
	if status == http.StatusNotFound {
		return &model.SearchResult{
			Market: model.MarketYahooFlea,
			Query:  query,
			Items:  []model.Item{},
			Page:   opts.Page,
		}, nil
	}
	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("HTTP %d from %s", status, searchURL)
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
			IsActive:  model.IsActive(model.MarketYahooFlea, it.ItemStatus, ""),
			Condition: it.Condition,
			Seller:    it.SellerID,
			SaleType:  model.SaleTypeFixedPrice,
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
						Item fleaItemDetail `json:"item"`
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

	return itemFromDetail(raw), nil
}

type fleaItemDetail struct {
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
}

func itemFromDetail(raw fleaItemDetail) *model.Item {
	imageURLs := make([]string, 0, len(raw.Images))
	for _, img := range raw.Images {
		if img.URL != "" {
			imageURLs = append(imageURLs, img.URL)
		}
	}
	imageURL := ""
	if len(imageURLs) > 0 {
		imageURL = imageURLs[0]
	}

	return &model.Item{
		Market:      model.MarketYahooFlea,
		ID:          raw.ID,
		Title:       raw.Title,
		Price:       raw.Price,
		Currency:    "JPY",
		URL:         itemBase + raw.ID,
		ImageURL:    imageURL,
		ImageURLs:   imageURLs,
		Description: strings.TrimSpace(raw.Description),
		SaleType:    model.SaleTypeFixedPrice,
		Status:      raw.Status,
		IsActive:    model.IsActive(model.MarketYahooFlea, raw.Status, ""),
		Condition:   raw.Condition.Text,
		Extra: map[string]any{
			"like_count": raw.LikeCount,
			"brand":      raw.Brand.Name,
		},
	}
}
