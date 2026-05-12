package fsys

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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

// SymlinkDir creates symlinks for all files in the source directory to the destination directory.
// Excludes directories and symlinks.
func (OsFS) SymlinkDir(oldDir, newDir string) error {
	return filepath.WalkDir(oldDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if d.Type()&os.ModeSymlink != 0 {
			return nil
		}

		oldName := filepath.Join(oldDir, d.Name())
		newName := filepath.Join(newDir, d.Name())

		if err := os.Symlink(oldName, newName); err != nil {
			if errors.Is(err, os.ErrExist) {
				if err := os.Remove(newName); err != nil {
					return fmt.Errorf("failed to unlink existing symlink: %w", err)
				}

				return os.Symlink(oldName, newName)
			}

			return fmt.Errorf("failed to symlink: %w", err)
		}

		return nil
	})
}

// WriteFile writes data to file with filename
func (OsFS) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}
