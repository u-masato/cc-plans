package cli

import (
	"errors"
	"fmt"

	"github.com/masato-uno/cc-plans/internal/config"
	"github.com/masato-uno/cc-plans/internal/pager"
	"github.com/masato-uno/cc-plans/internal/plan"
	"github.com/masato-uno/cc-plans/internal/renderer"
	"github.com/spf13/cobra"
)

var (
	showNoPager bool
	showRaw     bool
)

func init() {
	showCmd := &cobra.Command{
		Use:   "show <name>",
		Short: "プラン内容を表示",
		Long:  "指定したプランの内容を表示します。部分一致でも検索できます。",
		Args:  cobra.ExactArgs(1),
		RunE:  runShow,
	}

	showCmd.Flags().BoolVar(&showNoPager, "no-pager", false, "ページャーを使用しない")
	showCmd.Flags().BoolVar(&showRaw, "raw", false, "Markdownレンダリングを無効化")

	rootCmd.AddCommand(showCmd)
}

func runShow(cmd *cobra.Command, args []string) error {
	name := args[0]
	repo := plan.NewRepository()

	content, err := repo.GetContent(name)
	if err != nil {
		if errors.Is(err, plan.ErrNotFound) {
			return fmt.Errorf("プラン '%s' が見つかりません", name)
		}
		if errors.Is(err, plan.ErrAmbiguous) {
			return fmt.Errorf("'%s' に一致するプランが複数あります。より具体的な名前を指定してください", name)
		}
		return fmt.Errorf("プランの読み取りに失敗しました: %w", err)
	}

	if !showRaw && !config.RawMarkdown() && !pager.IsPiped() {
		content = renderer.Render(content)
	}

	return pager.Show(content, !showNoPager)
}
