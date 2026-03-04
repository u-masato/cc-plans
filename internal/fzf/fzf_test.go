package fzf

import (
	"os/exec"
	"testing"
)

func TestSelfPath(t *testing.T) {
	path := selfPath()
	if path == "" {
		t.Error("selfPath() returned empty string")
	}
}

func TestIsAvailable(t *testing.T) {
	// Test reflects actual system state - just verify no panic
	available := IsAvailable()

	_, err := exec.LookPath("fzf")
	expected := err == nil

	if available != expected {
		t.Errorf("IsAvailable() = %v, want %v", available, expected)
	}
}

func TestActionConstants(t *testing.T) {
	if ActionNone != 0 {
		t.Errorf("ActionNone = %d, want 0", ActionNone)
	}
	if ActionShow != 1 {
		t.Errorf("ActionShow = %d, want 1", ActionShow)
	}
	if ActionEdit != 2 {
		t.Errorf("ActionEdit = %d, want 2", ActionEdit)
	}
	if ActionDelete != 3 {
		t.Errorf("ActionDelete = %d, want 3", ActionDelete)
	}
}
