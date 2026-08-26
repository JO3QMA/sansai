package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jo3qma/yafuoku/internal/model"
)

var marketsCmd = &cobra.Command{
	Use:   "markets",
	Short: "対応マーケット一覧を表示",
	RunE: func(cmd *cobra.Command, args []string) error {
		type entry struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		return printJSON([]entry{
			{ID: string(model.MarketYahooAuction), Name: "Yahoo!オークション", Description: "日本最大級のオークションサイト"},
			{ID: string(model.MarketYahooFlea), Name: "Yahoo!フリマ", Description: "PayPayフリマ (旧PayPayフリマ)"},
			{ID: string(model.MarketMercari), Name: "メルカリ", Description: "日本最大級のフリマアプリ"},
		})
	},
}

func init() {
	rootCmd.AddCommand(marketsCmd)
}
