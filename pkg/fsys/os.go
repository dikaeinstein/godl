package fsys

import (
	"io/fs"
	"os"
)

// OsFS is an os based filesystem.
type OsFS struct{}

func (OsFS) Create(name string) (fs.File, error) {
	return os.Create(name)
}

func (OsFS) Open(name string) (fs.File, error) {
	return os.Open(name)
}

// RemoveAll removes path and it's children
func (OsFS) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

// Rename renames file from oldPath to newPath
func (OsFS) Rename(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

// Symlink creates newname as a symbolic link to oldname
func (OsFS) Symlink(oldName, newName string) error {
	return os.Symlink(oldName, newName)
}

// WriteFile writes data to file with filename
func (OsFS) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}
