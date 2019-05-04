package cmd

import (
	"testing"
)

func TestVersion(t *testing.T) {
	_, err := executeCommand(rootCmd, "version")
	if err != nil {
		t.Errorf("godl version failed: %v", err)
	}
}
