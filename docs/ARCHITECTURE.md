# cc-plans アーキテクチャ

## 概要

cc-plans は `~/.claude/plans/` 配下のMarkdownファイルを管理するCLIツールです。
クリーンアーキテクチャに基づき、関心事を分離した設計になっています。

## ディレクトリ構成

```
cc-plans/
├── cmd/cc-plans/main.go    # エントリーポイント
├── internal/
│   ├── cli/                # プレゼンテーション層（Cobraコマンド）
│   ├── plan/               # ドメイン層（Plan構造体、Repository）
│   ├── config/             # 設定管理
│   ├── fzf/                # 外部ツール連携（fzf）
│   └── pager/              # 外部ツール連携（ページャー）
├── go.mod
└── Makefile
```

## レイヤー構成

```
┌─────────────────────────────────────────────────────┐
│  CLI Layer (internal/cli)                           │
│  - root.go, list.go, show.go, search.go            │
│  - Cobraコマンド定義、ユーザー入出力               │
└─────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────┐
│  Domain Layer (internal/plan)                       │
│  - plan.go: Plan構造体、SearchResult               │
│  - repository.go: ファイルI/O、検索ロジック        │
└─────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────┐
│  Infrastructure Layer                               │
│  - config/: 設定値取得（環境変数、パス）           │
│  - fzf/: fzf外部プロセス呼び出し                   │
│  - pager/: ページャー外部プロセス呼び出し          │
└─────────────────────────────────────────────────────┘
```

## データフロー

### list コマンド

```
User → list.go → Repository.List() → ファイルシステム
                                    ↓
User ← list.go ← []Plan ←──────────┘
```

### show コマンド

```
User → show.go → Repository.Get() → ファイルシステム
              → Repository.GetContent()
                         ↓
User ← pager.Show() ← content
```

### インタラクティブモード

```
User → root.go → fzf.IsAvailable() → exec.LookPath
              → Repository.List()
              → fzf.Select() → fzfプロセス
                            ↓
User ← pager.Show() ← selected plan content
```

## 主要コンポーネント

### Plan (internal/plan/plan.go)

プランファイルのドメインモデル。

```go
type Plan struct {
    Name    string    // ファイル名（拡張子なし）
    Path    string    // フルパス
    ModTime time.Time // 更新日時
    Size    int64     // ファイルサイズ
    Title   string    // 先頭の # タイトル
}
```

### Repository (internal/plan/repository.go)

ファイルシステムからのプラン取得を担当。

| メソッド | 説明 |
|---------|------|
| `List()` | 全プラン取得 |
| `Get(name)` | 名前でプラン取得（部分一致対応） |
| `GetContent(name)` | プラン内容取得 |
| `Search(query, nameOnly)` | 検索 |

### fzf (internal/fzf/fzf.go)

fzf連携を担当。

| 関数 | 説明 |
|------|------|
| `IsAvailable()` | fzfがインストールされているか判定 |
| `Select(plans)` | fzfを起動し選択結果を返す |

### pager (internal/pager/pager.go)

ページャー連携を担当。

| 関数 | 説明 |
|------|------|
| `IsPiped()` | 標準出力がパイプかどうか判定 |
| `Show(content, usePager)` | ページャーで内容表示 |

## 設計判断

### 部分一致検索

`Repository.Get()` は以下の優先順位で検索：
1. 完全一致
2. 部分一致（大文字小文字無視）
3. 複数マッチ時は `ErrAmbiguous` を返す

### fzfフォールバック

fzf未インストール時は `list` コマンドにフォールバック。
ユーザーの環境に依存せず動作可能。

### ページャー自動無効化

- `--no-pager` フラグ指定時
- 標準出力がパイプの場合（`os.ModeCharDevice` で判定）

## 依存関係

```
cmd/cc-plans/main.go
    └── internal/cli
            ├── internal/plan
            │       └── internal/config
            ├── internal/fzf
            │       └── internal/plan
            ├── internal/pager
            │       └── internal/config
            └── github.com/spf13/cobra
```

## 拡張ポイント

| 機能 | 拡張方法 |
|------|----------|
| 新コマンド追加 | `internal/cli/` に新ファイル作成、`init()` で `rootCmd.AddCommand()` |
| 出力フォーマット | `list.go` に `--format` フラグ追加 |
| 複数ディレクトリ対応 | `config.PlansDir()` を複数パス対応に変更 |
