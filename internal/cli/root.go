package cli

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/masato-uno/cc-plans/internal/config"
	"github.com/masato-uno/cc-plans/internal/fzf"
	"github.com/masato-uno/cc-plans/internal/pager"
	"github.com/masato-uno/cc-plans/internal/plan"
	"github.com/masato-uno/cc-plans/internal/renderer"
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

	result, err := fzf.Select(plans)
	if err != nil {
		return fmt.Errorf("fzf エラー: %w", err)
	}

	if result.Name == "" {
		return nil // User cancelled
	}

	selectedPlan, err := repo.Get(result.Name)
	if err != nil {
		return fmt.Errorf("プランの取得に失敗しました: %w", err)
	}

	switch result.Action {
	case fzf.ActionEdit:
		return openInEditor(selectedPlan.Path)
	case fzf.ActionDelete:
		return deletePlan(repo, result.Name)
	default:
		content, err := repo.GetContent(result.Name)
		if err != nil {
			return fmt.Errorf("プランの読み取りに失敗しました: %w", err)
		}
		if !config.RawMarkdown() && !pager.IsPiped() {
			content = renderer.Render(content)
		}
		return pager.Show(content, true)
	}
}

func openInEditor(path string) error {
	editor := config.Editor()
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func deletePlan(repo *plan.Repository, name string) error {
	fmt.Printf("プラン '%s' を削除しますか? [y/N]: ", name)

	reader := bufio.NewReader(os.Stdin)
	answer, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("入力の読み取りに失敗しました: %w", err)
	}

	answer = strings.TrimSpace(strings.ToLower(answer))
	if answer != "y" && answer != "yes" {
		fmt.Println("削除をキャンセルしました")
		return nil
	}

	if err := repo.Delete(name); err != nil {
		return fmt.Errorf("プランの削除に失敗しました: %w", err)
	}

	fmt.Printf("プラン '%s' を削除しました\n", name)
	return nil
}
