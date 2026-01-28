package pager

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/masato-uno/cc-plans/internal/config"
)

// IsPiped returns true if stdout is piped (not a terminal).
func IsPiped() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) == 0
}

// Show displays content using the configured pager.
// If usePager is false or stdout is piped, it writes directly to stdout.
func Show(content string, usePager bool) error {
	if !usePager || IsPiped() {
		_, err := io.WriteString(os.Stdout, content)
		return err
	}

	pagerCmd := config.Pager()
	parts := strings.Fields(pagerCmd)
	if len(parts) == 0 {
		_, err := io.WriteString(os.Stdout, content)
		return err
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdin = strings.NewReader(content)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
