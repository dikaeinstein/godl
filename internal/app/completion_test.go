package app

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"testing/fstest"
)

const completionFileContent string = "this is a test completion"

type fakeSymLinkerFS struct{ fstest.MapFS }

func (fakeSymLinkerFS) Symlink(oldName, newName string) error {
	return nil
}

type fakeCompletionGenerator struct{}

func generateTestCompletion(out io.Writer) error {
	_, err := out.Write([]byte(completionFileContent))
	return err
}

func (fakeCompletionGenerator) GenerateBashCompletion(out io.Writer) error {
	return generateTestCompletion(out)
}

func (fakeCompletionGenerator) GenerateFishCompletion(out io.Writer, includeDesc bool) error {
	return generateTestCompletion(out)
}

func (fakeCompletionGenerator) GenerateZshCompletion(out io.Writer) error {
	return generateTestCompletion(out)
}

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
		AutocompleteDir:     t.TempDir(),
		BashSymlinkDir:      tmpSymDir,
		FishSymlinkDir:      tmpSymDir,
		FS:                  fakeSymLinkerFS{},
		HomeDir:             tmpHome,
		CompletionGenerator: fakeCompletionGenerator{},
		ZshSymlinkDir:       tmpSymDir,
	}
	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			err := completion.Run(tC.shell, out, tC.useDefault)
			if err != nil && err.Error() != tC.err.Error() {
				t.Errorf("expected completion(%#v, %#v) => %#v, got %v",
					tC.shell, tmpHome, tC.err, err)
			}
			content := out.String()
			if err == nil && content != completionFileContent {
				t.Errorf("expected completion file content %s, got %s",
					content, completionFileContent)
			}
		})
	}
}
