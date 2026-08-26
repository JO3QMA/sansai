---
name: sansai
description: >-
  LLM Agent 向けの日本 C2C マーケット（Yahoo!オークション・Yahoo!フリマ・メルカリ）商品検索 CLI。
  Use when searching or fetching items from mercari, Yahoo Auction, Yahoo Flea, PayPayフリマ, ヤフオク,
  developing sansai itself, or wiring Agent tool definitions for Japanese marketplace lookup.
---

# sansai

**sansai**（三サイト）は JSON 出力の CLI。Agent のツール呼び出しを主用途とする。

## 対応マーケット

| ID | サービス | エイリアス例 |
|---|---|---|
| `yahoo_auction` | Yahoo!オークション | `yahoo`, `auction`, `ヤフオク` |
| `yahoo_flea` | Yahoo!フリマ（PayPayフリマ） | `flea`, `paypay`, `フリマ` |
| `mercari` | メルカリ | `メルカリ` |

## CLI クイックリファレンス

```bash
# ビルド・インストール
go build -o bin/sansai ./cmd/sansai/
go install ./cmd/sansai

# 横断検索（全マーケット）
sansai search "ポケモンカード" -n 5

# 特定マーケット（カンマ区切り可）
sansai search "nintendo switch" -m mercari -n 10
sansai search "カメラ" -m yahoo_auction,yahoo_flea --min-price 10000 --max-price 50000

# 商品詳細
sansai get mercari m56797713000
sansai get yahoo_auction o1241906284
sansai get yahoo_flea z668531248

# マーケット一覧・ツール定義
sansai markets
sansai schema   # Function Calling 用 JSON
```

### search フラグ

| フラグ | 短縮 | 既定 | 説明 |
|---|---|---|---|
| `--market` | `-m` | `all` | 対象マーケット（`all` またはカンマ区切り） |
| `--limit` | `-n` | `10` | 取得件数（マーケットごと） |
| `--page` | | `1` | ページ番号 |
| `--min-price` | | `0` | 最低価格（円） |
| `--max-price` | | `0` | 最高価格（円） |

## 出力 JSON 型

```json
// SearchResponse（search）
{
  "query": "キーワード",
  "results": [
    {
      "market": "mercari",
      "query": "キーワード",
      "items": [
        {
          "market": "mercari",
          "id": "m56797713000",
          "title": "商品名",
          "price": 1000,
          "currency": "JPY",
          "url": "https://...",
          "image_url": "https://...",
          "status": "on_sale",
          "condition": "新品",
          "seller": "出品者",
          "extra": {}
        }
      ],
      "total": 100,
      "page": 1
    }
  ]
}

// Item（get）
{ "market": "mercari", "id": "...", "title": "...", "price": 1000, ... }
```

## Agent 連携

1. `sansai schema` の出力をツール定義として登録する（`sansai_search` / `sansai_get`）。
2. 実呼び出しは CLI を subprocess で実行し stdout の JSON をパースする。

```json
{
  "command": "sansai",
  "args": ["search", "{{query}}", "-n", "5"]
}
```

Agent が直接シェルを叩く場合は、検索→候補提示→`get` で詳細取得、の流れが自然。

## リポジトリ構造（sansai 本体を開発するとき）

```
cmd/sansai/          CLI エントリ
internal/cmd/        cobra サブコマンド
internal/market/     マーケット別クライアント（mercari, yahooauction, yahooflea）
internal/model/      共有型（Item, SearchResult 等）
internal/httpclient/ HTTP 共通
```

新マーケット追加: `internal/market/<name>/client.go` を実装し `registry.go` と `model.ParseMarket` に登録。

## 開発・検証

```bash
go test ./...
go vet ./...
go build -o /dev/null ./cmd/sansai

# 外部 API 実リクエスト（CI には含めない）
SANSAI_INTEGRATION=1 go test ./internal/market/mercari/ -run Integration
```

**CI gate**（`master` の PR/push）: `go test`, `go vet`, `go build ./cmd/sansai`。
**Integration test** は手動 workflow またはローカル専用。

用語は [CONTEXT.md](CONTEXT.md) を参照（Release / CI gate / Integration test 等）。

## 制約・注意

- 各サービスへの**非公式**アクセス。利用規約・robots.txt を確認し自己責任で利用。
- サイト側の HTML/API 変更で壊れる可能性がある。
- メルカリは DPoP 署名付き API（`goForMercari` 由来）。
- 過度なリクエストは避ける。Agent 利用時も件数・頻度に注意。

## 実装時の指針

- 最小差分。新マーケットは既存クライアント（`mercari`, `yahooauction`）のパターンに合わせる。
- CLI や出力形式を変えるときは README と `schema` コマンドを同期する。
- テストは HTTP モックで単体、実 API は Integration に分離。
