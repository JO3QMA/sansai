---
name: sansai
description: >-
  LLM Agent 向けの日本 C2C マーケット（Yahoo!オークション・Yahoo!フリマ・メルカリ）商品検索 CLI。
  Use when searching or fetching items from mercari, Yahoo Auction, Yahoo Flea, PayPayフリマ, ヤフオク,
  or wiring Agent tool definitions for Japanese marketplace lookup.
---

# sansai

**sansai**（三サイト）は JSON 出力の CLI。Agent のツール呼び出しを主用途とする。

## インストール（Cursor Agent）

配布元はリポジトリの `skills/sansai/`。インストール先は **ユーザーの** `~/.cursor/skills/sansai/`（全プロジェクトで有効）。

```bash
mkdir -p ~/.cursor/skills/sansai
cp skills/sansai/SKILL.md ~/.cursor/skills/sansai/
```

リポジトリを clone していない場合:

```bash
mkdir -p ~/.cursor/skills/sansai
curl -fsSL https://raw.githubusercontent.com/JO3QMA/sansai/master/skills/sansai/SKILL.md \
  -o ~/.cursor/skills/sansai/SKILL.md
```

**注意:** リポジトリ内の `.cursor/skills/` は Cursor のプロジェクトスキル用。sansai CLI 利用用スキルは上記の個人スキルへ置く。

## 対応マーケット

| ID | サービス | エイリアス例 |
|---|---|---|
| `yahoo_auction` | Yahoo!オークション | `yahoo`, `auction`, `ヤフオク` |
| `yahoo_flea` | Yahoo!フリマ（PayPayフリマ） | `flea`, `paypay`, `フリマ` |
| `mercari` | メルカリ | `メルカリ` |

## CLI クイックリファレンス

```bash
# ビルド・インストール
go install github.com/jo3qma/sansai/cmd/sansai@latest

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
{
  "market": "mercari",
  "id": "m56797713000",
  "title": "商品名",
  "price": 1000,
  "currency": "JPY",
  "url": "https://...",
  "image_url": "https://...",
  "image_urls": ["https://...", "https://..."],
  "description": "出品説明文",
  "sale_type": "fixed_price",
  "end_time": "2026-08-28T20:00:00+09:00",
  "status": "on_sale",
  "condition": "新品",
  "seller": "出品者",
  "extra": {}
}
```

`sale_type` は `auction` または `fixed_price`。`end_time` は RFC3339。メルカリオークションで入札前は `end_time` を省略する場合がある。

## Agent 連携

1. `sansai schema` の出力をツール定義として登録する（`sansai_search` / `sansai_get`）。
2. 実呼び出しは CLI を subprocess で実行し stdout の JSON をパースする。

```json
{
  "command": "sansai",
  "args": ["search", "{{query}}", "-n", "5"]
}
```

検索 → 候補提示 → `get` で詳細取得、の流れが自然。

## 制約・注意

- 各サービスへの**非公式**アクセス。利用規約・robots.txt を確認し自己責任で利用。
- サイト側の HTML/API 変更で壊れる可能性がある。
- メルカリは DPoP 署名付き API（`goForMercari` 由来）。
- 過度なリクエストは避ける。Agent 利用時も件数・頻度に注意。
