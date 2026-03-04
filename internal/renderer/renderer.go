package renderer

import (
	"os"
	"strconv"

	"github.com/charmbracelet/glamour"
	"golang.org/x/term"
)

// Render converts Markdown content to ANSI-styled terminal output.
// On failure, returns the original content unchanged (graceful degradation).
func Render(content string) (string, error) {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(terminalWidth()),
	)
	if err != nil {
		return content, nil
	}

	rendered, err := r.Render(content)
	if err != nil {
		return content, nil
	}

	return rendered, nil
}

func terminalWidth() int {
	// fzf preview pane sets FZF_PREVIEW_COLUMNS
	if cols := os.Getenv("FZF_PREVIEW_COLUMNS"); cols != "" {
		if w, err := strconv.Atoi(cols); err == nil && w > 0 {
			return w
		}
	}

	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		return w
	}

	return 80
}
