package cli

import (
	"fmt"

	"github.com/masato-uno/cc-plans/internal/plan"
	"github.com/spf13/cobra"
)

var searchNameOnly bool

func init() {
	searchCmd := &cobra.Command{
		Use:   "search <query>",
		Short: "プランを検索",
		Long:  "プランのファイル名または内容を検索します。",
		Args:  cobra.ExactArgs(1),
		RunE:  runSearch,
	}

	searchCmd.Flags().BoolVarP(&searchNameOnly, "name", "n", false, "ファイル名のみを検索")

	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := args[0]
	repo := plan.NewRepository()

	results, err := repo.Search(query, searchNameOnly)
	if err != nil {
		return fmt.Errorf("検索に失敗しました: %w", err)
	}

	if len(results) == 0 {
		fmt.Printf("'%s' に一致するプランが見つかりません\n", query)
		return nil
	}

	if searchNameOnly {
		for _, r := range results {
			fmt.Println(r.Plan.Name)
		}
		return nil
	}

	// Group by plan
	currentPlan := ""
	for _, r := range results {
		if r.Plan.Name != currentPlan {
			if currentPlan != "" {
				fmt.Println()
			}
			fmt.Printf("=== %s ===\n", r.Plan.Name)
			currentPlan = r.Plan.Name
		}
		fmt.Printf("%4d: %s\n", r.LineNumber, truncate(r.MatchLine, 80))
	}

	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
