package app

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/dikaeinstein/godl/internal/pkg/gv"
	"github.com/dikaeinstein/godl/pkg/text"
	"github.com/hashicorp/go-version"
	"github.com/tj/go-spin"
)

// ListBucketResult represents the list of objects result
type ListBucketResult struct {
	XMLNAME     xml.Name  `xml:"ListBucketResult"`
	Contents    []Content `xml:"Contents"`
	NextMarker  string
	IsTruncated bool
}

// Content represents a ListBucketResult object
type Content struct {
	XMLNAME      xml.Name `xml:"Contents"`
	Key          string
	LastModified time.Time
}

// ListRemote lists remote versions available for install
type ListRemote struct {
	Client  *http.Client
	Timeout time.Duration
}

func (lsRemote *ListRemote) Run(ctx context.Context, sortDirection gv.SortDirection) error {
	url := "https://storage.googleapis.com/golang/?prefix=go1"
	w := url
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var contents []Content
	go func(ctx context.Context) {
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
	}(ctx)

	for {
		listBucketResult, err := lsRemote.getBinaryReleases(url)
		if err != nil {
			return fmt.Errorf("\nerror fetching list: %v", err)
		}

		fl := selectDarwin(listBucketResult)
		contents = append(contents, fl.Contents...)

		// if there's nothing left to fetch
		if !listBucketResult.IsTruncated {
			break
		}

		// update url with marker to fetch next list
		url = w + "&marker=" + listBucketResult.NextMarker
	}

	versions := mapXMLContentToVersion(contents)
	// sort in-place comparing version numbers
	sort.Slice(versions, func(i, j int) bool {
		return gv.CompareVersions(versions[i], versions[j], sortDirection)
	})

	cancel()
	fmt.Println()

	for _, v := range versions {
		fmt.Println(v.Original())
	}

	return nil
}

func (lsRemote *ListRemote) getBinaryReleases(url string) (*ListBucketResult, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), lsRemote.Timeout)
	defer cancelFunc()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	res, err := lsRemote.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %v", url, res.Status)
	}

	var l ListBucketResult
	err = xml.NewDecoder(res.Body).Decode(&l)
	if err != nil {
		return nil, fmt.Errorf("error decoding xml: %v", err)
	}

	return &l, nil
}

func selectDarwin(l *ListBucketResult) ListBucketResult {
	var archiveList ListBucketResult

	for _, r := range l.Contents {
		if strings.Contains(r.Key, "darwin-amd64") && strings.HasSuffix(r.Key, "tar.gz") {
			archiveList.Contents = append(archiveList.Contents, r)
		}
	}

	return archiveList
}

func mapXMLContentToVersion(contents []Content) []*version.Version {
	versions := make([]*version.Version, len(contents))
	for i, c := range contents {
		versions[i] = gv.GetVersion(c.Key)
	}
	return versions
}
