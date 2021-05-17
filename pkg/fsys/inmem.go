package fsys

import (
	"bytes"
	"io/fs"
	"os"
	"testing/fstest"
)

// InMemFS is an in-memory based filesystem t.
type InMemFS struct {
	fstest.MapFS
}

type InMemFile struct {
	file *fstest.MapFile
	buf  *bytes.Buffer
}

func (*InMemFile) Close() error { return nil }
func (*InMemFile) Name() string { return "" }
func (inMemFile *InMemFile) Read(p []byte) (int, error) {
	return inMemFile.buf.Read(p)
}
func (*InMemFile) Stat() (fs.FileInfo, error) { return nil, nil }
func (inMemFile *InMemFile) Write(p []byte) (int, error) {
	return inMemFile.buf.Write(p)
}

// NewInMemFS returns a pointer to a new in-memory FS
func NewInMemFS(mapFs fstest.MapFS) *InMemFS {
	return &InMemFS{mapFs}
}

func (inmem *InMemFS) Open(name string) (fs.File, error) {
	return inmem.MapFS.Open(name)
}

func (inmem *InMemFS) Create(name string) (fs.File, error) {
	buf := new(bytes.Buffer)
	f := &InMemFile{&fstest.MapFile{}, buf}
	inmem.MapFS[name] = f.file
	return f, nil
}

func (inmem *InMemFS) Rename(oldPath, newPath string) error {
	f := inmem.MapFS[oldPath]
	inmem.MapFS[newPath] = f
	delete(inmem.MapFS, oldPath)

	return nil
}

// Content returns the stored data of the in-memory InMemFS as a buffer of bytes.
func (inmem *InMemFS) Content(name string) *bytes.Buffer {
	return bytes.NewBuffer(inmem.MapFS[name].Data)
}

func (*InMemFS) RemoveAll(path string) error { return nil }

func (inmem *InMemFS) WriteFile(filename string, data []byte, perm os.FileMode) error {
	f, err := inmem.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	inmem.MapFS[filename].Data = data
	inmem.MapFS[filename].Mode = perm

	return err
}
