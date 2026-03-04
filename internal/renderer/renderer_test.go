package renderer

import (
	"os"
	"strings"
	"testing"
)

func TestRender_BasicMarkdown(t *testing.T) {
	input := "# Hello\n\nThis is a paragraph.\n"
	rendered, err := Render(input)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	// glamour should produce some ANSI output different from raw input
	if rendered == input {
		t.Error("Render returned unchanged input; expected ANSI-styled output")
	}
	// Should still contain the text content
	if !strings.Contains(rendered, "Hello") {
		t.Error("rendered output missing 'Hello'")
	}
	if !strings.Contains(rendered, "paragraph") {
		t.Error("rendered output missing 'paragraph'")
	}
}

func TestRender_CodeBlock(t *testing.T) {
	input := "```go\nfmt.Println(\"hello\")\n```\n"
	rendered, err := Render(input)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	if !strings.Contains(rendered, "Println") {
		t.Error("rendered output missing code content")
	}
}

func TestRender_EmptyInput(t *testing.T) {
	rendered, err := Render("")
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	// Should not panic or error on empty input
	_ = rendered
}

func TestRender_List(t *testing.T) {
	input := "- item1\n- item2\n- item3\n"
	rendered, err := Render(input)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	if !strings.Contains(rendered, "item1") {
		t.Error("rendered output missing list item")
	}
}

func TestTerminalWidth_FZFPreviewColumns(t *testing.T) {
	t.Setenv("FZF_PREVIEW_COLUMNS", "120")
	w := terminalWidth()
	if w != 120 {
		t.Errorf("terminalWidth() = %d, want 120", w)
	}
}

func TestTerminalWidth_InvalidFZFPreviewColumns(t *testing.T) {
	t.Setenv("FZF_PREVIEW_COLUMNS", "notanumber")
	w := terminalWidth()
	// Should fall through to term.GetSize or default
	if w <= 0 {
		t.Errorf("terminalWidth() = %d, want > 0", w)
	}
}

func TestTerminalWidth_ZeroFZFPreviewColumns(t *testing.T) {
	t.Setenv("FZF_PREVIEW_COLUMNS", "0")
	w := terminalWidth()
	// 0 is invalid, should fall through
	if w <= 0 {
		t.Errorf("terminalWidth() = %d, want > 0", w)
	}
}

func TestTerminalWidth_Unset(t *testing.T) {
	os.Unsetenv("FZF_PREVIEW_COLUMNS")
	w := terminalWidth()
	// In test environment (no tty), should return 80 as default
	if w <= 0 {
		t.Errorf("terminalWidth() = %d, want > 0", w)
	}
}
