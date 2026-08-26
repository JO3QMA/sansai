# sansai

LLM Agent 向けの日本 C2C マーケット商品検索 CLI。

## Distribution

**Release**:
GitHub Release として公開する、バージョン付きの prebuilt `sansai` バイナリ一式。GoReleaser 成功後は Draft にせず即 Published。
_Avoid_: Deploy, publish（バイナリ配布の意味で使う場合）

**Release tag**:
`vMAJOR.MINOR.PATCH` 形式の git tag。push されると Release ワークフローが起動する。
_Avoid_: Version tag（git describe 等と紛らわしい）

**Release notes**:
GitHub Release 本文。前回 Release tag から今回 tag までの git commit 一覧を GoReleaser が自動生成する。
_Avoid_: Changelog（ファイルとして維持する意味ではない）

**Initial Release**:
CI/CD 導入後の最初の Release tag。`v0.1.0` とする。
_Avoid_: First tag, launch version

**Release target**:
Release バイナリを build する OS/アーキテクチャの組。初版は linux（amd64, arm64）、darwin（amd64, arm64）、windows（amd64）の 5 種。
_Avoid_: Platform, triplet

**CI/CD**:
PR・push 時の品質ゲート（CI）と、Release tag push 時の Release 配布（CD）を指す。
_Avoid_: Pipeline（このリポジトリでは CI/CD と呼ぶ）

**CI gate**:
Release tag 以外の push / PR で必ず通過させるチェック。初版は `go test ./...`、`go vet ./...`、`go build ./cmd/sansai`。`master` への merge 前に **branch protection で必須**とする。
_Avoid_: Pre-release check, quality gate

**Release gate**:
Release tag push 時、GoReleaser の前に CI gate と同じチェックを再実行する最終防衛線。
_Avoid_: Pre-release gate, smoke test（この文脈では Release gate と呼ぶ）

**Integration test**:
Mercari 等の外部 API に実リクエストするテスト。CI gate・Release gate には含めず、手動 workflow またはローカル（`SANSAI_INTEGRATION=1`）でのみ実行する。
_Avoid_: E2E test, live test

**Dependabot**:
`go.mod` と GitHub Actions workflow の依存更新 PR を自動作成する仕組み。監視対象は **gomod**（patch/minor のみ）と **github-actions**。チェック頻度は **週次**。新バージョン公開から **7 日間 cooldown** 経過後に PR を作成する（両 ecosystem 共通。セキュリティ更新は除外）。更新は **gomod 1 PR・github-actions 1 PR** にグループ化する。CI gate 通過後は **自動 merge** する。
_Avoid_: Renovate（このリポジトリでは Dependabot を使う）
