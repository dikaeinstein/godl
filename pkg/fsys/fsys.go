package fsys

import (
	"fmt"
	"io/fs"
	"os"
)

// CreatFS is a filesystem that can create a new file.
type CreatFS interface {
	fs.FS
	Create(name string) (fs.File, error)
}

// RenameFS is a filesystem that can rename a file.
type RenameFS interface {
	fs.FS
	Rename(oldPath, newPath string) error
}

// SymlinkFS is a filesystem that can symlink a file.
type SymlinkFS interface {
	fs.FS
	Symlink(oldName, newName string) error
}

// RemoveAllFS is a filesystem that remove path and it's children that it may contain.
type RemoveAllFS interface {
	fs.FS
	RemoveAll(path string) error
}

// WriteFileFS is a filesystem that can write data to a file.
type WriteFileFS interface {
	fs.FS
	WriteFile(filename string, data []byte, perm os.FileMode) error
}

// Create a new file using the given filesystem.
func Create(fsys fs.FS, name string) (fs.File, error) {
	if fsys, ok := fsys.(CreatFS); ok {
		return fsys.Create(name)
	}

	return nil, fmt.Errorf("create %s: operation not supported", name)
}

// Rename a file using the given filesystem.
func Rename(fsys fs.FS, oldPath, newPath string) error {
	if fsys, ok := fsys.(RenameFS); ok {
		return fsys.Rename(oldPath, newPath)
	}

	return fmt.Errorf("rename %s %s: operation not supported", oldPath, newPath)
}

// Symlink a file using the given filesystem.
func Symlink(fsys fs.FS, oldName, newName string) error {
	if fsys, ok := fsys.(SymlinkFS); ok {
		return fsys.Symlink(oldName, newName)
	}

	return fmt.Errorf("symlink %s %s: operation not supported", oldName, newName)
}

// RemoveAll path and it's children using the given filesystem.
func RemoveAll(fsys fs.FS, path string) error {
	if fsys, ok := fsys.(RemoveAllFS); ok {
		return fsys.RemoveAll(path)
	}

	return fmt.Errorf("removeAll %s: operation not supported", path)
}

func WriteFile(fsys fs.FS, filename string, data []byte, perm os.FileMode) error {
	if fsys, ok := fsys.(WriteFileFS); ok {
		return fsys.WriteFile(filename, data, perm)
	}

	return fmt.Errorf("writeFile %s %s %s: operation not supported", filename, data, perm)
}
