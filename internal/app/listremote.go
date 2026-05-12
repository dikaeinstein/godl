package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	go_version "github.com/hashicorp/go-version"
	"github.com/tj/go-spin"

	"github.com/dikaeinstein/godl/internal/version"
	"github.com/dikaeinstein/godl/pkg/text"
)

// GoRelease represents a go release as returned by https://go.dev/dl/
// excluding the `files` field.
type GoRelease struct {
	Version string
	Stable  bool
}

// ListRemote lists remote versions available for install
type ListRemote struct {
	Client  *http.Client
	Timeout time.Duration
}

const goDownloadURL = "https://go.dev/dl/?mode=json&include=all"

func (lsRemote *ListRemote) Run(ctx context.Context, sortDirection string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		duration := 50 * time.Millisecond
		s := spin.New()
	Loop:
		for {
			select {
			case <-ctx.Done():
				break Loop
			default:
				fmt.Printf("\r%s... %s", text.Green("fetching remote versions"), s.Next())
				// Pause current goroutine to reduce cpu workload
				time.Sleep(duration)
			}
		}
	}()

	releases, err := lsRemote.fetchGoReleases(ctx, goDownloadURL)
	if err != nil {
		return fmt.Errorf("could not fetch Go releases: %w", err)
	}

	versions := mapGoReleasesToVersions(releases)

	// sort in-place comparing version numbers
	sort.Slice(versions, func(i, j int) bool {
		return version.CompareVersions(versions[i], versions[j], sortDirection)
	})

	cancel()
	fmt.Println()

	for _, v := range versions {
		fmt.Println(v.Original())
	}

	return nil
}

func (lsRemote *ListRemote) fetchGoReleases(ctx context.Context, url string) ([]GoRelease, error) {
	ctx, cancelFunc := context.WithTimeout(ctx, lsRemote.Timeout)
	defer cancelFunc()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	res, err := lsRemote.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil {
			log.Printf("failed to close the resp body: %v\n", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %v", url, res.Status)
	}

	var releases []GoRelease
	err = json.NewDecoder(res.Body).Decode(&releases)
	if err != nil {
		return nil, fmt.Errorf("error decoding json: %v", err)
	}

	return releases, nil
}

func mapGoReleasesToVersions(releases []GoRelease) []*go_version.Version {
	versions := make([]*go_version.Version, len(releases))

	for i, r := range releases {
		versions[i] = version.GetVersion(r.Version)
	}

	return versions
}
