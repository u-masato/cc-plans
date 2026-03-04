package plan

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractTitle(t *testing.T) {
	dir := t.TempDir()

	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "standard heading",
			content: "# My Plan\n\nSome content\n",
			want:    "My Plan",
		},
		{
			name:    "heading after blank lines",
			content: "\n\n# Title After Blanks\n",
			want:    "Title After Blanks",
		},
		{
			name:    "no heading",
			content: "Just some text\nNo heading here\n",
			want:    "",
		},
		{
			name:    "h2 not h1",
			content: "## Not H1\n",
			want:    "",
		},
		{
			name:    "multiple headings returns first",
			content: "# First\n# Second\n",
			want:    "First",
		},
		{
			name:    "empty file",
			content: "",
			want:    "",
		},
		{
			name:    "heading with special chars",
			content: "# Plan: 日本語テスト (v2)\n",
			want:    "Plan: 日本語テスト (v2)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(dir, tt.name+".md")
			if err := os.WriteFile(path, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}
			got := extractTitle(path)
			if got != tt.want {
				t.Errorf("extractTitle() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestExtractTitle_NonexistentFile(t *testing.T) {
	got := extractTitle("/nonexistent/path/file.md")
	if got != "" {
		t.Errorf("extractTitle(nonexistent) = %q, want empty", got)
	}
}
