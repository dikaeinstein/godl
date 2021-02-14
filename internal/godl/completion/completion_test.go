package completion

import (
	"errors"
	"testing"

	"github.com/dikaeinstein/godl/internal/godl"
	"github.com/dikaeinstein/godl/pkg/fs"
	"github.com/dikaeinstein/godl/test"
	"github.com/spf13/cobra"
)

type fakeSymLinkerFS struct{}

func (fakeSymLinkerFS) Open(name string) (fs.File, error) {
	return nil, nil
}

func (fakeSymLinkerFS) Symlink(oldName, newName string) error {
	return nil
}

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

	tmpHome := t.TempDir()
	tmpSymDir := t.TempDir()

	completion := &completionCmd{
		bashSymlinkDir: tmpSymDir,
		fsys:           fakeSymLinkerFS{},
		homeDir:        tmpHome,
		rootCmd:        godl.New(),
		zshSymlinkDir:  tmpSymDir,
	}
	for name, tC := range testCases {
		t.Run(name, func(t *testing.T) {
			err := completion.Run(tC.shell)
			if err != nil {
				if err.Error() != tC.err.Error() {
					t.Errorf("expected completion(%#v, %#v) => %#v, got %v",
						tC.shell, tmpHome, tC.err, err)
				}
			}
		})
	}
}

func TestCompletionCmdCalledWithNoArgs(t *testing.T) {
	godlCmd := godl.New()
	completion := New(godlCmd)
	godlCmd.RegisterSubCommands([]*cobra.Command{completion})

	_, errOutput := test.ExecuteCommand(t, true, godlCmd, "completion")
	expected := "Error: provide shell to configure e.g bash or zsh\n"
	if errOutput != expected {
		t.Errorf("godl completion failed: expected: %s; got: %s", expected, errOutput)
	}
}
