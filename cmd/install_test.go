package cmd

import (
	"os"
	"testing"
)

type testGzUnArchiver struct{}

func (testGz testGzUnArchiver) Unarchive(source, target string) error {
	return nil
}

type installGoBinaryTest struct {
	downloadedVersion string
	installVersion    string
	success           bool
}

var installGoBinaryTestCases = []installGoBinaryTest{
	{"100.1", "1.17", false},
	{"1.16", "1.16", true},
}

func TestInstallGoBinary(t *testing.T) {
	tmpDir, err := createTempGodlDownloadDir()
	if err != nil {
		t.Errorf("InstallGoBinary failed: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	for _, c := range installGoBinaryTestCases {
		tmpFile, err := createTempGoBinaryArchive(tmpDir, c.downloadedVersion)
		defer tmpFile.Close()

		tgz := testGzUnArchiver{}
		err = installGoBinary(c.installVersion, tmpDir, tgz)
		var got bool
		if err != nil {
			got = false
		} else {
			got = true
		}

		if got != c.success {
			t.Errorf("Error installing go binary: %v", err)
		}
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
