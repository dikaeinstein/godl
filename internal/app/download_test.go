package app

import (
	"context"
	"errors"
	"testing"
	"time"
)

var errDownloadFailed = errors.New("download failed")

func TestDownload(t *testing.T) {
	testCases := []struct {
		name        string
		expectedErr error
		returnErr   bool
	}{
		{
			name:        "DownloadSucceeded",
			expectedErr: nil,
			returnErr:   false,
		},
		{
			name:        "DownloadFailed",
			expectedErr: errDownloadFailed,
			returnErr:   true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			dl := fakeDownloader{returnErr: tC.returnErr}

			d := Download{dl, 5 * time.Second}
			version := "1.12"
			err := d.Run(context.Background(), version)
			if !errors.Is(err, tC.expectedErr) {
				t.Fatalf("Download.Run(ctx, %s) = %v; want: %v",
					version, err, tC.expectedErr)
			}
		})
	}
}

type fakeDownloader struct {
	returnErr bool
}

func (d fakeDownloader) Download(_ context.Context, v string) error {
	if d.returnErr {
		return errDownloadFailed
	}

	return nil
}
