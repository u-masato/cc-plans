package fzf

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/masato-uno/cc-plans/internal/plan"
)

// IsAvailable returns true if fzf is installed and available.
func IsAvailable() bool {
	_, err := exec.LookPath("fzf")
	return err == nil
}

// Select launches fzf with the given plans and returns the selected plan name.
// Returns empty string if user cancels or no selection is made.
func Select(plans []plan.Plan) (string, error) {
	if len(plans) == 0 {
		return "", nil
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

	// Build fzf command with preview
	cmd := exec.Command("fzf",
		"--height=40%",
		"--reverse",
		"--delimiter=\t",
		"--with-nth=1,2",
		"--preview", "head -50 ~/.claude/plans/{1}.md",
		"--preview-window=right:50%:wrap",
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
				return "", nil // User cancelled
			}
		}
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", nil
	}

	// Extract just the name (first field before tab)
	parts := strings.Split(result, "\t")
	return parts[0], nil
}
