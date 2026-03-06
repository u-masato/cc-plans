package pager

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestAppendLessEnv_NotSet(t *testing.T) {
	env := []string{"HOME=/home/user", "PATH=/usr/bin"}
	result := appendLessEnv(env)

	found := false
	for _, e := range result {
		if e == "LESS=FRX" {
			found = true
			break
		}
	}
	if !found {
		t.Error("appendLessEnv did not add LESS=FRX when LESS was not set")
	}
	// Original entries should be preserved
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
	// Should keep original LESS value
	found := false
	for _, e := range result {
		if e == "LESS=-RS" {
			found = true
		}
		if e == "LESS=FRX" {
			t.Error("appendLessEnv overwrote existing LESS value")
		}
	}
	if !found {
		t.Error("original LESS value was lost")
	}
}

func TestAppendLessEnv_EmptyEnv(t *testing.T) {
	result := appendLessEnv(nil)
	if len(result) != 1 || result[0] != "LESS=FRX" {
		t.Errorf("appendLessEnv(nil) = %v, want [LESS=FRX]", result)
	}
}

func TestAppendLessEnv_LessEmptyValue(t *testing.T) {
	// LESS= (set but empty) should still count as "already set"
	env := []string{"LESS="}
	result := appendLessEnv(env)
	if len(result) != 1 {
		t.Error("appendLessEnv should not add LESS=FRX when LESS is already set (even if empty)")
	}
}

func TestShow_NoPager(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	content := "test content\n"
	err := Show(content, false)

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("Show returned error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	if buf.String() != content {
		t.Errorf("Show output = %q, want %q", buf.String(), content)
	}
}

func TestIsPiped(t *testing.T) {
	// In test environment, stdout is typically piped
	// We mainly verify it doesn't panic
	_ = IsPiped()
}

func TestShow_EmptyContent(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := Show("", false)

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("Show returned error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	if buf.String() != "" {
		t.Errorf("Show output = %q, want empty", buf.String())
	}
}

func TestShow_LargeContent(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	content := strings.Repeat("line\n", 1000)
	err := Show(content, false)

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("Show returned error: %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	if buf.String() != content {
		t.Error("Show did not output all content")
	}
}
