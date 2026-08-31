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
					"description": "キーワードで複数マーケットを横断検索する。返却 Item の is_active は検索時点の推定値（get で確定）",
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
								"description": "取得件数 (マーケットごと)",
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
				},
				{
					"name":        "sansai_get",
					"description": "商品IDで詳細情報を取得する。返却 Item には description, image_urls, sale_type (auction|fixed_price), end_time (オークション・入札後), is_active (購入・入札可能か) を含む",
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
								"description": "商品ID",
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
