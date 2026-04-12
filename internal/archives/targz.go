package archives

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/mholt/archives"
)

// TarGZ provides a method for unarchiving .tar.gz files.
type TarGZ struct {
	format archives.CompressedArchive
}

// NewTarGZ returns a new TarGZ instance with default settings.
func NewTarGZ() TarGZ {
	return TarGZ{
		format: archives.CompressedArchive{
			Compression: archives.Gz{},
			Archival:    archives.Tar{},
			Extraction:  archives.Tar{},
		},
	}
}

// Unarchive unarchives the given archive file into the destination folder.
func (tgz TarGZ) Unarchive(
	ctx context.Context,
	source, destination string,
) error {
	if !fileExists(destination) {
		if err := os.MkdirAll(destination, 0o755); err != nil {
			return fmt.Errorf("%s: making directory: %w", destination, err)
		}
	}

	sourceArchive, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("%s: opening archive: %w", source, err)
	}
	defer sourceArchive.Close()

	handleFile := func(ctx context.Context, fileInfo archives.FileInfo) error {
		return extractArchive(fileInfo, destination)
	}

	return tgz.format.Extract(ctx, sourceArchive, handleFile)
}

func extractArchive(fileInfo archives.FileInfo, destination string) error {
	to := filepath.Join(destination, path.Clean(fileInfo.Name()))

	if fileInfo.IsDir() {
		return os.MkdirAll(to, fileInfo.Mode())
	}

	if fileInfo.Mode()&os.ModeSymlink != 0 {
		if fileInfo.LinkTarget == "" {
			return fmt.Errorf("symlink target is empty")
		}
		return os.Symlink(fileInfo.LinkTarget, to)
	}

	handle, err := fileInfo.Open()
	if err != nil {
		return err
	}
	defer handle.Close()

	dest, err := os.Create(to)
	if err != nil {
		return err
	}
	defer dest.Close()

	if _, err := io.Copy(dest, handle); err != nil {
		return err
	}

	return os.Chmod(to, fileInfo.Mode())
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}
