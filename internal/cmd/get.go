package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jo3qma/sansai/internal/market"
	"github.com/jo3qma/sansai/internal/model"
)

var getCmd = &cobra.Command{
	Use:   "get <market> <id>",
	Short: "商品IDで詳細を取得",
	Long: `market: yahoo_auction | yahoo_flea | mercari
id: 各マーケットの商品ID`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		m, ok := model.ParseMarket(args[0])
		if !ok {
			return fmt.Errorf("unknown market: %s", args[0])
		}

		client, err := market.New(m)
		if err != nil {
			return err
		}

		item, err := client.Get(context.Background(), args[1])
		if err != nil {
			return err
		}
		return printJSON(item)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
