package inmem

import (
	"bytes"

	"github.com/dikaeinstein/godl/pkg/fs"
)

// File is an in-memory File and doesn't persist it's content to disk.
// It is good for testing purposes.
type File struct {
	*bytes.Buffer
}

func (f *File) Write(p []byte) (int, error) {
	f.Write(p)
	return f.Len(), nil
}

func (*File) Close() error { return nil }

func (*File) Name() string { return "in-mem" }

// FS is an in-memory based filesystem.
type FS struct {
	storage *bytes.Buffer
}

// NewFS returns a pointer to a new in-memory FS
func NewFS(storage *bytes.Buffer) *FS {
	return &FS{storage: storage}
}

func (*FS) Open(name string) (fs.File, error) {
	return &File{new(bytes.Buffer)}, nil
}

func (inmem *FS) Create(name string) (fs.File, error) {
	w := &File{inmem.storage}
	return w, nil
}

func (*FS) Rename(oldPath, newPath string) error {
	return nil
}

// Content returns the stored data of the in-memory FS as a buffer of bytes.
func (inmem *FS) Content() *bytes.Buffer {
	return inmem.storage
}

func (*FS) RemoveAll(path string) error { return nil }
