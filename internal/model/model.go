package model

type Market string

const (
	MarketYahooAuction Market = "yahoo_auction"
	MarketYahooFlea    Market = "yahoo_flea"
	MarketMercari      Market = "mercari"
)

func AllMarkets() []Market {
	return []Market{MarketYahooAuction, MarketYahooFlea, MarketMercari}
}

func ParseMarket(s string) (Market, bool) {
	switch s {
	case string(MarketYahooAuction), "yahoo", "auction", "ヤフオク":
		return MarketYahooAuction, true
	case string(MarketYahooFlea), "flea", "paypay", "フリマ":
		return MarketYahooFlea, true
	case string(MarketMercari), "メルカリ":
		return MarketMercari, true
	default:
		return "", false
	}
}

type SaleType string

const (
	SaleTypeAuction    SaleType = "auction"
	SaleTypeFixedPrice SaleType = "fixed_price"
)

type Item struct {
	Market      Market   `json:"market"`
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Price       int      `json:"price"`
	Currency    string   `json:"currency"`
	URL         string   `json:"url"`
	ImageURL    string   `json:"image_url,omitempty"`
	ImageURLs   []string `json:"image_urls,omitempty"`
	Description string   `json:"description,omitempty"`
	SaleType    SaleType `json:"sale_type,omitempty"`
	EndTime     string   `json:"end_time,omitempty"`
	Status      string   `json:"status,omitempty"`
	Condition   string   `json:"condition,omitempty"`
	Seller      string   `json:"seller,omitempty"`
	Extra       any      `json:"extra,omitempty"`
}

type SearchOptions struct {
	Limit    int
	Page     int
	MinPrice int
	MaxPrice int
}

type SearchResult struct {
	Market Market `json:"market"`
	Query  string `json:"query"`
	Items  []Item `json:"items"`
	Total  int    `json:"total,omitempty"`
	Page   int    `json:"page,omitempty"`
}

type SearchResponse struct {
	Query   string         `json:"query"`
	Results []SearchResult `json:"results"`
}
