package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jo3qma/sansai/internal/market"
	"github.com/jo3qma/sansai/internal/model"
)

var (
	searchMarkets string
	searchLimit   int
	searchPage    int
	searchMin     int
	searchMax     int
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "キーワードで商品を検索",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		markets, err := market.ParseMarkets(searchMarkets)
		if err != nil {
			return err
		}

		opts := model.SearchOptions{
			Limit:    searchLimit,
			Page:     searchPage,
			MinPrice: searchMin,
			MaxPrice: searchMax,
		}

		resp := model.SearchResponse{Query: args[0]}
		ctx := context.Background()

		for _, m := range markets {
			client, err := market.New(m)
			if err != nil {
				return err
			}
			result, err := client.Search(ctx, args[0], opts)
			if err != nil {
				return fmt.Errorf("%s: %w", m, err)
			}
			resp.Results = append(resp.Results, *result)
		}

		return printJSON(resp)
	},
}

func init() {
	searchCmd.Flags().StringVarP(&searchMarkets, "market", "m", "all", "対象マーケット (yahoo_auction,yahoo_flea,mercari,all)")
	searchCmd.Flags().IntVarP(&searchLimit, "limit", "n", 10, "取得件数 (マーケットごと)")
	searchCmd.Flags().IntVar(&searchPage, "page", 1, "ページ番号")
	searchCmd.Flags().IntVar(&searchMin, "min-price", 0, "最低価格 (円)")
	searchCmd.Flags().IntVar(&searchMax, "max-price", 0, "最高価格 (円)")
	rootCmd.AddCommand(searchCmd)
}
