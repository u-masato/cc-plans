package plan

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractMeta_Title(t *testing.T) {
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
			gotTitle, _ := extractMeta(path)
			if gotTitle != tt.want {
				t.Errorf("extractMeta() title = %q, want %q", gotTitle, tt.want)
			}
		})
	}
}

func TestExtractMeta_NonexistentFile(t *testing.T) {
	title, preview := extractMeta("/nonexistent/path/file.md")
	if title != "" {
		t.Errorf("extractMeta(nonexistent) title = %q, want empty", title)
	}
	if preview != "" {
		t.Errorf("extractMeta(nonexistent) preview = %q, want empty", preview)
	}
}

func TestExtractMeta_Preview(t *testing.T) {
	dir := t.TempDir()

	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "basic preview",
			content: "# Title\n\nFirst line\nSecond line\n",
			want:    "First line / Second line",
		},
		{
			name:    "skips headings and blanks",
			content: "# Title\n\n## Section\n\nContent line\n",
			want:    "Content line",
		},
		{
			name:    "truncates long lines",
			content: "# T\n\nあいうえおかきくけこさしすせそたちつてとなにぬねのはひふへほまみむめもやゆよらりるれろわをんアイウエオカキクケコサシスセソ\n",
			want:    "あいうえおかきくけこさしすせそたちつてとなにぬねのはひふへほまみむめもやゆよらりるれろわをんアイウエオカキクケコサシスセ...",
		},
		{
			name:    "empty file",
			content: "",
			want:    "",
		},
		{
			name:    "skips table rows",
			content: "# Title\n\n| Name | Value |\n|------|-------|\n| a | b |\nReal content\n",
			want:    "Real content",
		},
		{
			name:    "max two lines",
			content: "Line one\nLine two\nLine three\n",
			want:    "Line one / Line two",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(dir, tt.name+".md")
			if err := os.WriteFile(path, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}
			_, gotPreview := extractMeta(path)
			if gotPreview != tt.want {
				t.Errorf("extractMeta() preview = %q, want %q", gotPreview, tt.want)
			}
		})
	}
}
