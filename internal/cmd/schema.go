package cmd

import (
	"github.com/spf13/cobra"
)

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "LLM Agent向けツール定義をJSONで出力",
	RunE: func(cmd *cobra.Command, args []string) error {
		schema := map[string]any{
			"name":        "sansai",
			"description": "日本のC2Cマーケット (Yahoo!オークション, Yahoo!フリマ, メルカリ) から商品情報を取得する",
			"tools": []map[string]any{
				{
					"name":        "sansai_search",
					"description": "キーワードで複数マーケットを横断検索する。各 results[] に market ごとの items が入る。1市場が失敗しても他市場は続行し JSON は返す（全市場失敗時のみ非0終了）。results[].error に市場ローカルな失敗理由が入ることがある（404等の0件は空 items で error なし）。search の is_active は検索時点の推定値で、市場によって get と揃いやすい。yahoo_auction の -n は内部で 20/50/100 にクランプされる",
					"parameters": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"query": map[string]any{
								"type":        "string",
								"description": "検索キーワード",
							},
							"market": map[string]any{
								"type":        "string",
								"description": "対象マーケット (yahoo_auction, yahoo_flea, mercari, all)",
								"default":     "all",
							},
							"limit": map[string]any{
								"type":        "integer",
								"description": "取得件数 (マーケットごと)。yahoo_auction は内部で 20/50/100 のページサイズに丸める",
								"default":     10,
							},
							"min_price": map[string]any{
								"type":        "integer",
								"description": "最低価格 (円)",
							},
							"max_price": map[string]any{
								"type":        "integer",
								"description": "最高価格 (円)",
							},
						},
						"required": []string{"query"},
					},
					"returns": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"query": map[string]any{"type": "string"},
							"results": map[string]any{
								"type": "array",
								"items": map[string]any{
									"type": "object",
									"properties": map[string]any{
										"market": map[string]any{"type": "string"},
										"items":  map[string]any{"type": "array"},
										"error": map[string]any{
											"type":        "string",
											"description": "市場ローカルな失敗。空なら成功または0件",
										},
									},
								},
							},
						},
					},
				},
				{
					"name":        "sansai_get",
					"description": "商品IDで詳細情報を取得する。返却 Item には description, image_urls, sale_type (auction|fixed_price), end_time (オークション・入札後), is_active (購入・入札可能か) を含む。mercari は個人出品 (id が m+数字) のみ対応。search で extra.shops=true の Shops 出品は get 不可（unsupported id エラー）— search 結果の title/price/url を使う",
					"parameters": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"market": map[string]any{
								"type":        "string",
								"enum":        []string{"yahoo_auction", "yahoo_flea", "mercari"},
								"description": "マーケット名",
							},
							"id": map[string]any{
								"type":        "string",
								"description": "商品ID。mercari は m+数字のみ（例: m29820723477）。Shops 出品IDは get 不可",
							},
						},
						"required": []string{"market", "id"},
					},
				},
			},
		}
		return printJSON(schema)
	},
}

func init() {
	rootCmd.AddCommand(schemaCmd)
}
