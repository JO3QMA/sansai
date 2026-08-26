# yafuoku

**yafuoku**（ヤフオク + フリマ + メルカリ）は、LLM Agent 向けの日本 C2C マーケット商品検索 CLI です。

対応マーケット:

| ID | サービス |
|---|---|
| `yahoo_auction` | Yahoo!オークション |
| `yahoo_flea` | Yahoo!フリマ（PayPayフリマ） |
| `mercari` | メルカリ |

出力はすべて JSON です。Agent からのツール呼び出しを想定しています。

## インストール

```bash
go install ./cmd/yafuoku
# または
go build -o bin/yafuoku ./cmd/yafuoku/
```

## 使い方

### 検索

```bash
# 全マーケット横断検索
yafuoku search "ポケモンカード" -n 5

# 特定マーケットのみ
yafuoku search "nintendo switch" -m mercari -n 10
yafuoku search "nintendo switch" -m yahoo_auction,yahoo_flea -n 5

# 価格フィルタ
yafuoku search "カメラ" --min-price 10000 --max-price 50000
```

### 商品詳細

```bash
yafuoku get mercari m56797713000
yafuoku get yahoo_auction o1241906284
yafuoku get yahoo_flea z668531248
```

### Agent 向け

```bash
# 対応マーケット一覧
yafuoku markets

# ツール定義 (Function Calling 用)
yafuoku schema
```

## Agent 連携例

Cursor / Claude などの Agent に次のように登録します:

```json
{
  "command": "yafuoku",
  "args": ["search", "{{query}}", "-n", "5"]
}
```

`yafuoku schema` の出力をそのままツール定義として使えます。

## 注意事項

- 各サービスの**非公式**アクセスです。利用規約・robots.txt を確認のうえ自己責任でご利用ください。
- サイト側の HTML / API 変更で動かなくなる可能性があります。
- メルカリは DPoP 署名付き API を使用しています（`goForMercari` 由来の実装）。
- 過度なリクエストは避けてください。

## 開発

```bash
go test ./...
YAFUOKU_INTEGRATION=1 go test ./internal/market/mercari/ -run Integration
```

## ライセンス

MIT
