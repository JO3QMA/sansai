package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yafuoku",
	Short: "LLM Agent向け 日本C2Cマーケット商品検索CLI",
	Long: `yafuoku は Yahoo!オークション・Yahoo!フリマ・メルカリの商品情報を
JSON形式で取得するCLIツールです。LLM Agentからのツール呼び出しを想定しています。`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(v)
}
