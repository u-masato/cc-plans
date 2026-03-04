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
	p := Pager()
	if p != DefaultPager {
		t.Errorf("Pager() = %q, want %q", p, DefaultPager)
	}
}

func TestPager_EnvOverride(t *testing.T) {
	t.Setenv("PAGER", "more")
	p := Pager()
	if p != "more" {
		t.Errorf("Pager() = %q, want 'more'", p)
	}
}

func TestEditor_Default(t *testing.T) {
	t.Setenv("EDITOR", "")
	os.Unsetenv("EDITOR")
	e := Editor()
	if e != "vim" {
		t.Errorf("Editor() = %q, want 'vim'", e)
	}
}

func TestEditor_EnvOverride(t *testing.T) {
	t.Setenv("EDITOR", "nano")
	e := Editor()
	if e != "nano" {
		t.Errorf("Editor() = %q, want 'nano'", e)
	}
}

func TestRawMarkdown_Unset(t *testing.T) {
	os.Unsetenv("CC_PLANS_RAW")
	if RawMarkdown() {
		t.Error("RawMarkdown() = true, want false when CC_PLANS_RAW is unset")
	}
}

func TestRawMarkdown_Enabled(t *testing.T) {
	t.Setenv("CC_PLANS_RAW", "1")
	if !RawMarkdown() {
		t.Error("RawMarkdown() = false, want true when CC_PLANS_RAW=1")
	}
}

func TestRawMarkdown_OtherValues(t *testing.T) {
	tests := []struct {
		value string
		want  bool
	}{
		{"0", false},
		{"true", false},
		{"yes", false},
		{"", false},
		{"1", true},
	}
	for _, tt := range tests {
		t.Run("CC_PLANS_RAW="+tt.value, func(t *testing.T) {
			if tt.value == "" {
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
