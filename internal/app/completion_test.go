package app

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
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

func (fakeCompletionGenerator) GenBashCompletion(out io.Writer) error {
	return generateTestCompletion(out)
}

func (fakeCompletionGenerator) GenFishCompletion(out io.Writer, includeDesc bool) error {
	return generateTestCompletion(out)
}

func (fakeCompletionGenerator) GenZshCompletion(out io.Writer) error {
	return generateTestCompletion(out)
}

func TestCompletion(t *testing.T) {
	testCases := []struct {
		err         error
		description string
		name        string
		shell       string
		useDefault  bool
	}{
		{
			err:         nil,
			description: "creates completion file when bash is passed",
			name:        "ShellBash",
			shell:       "bash",
			useDefault:  false,
		},
		{
			err:         nil,
			description: "creates completion file when zsh is passed",
			name:        "ShellZsh",
			shell:       "zsh",
			useDefault:  true,
		},
		{
			err:         nil,
			description: "creates completion file when fish is passed",
			name:        "ShellFish",
			shell:       "fish",
			useDefault:  true,
		},
		{
			err:         errors.New("unknown shell passed"),
			description: "returns an error when unknown shell name is passed",
			name:        "ShellUnknown",
			shell:       "unknown",
			useDefault:  true,
		},
	}

	tmpSymDir := t.TempDir()

	completion := &Completion{
		AutocompleteDir:     t.TempDir(),
		BashSymlinkDir:      tmpSymDir,
		FishSymlinkDir:      tmpSymDir,
		FS:                  fakeSymLinkerFS{},
		HomeDir:             t.TempDir(),
		CompletionGenerator: fakeCompletionGenerator{},
		ZshSymlinkDir:       tmpSymDir,
	}

	for i := range testCases {
		tC := testCases[i]
		t.Run(tC.name, func(t *testing.T) {
			out := new(bytes.Buffer)

			err := completion.Run(tC.shell, out, tC.useDefault)
			assert.Equal(t, tC.err, err, tC.description)

			if err == nil {
				content := out.String()
				assert.Equal(t, completionFileContent, content, tC.description)
			}
		})
	}
}
