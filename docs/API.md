# cc-plans API リファレンス

## パッケージ一覧

| パッケージ | パス | 説明 |
|-----------|------|------|
| main | `cmd/cc-plans` | エントリーポイント |
| cli | `internal/cli` | CLIコマンド定義 |
| plan | `internal/plan` | ドメインモデル・リポジトリ |
| config | `internal/config` | 設定管理 |
| fzf | `internal/fzf` | fzf連携 |
| pager | `internal/pager` | ページャー連携 |

---

## internal/plan

### 型

#### Plan

```go
type Plan struct {
    Name    string    // ファイル名（拡張子なし）
    Path    string    // フルパス
    ModTime time.Time // 更新日時
    Size    int64     // ファイルサイズ（バイト）
    Title   string    // ファイル先頭の # タイトル
}
```

#### SearchResult

```go
type SearchResult struct {
    Plan       Plan   // マッチしたプラン
    MatchLine  string // マッチした行（内容検索時）
    LineNumber int    // 行番号（1-indexed）
}
```

#### Repository

```go
type Repository struct {
    plansDir string
}
```

### 関数

#### NewRepository

```go
func NewRepository() *Repository
```

新しいRepositoryインスタンスを作成します。
プランディレクトリは `~/.claude/plans/` がデフォルトです。

#### (*Repository) List

```go
func (r *Repository) List() ([]Plan, error)
```

全プランを取得します。

**戻り値:**
- `[]Plan`: プラン一覧
- `error`: `ErrPlansNotExist`（ディレクトリが存在しない場合）

#### (*Repository) Get

```go
func (r *Repository) Get(name string) (*Plan, error)
```

名前でプランを取得します。部分一致に対応しています。

**引数:**
- `name`: プラン名（部分一致可）

**戻り値:**
- `*Plan`: マッチしたプラン
- `error`: `ErrNotFound`（見つからない）、`ErrAmbiguous`（複数マッチ）

#### (*Repository) GetContent

```go
func (r *Repository) GetContent(name string) (string, error)
```

プランの内容を取得します。

**引数:**
- `name`: プラン名（部分一致可）

**戻り値:**
- `string`: ファイル内容
- `error`: `Get()` と同様のエラー

#### (*Repository) Search

```go
func (r *Repository) Search(query string, nameOnly bool) ([]SearchResult, error)
```

プランを検索します。

**引数:**
- `query`: 検索クエリ（大文字小文字無視）
- `nameOnly`: `true` の場合ファイル名のみ検索

**戻り値:**
- `[]SearchResult`: 検索結果

### ユーティリティ関数

#### SortByModTime

```go
func SortByModTime(plans []Plan)
```

プランを更新日時の降順（新しい順）でソートします。

#### SortByName

```go
func SortByName(plans []Plan)
```

プランを名前の昇順でソートします。

### エラー

```go
var (
    ErrNotFound      = errors.New("plan not found")
    ErrAmbiguous     = errors.New("ambiguous plan name: multiple matches found")
    ErrPlansNotExist = errors.New("plans directory does not exist")
)
```

---

## internal/config

### 定数

```go
const DefaultPager = "less"
```

### 関数

#### PlansDir

```go
func PlansDir() string
```

プランディレクトリのパスを返します。

**戻り値:** `~/.claude/plans/`

#### Pager

```go
func Pager() string
```

使用するページャーコマンドを返します。

**戻り値:** 環境変数 `$PAGER` または `"less"`

---

## internal/fzf

### 関数

#### IsAvailable

```go
func IsAvailable() bool
```

fzfがPATH上にインストールされているか判定します。

#### Select

```go
func Select(plans []plan.Plan) (string, error)
```

fzfを起動してプランを選択させます。

**引数:**
- `plans`: 選択肢となるプラン一覧

**戻り値:**
- `string`: 選択されたプラン名（キャンセル時は空文字）
- `error`: fzf実行エラー

**fzfオプション:**
- `--height=40%`
- `--reverse`
- `--preview`: ファイル先頭50行をプレビュー

---

## internal/pager

### 関数

#### IsPiped

```go
func IsPiped() bool
```

標準出力がパイプされているか判定します。

#### Show

```go
func Show(content string, usePager bool) error
```

内容を表示します。

**引数:**
- `content`: 表示する内容
- `usePager`: ページャーを使用するか

**動作:**
- `usePager=false` または `IsPiped()=true` の場合: 標準出力に直接出力
- それ以外: `$PAGER` を起動して表示

---

## internal/cli

### コマンド

#### root

```
cc-plans [flags]
```

引数なしで実行するとインタラクティブモード（fzf）を起動します。
fzf未インストール時は `list` にフォールバックします。

#### list

```
cc-plans list [flags]
cc-plans ls [flags]
```

| フラグ | 短縮 | 説明 |
|--------|------|------|
| `--long` | `-l` | 詳細表示 |
| `--time` | `-t` | 更新順ソート |

#### show

```
cc-plans show <name> [flags]
```

| フラグ | 説明 |
|--------|------|
| `--no-pager` | ページャーを使用しない |

#### search

```
cc-plans search <query> [flags]
```

| フラグ | 短縮 | 説明 |
|--------|------|------|
| `--name` | `-n` | ファイル名のみ検索 |
