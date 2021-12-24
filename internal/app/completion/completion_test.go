package completion

import (
	"errors"
	"io"
	"testing"
	"testing/fstest"
)

type fakeSymLinkerFS struct{ fstest.MapFS }

func (fakeSymLinkerFS) Symlink(oldName, newName string) error {
	return nil
}

type fakeCompletionGenerator struct{}

func (fakeCompletionGenerator) GenerateBashCompletion(io.Writer) error       { return nil }
func (fakeCompletionGenerator) GenerateFishCompletion(io.Writer, bool) error { return nil }
func (fakeCompletionGenerator) GenerateZshCompletion(io.Writer) error        { return nil }

func TestCompletion(t *testing.T) {
	testCases := []struct {
		name       string
		shell      string
		useDefault bool
		err        error
	}{
		{"creates completion file when bash is passed", "bash", false, nil},
		{"creates completion file when zsh is passed", "zsh", true, nil},
		{"creates completion file when fish is passed", "fish", true, nil},
		{
			"returns an error when unknown shell name is passed", "unknown", true, errors.New("unknown shell passed"),
		},
	}

	tmpHome := t.TempDir()
	tmpSymDir := t.TempDir()

	completion := &Completion{
		AutocompleteDir: t.TempDir(),
		BashSymlinkDir:  tmpSymDir,
		FishSymlinkDir:  tmpSymDir,
		FS:              fakeSymLinkerFS{},
		HomeDir:         tmpHome,
		Generator:       fakeCompletionGenerator{},
		ZshSymlinkDir:   tmpSymDir,
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.name, func(t *testing.T) {
			err := completion.Run(tC.shell, io.Discard, tC.useDefault)
			if err != nil && err.Error() != tC.err.Error() {
				t.Errorf("expected completion(%#v, %#v) => %#v, got %v",
					tC.shell, tmpHome, tC.err, err)
			}
		})
	}
}
