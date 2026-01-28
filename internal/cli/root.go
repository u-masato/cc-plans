package cli

import (
	"fmt"
	"os"

	"github.com/masato-uno/cc-plans/internal/fzf"
	"github.com/masato-uno/cc-plans/internal/pager"
	"github.com/masato-uno/cc-plans/internal/plan"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cc-plans",
	Short: "Claude Code プラン管理CLIツール",
	Long:  `cc-plans は ~/.claude/plans/ にあるプランファイルを管理するCLIツールです。`,
	RunE:  runInteractive,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runInteractive(cmd *cobra.Command, args []string) error {
	// fzf 未インストール時は list にフォールバック
	if !fzf.IsAvailable() {
		return runList(cmd, args)
	}

	repo := plan.NewRepository()
	plans, err := repo.List()
	if err != nil {
		return fmt.Errorf("プラン一覧の取得に失敗しました: %w", err)
	}

	if len(plans) == 0 {
		fmt.Println("プランがありません")
		return nil
	}

	plan.SortByModTime(plans)

	selected, err := fzf.Select(plans)
	if err != nil {
		return fmt.Errorf("fzf エラー: %w", err)
	}

	if selected == "" {
		return nil // User cancelled
	}

	// Show the selected plan
	content, err := repo.GetContent(selected)
	if err != nil {
		return fmt.Errorf("プランの読み取りに失敗しました: %w", err)
	}

	return pager.Show(content, true)
}
