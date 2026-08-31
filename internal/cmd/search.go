package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/jo3qma/sansai"
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
		resp, err := sansai.Search(context.Background(), args[0], searchMarkets, sansai.SearchOptions{
			Limit:    searchLimit,
			Page:     searchPage,
			MinPrice: searchMin,
			MaxPrice: searchMax,
		})
		if resp != nil {
			if encErr := printJSON(resp); encErr != nil {
				return encErr
			}
		}
		return err
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
