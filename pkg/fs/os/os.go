package os

import (
	"os"

	"github.com/dikaeinstein/godl/pkg/fs"
)

// FS is an os based filesystem.
type FS struct{}

func (FS) Create(name string) (fs.File, error) {
	return os.Create(name)
}

func (FS) Open(name string) (fs.File, error) {
	return os.Open(name)
}

// RemoveAll removes path and it's children
func (FS) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

// Rename renames file from
func (FS) Rename(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}