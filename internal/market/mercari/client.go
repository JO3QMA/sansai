package mercari

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/jo3qma/sansai/internal/httpclient"
	"github.com/jo3qma/sansai/internal/model"
)

const (
	searchURL = "https://api.mercari.jp/v2/entities:search"
	itemURL   = "https://api.mercari.jp/items/get"
	itemBase  = "https://jp.mercari.com/item/"
)

type Client struct{}

type ecdsaSignature struct {
	R, S *big.Int
}

func (c *Client) Search(_ context.Context, query string, opts model.SearchOptions) (*model.SearchResult, error) {
	if opts.Limit <= 0 {
		opts.Limit = 20
	}
	if opts.Limit > 120 {
		opts.Limit = 120
	}

	body := map[string]any{
		"pageSize":           opts.Limit,
		"searchSessionId":    uuid.New().String(),
		"serviceFrom":        "suruga",
		"withItemBrand":      true,
		"withItemSize":       false,
		"withItemPromotions": true,
		"withItemSizes":      true,
		"withShopname":       false,
		"withAuction":        true,
		"searchCondition": map[string]any{
			"keyword":          query,
			"sort":             "SORT_SCORE",
			"order":            "ORDER_DESC",
			"status":           []string{"STATUS_ON_SALE"},
			"sizeId":           []any{},
			"categoryId":       []any{},
			"brandId":          []any{},
			"sellerId":         []any{},
			"priceMin":         opts.MinPrice,
			"priceMax":         opts.MaxPrice,
			"itemConditionId":  []any{},
			"shippingPayerId":  []any{},
			"shippingFromArea": []any{},
			"shippingMethod":   []any{},
			"colorId":          []any{},
			"hasCoupon":        false,
			"attributes":       []any{},
			"itemTypes":        []any{},
			"skuIds":           []any{},
			"excludeKeyword":   "",
		},
	}

	raw, err := c.postJSON(searchURL, body)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Items []struct {
			ID               string   `json:"id"`
			Name             string   `json:"name"`
			Price            string   `json:"price"`
			Status           string   `json:"status"`
			Thumbnails       []string `json:"thumbnails"`
			ItemConditionID  string   `json:"itemConditionId"`
			SellerID         string   `json:"sellerId"`
			Auction          *struct {
				BidDeadline string `json:"bidDeadline"`
				TotalBid    int    `json:"totalBid"`
			} `json:"auction"`
		} `json:"items"`
		Meta struct {
			NumFound string `json:"numFound"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("parse mercari search response: %w", err)
	}

	total, _ := strconv.Atoi(resp.Meta.NumFound)
	items := make([]model.Item, 0, len(resp.Items))
	for _, it := range resp.Items {
		price, _ := strconv.Atoi(it.Price)
		imageURL := ""
		if len(it.Thumbnails) > 0 {
			imageURL = it.Thumbnails[0]
		}
		item := model.Item{
			Market:    model.MarketMercari,
			ID:        it.ID,
			Title:     it.Name,
			Price:     price,
			Currency:  "JPY",
			URL:       itemBase + it.ID,
			ImageURL:  imageURL,
			ImageURLs: it.Thumbnails,
			Status:    it.Status,
			Condition: conditionLabel(it.ItemConditionID),
			Seller:    it.SellerID,
			SaleType:  model.SaleTypeFixedPrice,
		}
		if it.Auction != nil {
			item.SaleType = model.SaleTypeAuction
			if it.Auction.TotalBid > 0 {
				item.EndTime = it.Auction.BidDeadline
			}
		}
		items = append(items, item)
	}

	return &model.SearchResult{
		Market: model.MarketMercari,
		Query:  query,
		Items:  items,
		Total:  total,
		Page:   opts.Page,
	}, nil
}

func (c *Client) Get(_ context.Context, id string) (*model.Item, error) {
	target := itemURL + "?id=" + id + "&include_auction=true"
	req, err := http.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req, http.MethodGet, target)

	raw, status, err := httpclient.Do(req)
	if err != nil {
		return nil, err
	}
	if status == http.StatusNotFound {
		return nil, fmt.Errorf("item %s not found", id)
	}
	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("mercari item API HTTP %d", status)
	}

	var resp struct {
		Data mercariItemDetail `json:"data"`
	}
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("parse mercari item response: %w", err)
	}
	if resp.Data.ID == "" {
		return nil, fmt.Errorf("item %s not found", id)
	}

	return itemFromDetail(resp.Data), nil
}

type mercariItemDetail struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Price         int      `json:"price"`
	Status        string   `json:"status"`
	Description   string   `json:"description"`
	Photos        []string `json:"photos"`
	Thumbnails    []string `json:"thumbnails"`
	ItemCondition struct {
		Name string `json:"name"`
	} `json:"item_condition"`
	Seller struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"seller"`
	AuctionInfo *struct {
		ExpectedEndTime int64  `json:"expected_end_time"`
		TotalBids       int    `json:"total_bids"`
		State           string `json:"state"`
	} `json:"auction_info"`
}

func itemFromDetail(raw mercariItemDetail) *model.Item {
	imageURLs := raw.Photos
	if len(imageURLs) == 0 {
		imageURLs = raw.Thumbnails
	}
	imageURL := ""
	if len(imageURLs) > 0 {
		imageURL = imageURLs[0]
	}

	item := &model.Item{
		Market:      model.MarketMercari,
		ID:          raw.ID,
		Title:       raw.Name,
		Price:       raw.Price,
		Currency:    "JPY",
		URL:         itemBase + raw.ID,
		ImageURL:    imageURL,
		ImageURLs:   imageURLs,
		Description: strings.TrimSpace(raw.Description),
		SaleType:    model.SaleTypeFixedPrice,
		Status:      raw.Status,
		Condition:   raw.ItemCondition.Name,
		Seller:      strconv.FormatInt(raw.Seller.ID, 10),
	}

	if raw.AuctionInfo != nil {
		item.SaleType = model.SaleTypeAuction
		if raw.AuctionInfo.TotalBids > 0 {
			item.EndTime = unixRFC3339(raw.AuctionInfo.ExpectedEndTime)
		}
	}

	return item
}

func unixRFC3339(sec int64) string {
	if sec <= 0 {
		return ""
	}
	return time.Unix(sec, 0).UTC().Format(time.RFC3339)
}

func (c *Client) postJSON(url string, payload any) ([]byte, error) {
	buf, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	c.setHeaders(req, http.MethodPost, url)

	raw, status, err := httpclient.Do(req)
	if err != nil {
		return nil, err
	}
	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("mercari search API HTTP %d: %s", status, string(raw))
	}
	return raw, nil
}

func (c *Client) setHeaders(req *http.Request, method, url string) {
	req.Header.Set("DPoP", generateDPoP(method, url))
	req.Header.Set("X-Platform", "web")
	req.Header.Set("Accept-Language", "ja-JP,ja;q=0.9")
}

// ponytail: DPoP signing ported from goForMercari; breaks if Mercari changes auth.
func generateDPoP(method, url string) string {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	header := map[string]any{
		"typ": "dpop+jwt",
		"alg": "ES256",
		"jwk": map[string]string{
			"crv": "P-256",
			"kty": "EC",
			"x":   byteToBase64URL(privateKey.PublicKey.X.Bytes()),
			"y":   byteToBase64URL(privateKey.PublicKey.Y.Bytes()),
		},
	}
	payload := map[string]any{
		"iat": time.Now().Unix(),
		"jti": uuid.New().String(),
		"htu": url,
		"htm": strings.ToUpper(method),
	}

	headerJSON, _ := json.Marshal(header)
	payloadJSON, _ := json.Marshal(payload)
	unsigned := byteToBase64URL(headerJSON) + "." + byteToBase64URL(payloadJSON)

	hash := sha256.Sum256([]byte(unsigned))
	sigASN1, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	if err != nil {
		panic(err)
	}

	var sig ecdsaSignature
	if _, err := asn1.Unmarshal(sigASN1, &sig); err != nil {
		panic(err)
	}

	signature := byteToBase64URL(append(sig.R.Bytes(), sig.S.Bytes()...))
	return unsigned + "." + signature
}

func byteToBase64URL(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func conditionLabel(id string) string {
	switch id {
	case "1":
		return "新品、未使用"
	case "2":
		return "未使用に近い"
	case "3":
		return "目立った傷や汚れなし"
	case "4":
		return "やや傷や汚れあり"
	case "5":
		return "傷や汚れあり"
	case "6":
		return "全体的に状態が悪い"
	default:
		return id
	}
}
