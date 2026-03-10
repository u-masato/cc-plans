package pager

import (
	"bytes"
	"os"
	"slices"
	"strings"
	"testing"
)

func TestAppendLessEnv_NotSet(t *testing.T) {
	env := []string{"HOME=/home/user", "PATH=/usr/bin"}
	result := appendLessEnv(env)

	if !slices.Contains(result, "LESS=RX") {
		t.Error("appendLessEnv did not add LESS=RX when LESS was not set")
	}
	if len(result) != len(env)+1 {
		t.Errorf("expected %d entries, got %d", len(env)+1, len(result))
	}
}

func TestAppendLessEnv_AlreadySet(t *testing.T) {
	env := []string{"HOME=/home/user", "LESS=-RS", "PATH=/usr/bin"}
	result := appendLessEnv(env)

	if len(result) != len(env) {
		t.Error("appendLessEnv should not modify env when LESS is already set")
	}
	if !slices.Contains(result, "LESS=-RS") {
		t.Error("original LESS value was lost")
	}
	if slices.Contains(result, "LESS=RX") {
		t.Error("appendLessEnv overwrote existing LESS value")
	}
}

func TestAppendLessEnv_EmptyEnv(t *testing.T) {
	result := appendLessEnv(nil)
	if len(result) != 1 || result[0] != "LESS=RX" {
		t.Errorf("appendLessEnv(nil) = %v, want [LESS=RX]", result)
	}
}

func TestAppendLessEnv_LessEmptyValue(t *testing.T) {
	// LESS= (set but empty) should still count as "already set"
	env := []string{"LESS="}
	result := appendLessEnv(env)
	if len(result) != 1 {
		t.Error("appendLessEnv should not add LESS=RX when LESS is already set (even if empty)")
	}
}

// captureShow calls Show with usePager=false and captures its stdout output.
func captureShow(t *testing.T, content string) string {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := Show(content, false)

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("Show returned error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func TestShow_NoPager(t *testing.T) {
	content := "test content\n"
	got := captureShow(t, content)
	if got != content {
		t.Errorf("Show output = %q, want %q", got, content)
	}
}

func TestIsPiped(t *testing.T) {
	// In test environment, stdout is typically piped.
	// We verify it returns without panic.
	_ = IsPiped()
}

func TestShow_EmptyContent(t *testing.T) {
	got := captureShow(t, "")
	if got != "" {
		t.Errorf("Show output = %q, want empty", got)
	}
}

func TestShow_LargeContent(t *testing.T) {
	content := strings.Repeat("line\n", 1000)
	got := captureShow(t, content)
	if got != content {
		t.Error("Show did not output all content")
	}
}
