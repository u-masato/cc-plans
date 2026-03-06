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
	// t.Setenv captures original value for cleanup, then we unset
	t.Setenv("PAGER", "")
	os.Unsetenv("PAGER")
	t.Cleanup(func() {
		// t.Setenv's cleanup will restore original value
	})
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
	t.Cleanup(func() {})
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
	t.Setenv("CC_PLANS_RAW", "")
	os.Unsetenv("CC_PLANS_RAW")
	t.Cleanup(func() {})
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
			t.Setenv("CC_PLANS_RAW", "")
			if tt.unset {
				os.Unsetenv("CC_PLANS_RAW")
			} else {
				os.Setenv("CC_PLANS_RAW", tt.value)
			}
			t.Cleanup(func() {})
			if got := RawMarkdown(); got != tt.want {
				t.Errorf("RawMarkdown() = %v, want %v", got, tt.want)
			}
		})
	}
}
