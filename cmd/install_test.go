package cmd

import "testing"

type testGzUnArchiver struct{}

func (testGz testGzUnArchiver) Unarchive(source, target string) error {
	return nil
}

func TestInstallGoBinary(t *testing.T) {
	tgz := testGzUnArchiver{}
	err := installGoBinary("1.12", tgz)
	if err != nil {
		t.Errorf("Error installing go binary: %v", err)
	}
}
