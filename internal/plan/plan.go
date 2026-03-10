package plan

import (
	"bufio"
	"os"
	"strings"
	"time"
)

// Plan represents a plan file.
type Plan struct {
	Name    string    // filename without extension
	Path    string    // full path
	ModTime time.Time // modification time
	Size    int64     // file size
	Title   string    // first # heading from file
	Preview string    // first meaningful lines as preview
}

// extractMeta reads the title and preview from a single file scan.
func extractMeta(path string) (title, preview string) {
	f, err := os.Open(path)
	if err != nil {
		return "", ""
	}
	defer f.Close()

	const maxPreviewLines = 2
	const maxLen = 60

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if title == "" && strings.HasPrefix(line, "# ") {
			title = strings.TrimPrefix(line, "# ")
			continue
		}

		if len(lines) >= maxPreviewLines {
			break
		}

		if isSkippableLine(line) {
			continue
		}

		if len([]rune(line)) > maxLen {
			line = string([]rune(line)[:maxLen]) + "..."
		}
		lines = append(lines, line)
	}

	preview = strings.Join(lines, " / ")
	return title, preview
}

// isSkippableLine returns true for lines that should not appear in preview.
func isSkippableLine(line string) bool {
	if line == "" {
		return true
	}
	// headings
	if strings.HasPrefix(line, "#") {
		return true
	}
	// table separator (e.g. |---|---|)
	trimmed := strings.ReplaceAll(strings.ReplaceAll(line, "-", ""), "|", "")
	trimmed = strings.ReplaceAll(trimmed, ":", "")
	trimmed = strings.TrimSpace(trimmed)
	if trimmed == "" && strings.Contains(line, "|") {
		return true
	}
	// table header rows (e.g. | Name | Value |)
	if strings.HasPrefix(line, "|") && strings.HasSuffix(line, "|") {
		return true
	}
	return false
}

// SearchResult represents a search match.
type SearchResult struct {
	Plan       Plan
	MatchLine  string // the line that matched (for content search)
	LineNumber int    // line number of the match (1-indexed)
}
