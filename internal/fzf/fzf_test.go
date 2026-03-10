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
	available := IsAvailable()

	_, err := exec.LookPath("fzf")
	expected := err == nil

	if available != expected {
		t.Errorf("IsAvailable() = %v, want %v", available, expected)
	}
}

func TestActionFromKey(t *testing.T) {
	tests := []struct {
		key  string
		want Action
	}{
		{"ctrl-e", ActionEdit},
		{"ctrl-d", ActionDelete},
		{"", ActionShow},
		{"enter", ActionShow},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			if got := actionFromKey(tt.key); got != tt.want {
				t.Errorf("actionFromKey(%q) = %d, want %d", tt.key, got, tt.want)
			}
		})
	}
}
