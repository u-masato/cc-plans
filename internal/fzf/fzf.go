package fzf

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/masato-uno/cc-plans/internal/plan"
)

func selfPath() string {
	exe, err := os.Executable()
	if err != nil {
		return "cc-plans"
	}
	return exe
}

// Action represents the action to take on the selected plan.
type Action int

const (
	ActionNone Action = iota
	ActionShow
	ActionEdit
	ActionDelete
)

// SelectResult represents the result of fzf selection.
type SelectResult struct {
	Name   string
	Action Action
}

// IsAvailable returns true if fzf is installed and available.
func IsAvailable() bool {
	_, err := exec.LookPath("fzf")
	return err == nil
}

// Select launches fzf with the given plans and returns the selected plan name and action.
// Returns empty result if user cancels or no selection is made.
func Select(plans []plan.Plan) (SelectResult, error) {
	if len(plans) == 0 {
		return SelectResult{}, nil
	}

	// Build input for fzf: "name\ttitle"
	var input strings.Builder
	for _, p := range plans {
		title := p.Title
		if len(title) > 50 {
			title = title[:47] + "..."
		}
		fmt.Fprintf(&input, "%s\t%s\n", p.Name, title)
	}

	// Build fzf command with preview and key bindings
	cmd := exec.Command("fzf",
		"--height=40%",
		"--reverse",
		"--delimiter=\t",
		"--with-nth=1,2",
		"--preview", fmt.Sprintf("'%s' show --no-pager {1}", selfPath()),
		"--preview-window=right:50%:wrap",
		"--expect=ctrl-e,ctrl-d",
		"--header=Enter: 表示 | ctrl-e: 編集 | ctrl-d: 削除",
	)

	cmd.Stdin = strings.NewReader(input.String())
	cmd.Stderr = os.Stderr

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		// Check if it's just a user cancel (exit code 130 or 1)
		if exitErr, ok := err.(*exec.ExitError); ok {
			code := exitErr.ExitCode()
			if code == 130 || code == 1 {
				return SelectResult{}, nil // User cancelled
			}
		}
		return SelectResult{}, err
	}

	output := stdout.String()
	lines := strings.Split(output, "\n")

	if len(lines) < 2 {
		return SelectResult{}, nil
	}

	// First line is the key pressed (from --expect)
	key := strings.TrimSpace(lines[0])
	// Second line is the selected item
	selection := strings.TrimSpace(lines[1])

	if selection == "" {
		return SelectResult{}, nil
	}

	// Extract just the name (first field before tab)
	name := strings.Split(selection, "\t")[0]

	return SelectResult{Name: name, Action: actionFromKey(key)}, nil
}

// actionFromKey maps an fzf --expect key to an Action.
func actionFromKey(key string) Action {
	switch key {
	case "ctrl-e":
		return ActionEdit
	case "ctrl-d":
		return ActionDelete
	default:
		return ActionShow
	}
}
