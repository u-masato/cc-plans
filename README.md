# cc-plans

Claude Code プラン管理CLIツール。`~/.claude/plans/` にあるプランファイルを管理します。

## インストール

```bash
# ビルド
make build

# ~/.local/bin にインストール
make install
```

## 使い方

### インタラクティブモード

```bash
cc-plans
```

引数なしで実行すると、fzfでプランを選択し、内容を表示します。
fzf未インストール時は一覧表示にフォールバックします。

### 一覧表示

```bash
cc-plans list          # プラン一覧
cc-plans ls            # エイリアス
cc-plans ls -l         # 詳細表示（更新日時、サイズ、タイトル）
cc-plans ls -t         # 更新順ソート
cc-plans ls -lt        # 詳細 + 更新順
```

### 内容表示

```bash
cc-plans show <name>              # プラン内容を表示（$PAGERを使用）
cc-plans show peaceful            # 部分一致OK
cc-plans show peaceful --no-pager # ページャーなし
```

### 検索

```bash
cc-plans search <query>     # 内容検索
cc-plans search "HOGE"      # 例: HOGEを含むプランを検索
cc-plans search -n foo      # ファイル名のみ検索
```

## 環境変数

| 変数 | 説明 | デフォルト |
|------|------|-----------|
| `PAGER` | ページャーコマンド | `less` |

## 依存関係

- Go 1.21+
- fzf（オプション、インタラクティブモード用）

## ライセンス

MIT
