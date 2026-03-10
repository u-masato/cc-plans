package config

import (
	"os"
	"strings"
	"testing"
)

func TestPlansDir(t *testing.T) {
	dir := PlansDir()
	if dir == "" {
		t.Fatal("PlansDir() returned empty string")
	}
	if !strings.HasSuffix(dir, ".claude/plans") {
		t.Errorf("PlansDir() = %q, want suffix '.claude/plans'", dir)
	}
}

func TestPager_Default(t *testing.T) {
	t.Setenv("PAGER", "")
	os.Unsetenv("PAGER")

	if got := Pager(); got != DefaultPager {
		t.Errorf("Pager() = %q, want %q", got, DefaultPager)
	}
}

func TestPager_EnvOverride(t *testing.T) {
	t.Setenv("PAGER", "more")

	if got := Pager(); got != "more" {
		t.Errorf("Pager() = %q, want 'more'", got)
	}
}

func TestEditor_Default(t *testing.T) {
	t.Setenv("EDITOR", "")
	os.Unsetenv("EDITOR")

	if got := Editor(); got != "vim" {
		t.Errorf("Editor() = %q, want 'vim'", got)
	}
}

func TestEditor_EnvOverride(t *testing.T) {
	t.Setenv("EDITOR", "nano")

	if got := Editor(); got != "nano" {
		t.Errorf("Editor() = %q, want 'nano'", got)
	}
}

func TestRawMarkdown(t *testing.T) {
	tests := []struct {
		value string
		unset bool
		want  bool
	}{
		{"0", false, false},
		{"true", false, false},
		{"yes", false, false},
		{"", true, false},
		{"1", false, true},
	}
	for _, tt := range tests {
		name := "CC_PLANS_RAW=" + tt.value
		if tt.unset {
			name = "CC_PLANS_RAW=<unset>"
		}
		t.Run(name, func(t *testing.T) {
			if tt.unset {
				t.Setenv("CC_PLANS_RAW", "")
				os.Unsetenv("CC_PLANS_RAW")
			} else {
				t.Setenv("CC_PLANS_RAW", tt.value)
			}
			if got := RawMarkdown(); got != tt.want {
				t.Errorf("RawMarkdown() = %v, want %v", got, tt.want)
			}
		})
	}
}
