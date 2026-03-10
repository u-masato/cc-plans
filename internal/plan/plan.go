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
	Preview string   // first meaningful lines as preview
}

// extractTitle reads the first # heading from the file.
func extractTitle(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}
	return ""
}

// extractPreview reads the first meaningful lines from the file, skipping headings and blank lines.
func extractPreview(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	const maxLines = 2
	const maxLen = 60

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() && len(lines) < maxLines {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if len([]rune(line)) > maxLen {
			line = string([]rune(line)[:maxLen]) + "..."
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, " / ")
}

// SearchResult represents a search match.
type SearchResult struct {
	Plan       Plan
	MatchLine  string // the line that matched (for content search)
	LineNumber int    // line number of the match (1-indexed)
}
