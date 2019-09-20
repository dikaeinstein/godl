package cmd

import (
	"os"
	"testing"
)

type testGzUnArchiver struct{}

func (testGzUnArchiver) Unarchive(source, target string) error { return nil }

type fakeRemover struct{}

func (fr fakeRemover) RemoveAll(path string) error { return nil }

func TestInstallRelease(t *testing.T) {
	tests := map[string]struct {
		downloadedVersion string
		installVersion    string
		success           bool
	}{
		"installRelease handles error due to missing downloaded version": {
			"100.1", "1.17", false,
		},
		"installRelease succeeds": {"1.16", "1.16", true},
	}

	tmpDir, err := createTempGodlDownloadDir()
	if err != nil {
		t.Errorf("TestInstallRelease failed: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tmpFile, err := createTempGoBinaryArchive(tmpDir, tc.downloadedVersion)
			defer tmpFile.Close()

			ua := testGzUnArchiver{}
			fr := fakeRemover{}
			err = installRelease(tc.installVersion, tmpDir, ua, fr)
			var got bool
			if err != nil {
				got = false
			} else {
				got = true
			}

			if got != tc.success {
				t.Errorf("Error installing go binary: %v", err)
			}
		})
	}
}

func TestInstallCmdCalledWithNoArgs(t *testing.T) {
	_, err := executeCommand(rootCmd, "install")
	expected := "provide binary archive version to install"
	got := err.Error()
	if got != expected {
		t.Errorf("godl install Unknown error: %v", err)
	}
}

func TestInstallCommandHelp(t *testing.T) {
	_, err := executeCommand(rootCmd, "install", "-h")
	if err != nil {
		t.Errorf("godl install failed: %v", err)
	}
}
