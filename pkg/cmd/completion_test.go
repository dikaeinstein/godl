package cmd

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"
)

func TestCompletion(t *testing.T) {
	testCases := map[string]struct {
		shell string
		err   error
	}{
		"creates godl file when bash is passed": {"bash", nil},
		"creates _godl file when zsh is passed": {"zsh", nil},
		"returns an error when unknown shell name is passed": {
			"unknown", errors.New("unknown shell passed")},
	}

	tmpDir, err := ioutil.TempDir(".", "_tmpHome")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	for name, tC := range testCases {
		t.Run(name, func(t *testing.T) {
			err := completion(tC.shell, tmpDir)
			if err != nil {
				if err.Error() != tC.err.Error() {
					t.Errorf("expected completion(%#v, %#v) => %#v, got %v",
						tC.shell, tmpDir, tC.err, err)
				}
			}
		})
	}
}

func TestCompletionCmdCalledWithNoArgs(t *testing.T) {
	_, err := executeCommand(rootCmd, "completion")
	expected := "provide shell to configure e.g bash or zsh"
	got := err.Error()
	if got != expected {
		t.Errorf("godl completion: %v", err)
	}
}
