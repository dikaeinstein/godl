package completion

import (
	"errors"
	"testing"

	"github.com/dikaeinstein/godl/pkg/fs"
)

type fakeSymLinkerFS struct{}

func (fakeSymLinkerFS) Open(name string) (fs.File, error) {
	return nil, nil
}

func (fakeSymLinkerFS) Symlink(oldName, newName string) error {
	return nil
}

type fakeCompletionGenerator struct{}

func (fakeCompletionGenerator) GenerateBashCompletionFile(string) error { return nil }
func (fakeCompletionGenerator) GenerateZshCompletionFile(string) error  { return nil }

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

	completion := &Completion{
		BashSymlinkDir: tmpSymDir,
		FSys:           fakeSymLinkerFS{},
		HomeDir:        tmpHome,
		Generator:      fakeCompletionGenerator{},
		ZshSymlinkDir:  tmpSymDir,
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
