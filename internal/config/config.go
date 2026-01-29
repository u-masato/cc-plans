package config

import (
	"os"
	"path/filepath"
)

const (
	DefaultPager = "less"
)

// PlansDir returns the path to the plans directory.
func PlansDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".claude", "plans")
}

// Pager returns the pager command to use.
func Pager() string {
	if p := os.Getenv("PAGER"); p != "" {
		return p
	}
	return DefaultPager
}

// Editor returns the editor command to use.
func Editor() string {
	if e := os.Getenv("EDITOR"); e != "" {
		return e
	}
	return "vim"
}
