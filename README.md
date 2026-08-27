# sansai

**sansai**（三サイト）は、LLM Agent 向けの日本 C2C マーケット商品検索 CLI です。

対応マーケット:

| ID | サービス |
|---|---|
| `yahoo_auction` | Yahoo!オークション |
| `yahoo_flea` | Yahoo!フリマ（PayPayフリマ） |
| `mercari` | メルカリ |

出力はすべて JSON です。Agent からのツール呼び出しを想定しています。

`get` で返る `Item` には Bot 連携向けの正規フィールドがあります:

| フィールド | 型 | 説明 |
|---|---|---|
| `description` | string | 出品説明文 |
| `image_urls` | []string | 画像 URL 一覧 |
| `sale_type` | `"auction"` \| `"fixed_price"` | 販売形式 |
| `end_time` | string (RFC3339) | オークション終了時刻（該当時のみ） |
| `is_active` | bool | 購入・入札可能なら `true`。売切れ・終了・取消等は `false` |

## インストール

```bash
go install ./cmd/sansai
# または
go build -o bin/sansai ./cmd/sansai/
```

## 使い方

### 検索

```bash
# 全マーケット横断検索
sansai search "ポケモンカード" -n 5

# 特定マーケットのみ
sansai search "nintendo switch" -m mercari -n 10
sansai search "nintendo switch" -m yahoo_auction,yahoo_flea -n 5

# 価格フィルタ
sansai search "カメラ" --min-price 10000 --max-price 50000
```

### 商品詳細

```bash
sansai get mercari m56797713000
sansai get yahoo_auction o1241906284
sansai get yahoo_flea z668531248
```

### Agent 向け

```bash
# 対応マーケット一覧
sansai markets

# ツール定義 (Function Calling 用)
sansai schema
```

## Go ライブラリとして使う

他の Go ツールからは `github.com/jo3qma/sansai` を import してください。

```go
import (
    "context"
    "github.com/jo3qma/sansai"
)

func example() error {
    // 横断検索
    resp, err := sansai.Search(context.Background(), "ポケモンカード", "all", sansai.SearchOptions{
        Limit: 5,
    })
    if err != nil {
        return err
    }
    _ = resp

    // 単品取得
    item, err := sansai.Get(context.Background(), sansai.MarketMercari, "m56797713000")
    if err != nil {
        return err
    }
    _ = item

    // 単一マーケット用クライアント
    client, err := sansai.NewClient(sansai.MarketYahooAuction)
    if err != nil {
        return err
    }
    result, err := client.Search(context.Background(), "nintendo switch", sansai.SearchOptions{Limit: 10})
    _ = result
    return err
}
```

実装の詳細は `internal/` に閉じ込めてあり、公開 API はルートの `sansai` パッケージのみです。

## Agent 連携例

### Cursor Agent スキル

CLI 利用用スキルは `skills/sansai/SKILL.md` にあります。Cursor が全プロジェクトで読む場所へコピーしてください:

```bash
mkdir -p ~/.cursor/skills/sansai
cp skills/sansai/SKILL.md ~/.cursor/skills/sansai/
```

### ツール登録

Cursor / Claude などの Agent に次のように登録します:

```json
{
  "command": "sansai",
  "args": ["search", "{{query}}", "-n", "5"]
}
```

`sansai schema` の出力をそのままツール定義として使えます。

## 注意事項

- 各サービスの**非公式**アクセスです。利用規約・robots.txt を確認のうえ自己責任でご利用ください。
- サイト側の HTML / API 変更で動かなくなる可能性があります。
- メルカリは DPoP 署名付き API を使用しています（`goForMercari` 由来の実装）。
- 過度なリクエストは避けてください。

## 開発

```bash
go test ./...
SANSAI_INTEGRATION=1 go test ./internal/market/mercari/ -run Integration
```

## ライセンス

MIT
