package yahooauction

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jo3qma/sansai/internal/httpclient"
	"github.com/jo3qma/sansai/internal/model"
	"github.com/jo3qma/sansai/internal/nextdata"
)

const (
	searchBase = "https://auctions.yahoo.co.jp/search/search"
	itemBase   = "https://auctions.yahoo.co.jp/jp/auction/"
)

var productBlockRe = regexp.MustCompile(`(?s)<li class="Product">.*?</li>`)

type Client struct{}

func (c *Client) Search(_ context.Context, query string, opts model.SearchOptions) (*model.SearchResult, error) {
	if opts.Limit <= 0 {
		opts.Limit = 20
	}
	if opts.Page <= 0 {
		opts.Page = 1
	}

	offset := (opts.Page-1)*opts.Limit + 1
	params := url.Values{
		"p":  {query},
		"va": {query},
		"b":  {strconv.Itoa(offset)},
		"n":  {strconv.Itoa(opts.Limit)},
	}
	if opts.MinPrice > 0 {
		params.Set("min", strconv.Itoa(opts.MinPrice))
	}
	if opts.MaxPrice > 0 {
		params.Set("max", strconv.Itoa(opts.MaxPrice))
	}

	body, err := httpclient.Get(searchBase + "?" + params.Encode())
	if err != nil {
		return nil, err
	}

	items := parseSearchHTML(string(body), opts.MinPrice, opts.MaxPrice)
	if len(items) > opts.Limit {
		items = items[:opts.Limit]
	}

	return &model.SearchResult{
		Market: model.MarketYahooAuction,
		Query:  query,
		Items:  items,
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
			PageProps struct {
				InitialState struct {
					Item struct {
						Detail struct {
							Item auctionItemDetail `json:"item"`
						} `json:"detail"`
					} `json:"item"`
				} `json:"initialState"`
			} `json:"pageProps"`
		} `json:"props"`
	}
	if err := nextdata.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	raw := data.Props.PageProps.InitialState.Item.Detail.Item
	if raw.AuctionID == "" {
		return nil, fmt.Errorf("item %s not found", id)
	}

	return itemFromDetail(raw), nil
}

type auctionItemDetail struct {
	AuctionID     string   `json:"auctionId"`
	Title         string   `json:"title"`
	Price         int      `json:"price"`
	BuyNowPrice   int      `json:"buyNowPrice"`
	EndTime       string   `json:"endTime"`
	ItemStatus    string   `json:"itemStatus"`
	ConditionName string   `json:"conditionName"`
	Description   []string `json:"description"`
	Img           []struct {
		Image string `json:"image"`
	} `json:"img"`
	Images []struct {
		Image string `json:"image"`
	} `json:"images"`
	Seller struct {
		AucUserID string `json:"aucUserId"`
	} `json:"seller"`
}

func itemFromDetail(raw auctionItemDetail) *model.Item {
	imageURLs := imageURLsFromAuction(raw.Img, raw.Images)
	imageURL := ""
	if len(imageURLs) > 0 {
		imageURL = imageURLs[0]
	}

	item := &model.Item{
		Market:      model.MarketYahooAuction,
		ID:          raw.AuctionID,
		Title:       raw.Title,
		Price:       raw.Price,
		Currency:    "JPY",
		URL:         itemBase + raw.AuctionID,
		ImageURL:    imageURL,
		ImageURLs:   imageURLs,
		Description: joinDescription(raw.Description),
		SaleType:    model.SaleTypeAuction,
		Status:      raw.ItemStatus,
		Condition:   raw.ConditionName,
		Seller:      raw.Seller.AucUserID,
		Extra: map[string]any{
			"buy_now_price": raw.BuyNowPrice,
		},
	}
	if raw.EndTime != "" {
		item.EndTime = raw.EndTime
	}
	item.IsActive = model.IsActive(model.MarketYahooAuction, raw.ItemStatus, raw.EndTime)
	return item
}

func imageURLsFromAuction(img, images []struct {
	Image string `json:"image"`
}) []string {
	source := img
	if len(source) == 0 {
		source = images
	}
	urls := make([]string, 0, len(source))
	for _, entry := range source {
		if entry.Image != "" {
			urls = append(urls, entry.Image)
		}
	}
	return urls
}

func joinDescription(parts []string) string {
	nonEmpty := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			nonEmpty = append(nonEmpty, part)
		}
	}
	return strings.Join(nonEmpty, "\n")
}

func parseSearchHTML(html string, minPrice, maxPrice int) []model.Item {
	blocks := productBlockRe.FindAllString(html, -1)
	items := make([]model.Item, 0, len(blocks))

	attr := func(block, name string) string {
		re := regexp.MustCompile(name + `="([^"]*)"`)
		m := re.FindStringSubmatch(block)
		if len(m) < 2 {
			return ""
		}
		return m[1]
	}

	for _, block := range blocks {
		id := attr(block, `data-auction-id`)
		if id == "" {
			continue
		}
		price, _ := strconv.Atoi(attr(block, `data-auction-price`))
		if minPrice > 0 && price > 0 && price < minPrice {
			continue
		}
		if maxPrice > 0 && price > maxPrice {
			continue
		}
		endUnix := attr(block, `data-auction-endtime`)
		endTime := unixToRFC3339(endUnix)
		status := yahooSearchStatus(endUnix)
		items = append(items, model.Item{
			Market:   model.MarketYahooAuction,
			ID:       id,
			Title:    attr(block, `data-auction-title`),
			Price:    price,
			Currency: "JPY",
			URL:      itemBase + id,
			ImageURL: attr(block, `data-auction-img`),
			SaleType: model.SaleTypeAuction,
			EndTime:  endTime,
			Status:   status,
			IsActive: model.IsActive(model.MarketYahooAuction, status, endUnix),
		})
	}
	return items
}

func yahooSearchStatus(endUnix string) string {
	if endUnix == "" {
		return ""
	}
	if model.IsActive(model.MarketYahooAuction, "", endUnix) {
		return "open"
	}
	return "closed"
}

func unixToRFC3339(unix string) string {
	sec, err := strconv.ParseInt(unix, 10, 64)
	if err != nil || sec <= 0 {
		return ""
	}
	return time.Unix(sec, 0).Format(time.RFC3339)
}
