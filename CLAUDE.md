# cc-plans

Claude Code plans CLI tool (Go 1.25.5)

## ドキュメント
→ docs/ARCHITECTURE.md  -- アーキテクチャ設計
→ docs/API.md           -- 各パッケージAPI
→ README.md             -- ユーザー向け説明

## 構成
cmd/cc-plans/main.go    -- エントリポイント
internal/cli/            -- Cobraコマンド（プレゼンテーション層）
internal/plan/           -- ドメインモデル+リポジトリ
internal/config/         -- 設定管理
internal/fzf/            -- fzf統合
internal/pager/          -- ページャー統合
internal/renderer/       -- Glamour Markdownレンダリング

## 開発コマンド
make build    -- ビルド
make test     -- 全テスト実行
make lint     -- golangci-lint実行
make fmt      -- goimports実行
make setup    -- 開発環境セットアップ (goimports + lefthook)
make install  -- ~/.local/bin にインストール

## コード規約
- internal/ 配下のみ（公開APIなし）
- テスト: *_test.go をソースと同じディレクトリに配置
- エラー: Sentinel errors (ErrNotFound, ErrAmbiguous 等)
- 新コマンド: internal/cli/ にファイル追加、init() + rootCmd.AddCommand() で登録
- 外部ツール連携: 専用の internal/ パッケージでラップ

## Git ワークフロー
- プリコミットフック: Lefthook (lefthook.yml)
- フック内容: goimports, go vet, golangci-lint, go test
- コミットメッセージ: feat:, fix:, docs: 等の接頭辞

## ツール設定ポインタ
→ .golangci.yml                 -- リンター設定
→ lefthook.yml                  -- Gitフック設定
→ .claude/settings.local.json   -- Claude Code Hook設定
