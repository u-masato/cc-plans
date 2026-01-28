package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/masato-uno/cc-plans/internal/plan"
	"github.com/spf13/cobra"
)

var (
	listLong    bool
	listSortMod bool
)

func init() {
	listCmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "プラン一覧を表示",
		RunE:    runList,
	}

	listCmd.Flags().BoolVarP(&listLong, "long", "l", false, "詳細表示（更新日時、サイズ、タイトル）")
	listCmd.Flags().BoolVarP(&listSortMod, "time", "t", false, "更新順でソート")

	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	repo := plan.NewRepository()
	plans, err := repo.List()
	if err != nil {
		return fmt.Errorf("プラン一覧の取得に失敗しました: %w", err)
	}

	if len(plans) == 0 {
		fmt.Println("プランがありません")
		return nil
	}

	if listSortMod {
		plan.SortByModTime(plans)
	} else {
		plan.SortByName(plans)
	}

	if listLong {
		return printListLong(plans)
	}

	for _, p := range plans {
		fmt.Println(p.Name)
	}
	return nil
}

func printListLong(plans []plan.Plan) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, p := range plans {
		title := p.Title
		if len(title) > 40 {
			title = title[:37] + "..."
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			p.ModTime.Format("2006-01-02 15:04"),
			formatSize(p.Size),
			p.Name,
			title,
		)
	}
	return w.Flush()
}

func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%dB", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(size)/float64(div), "KMGTPE"[exp])
}
