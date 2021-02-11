package listremote

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/dikaeinstein/godl/internal/pkg/godlutil"
	"github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
	"github.com/tj/go-spin"
)

// New returns the list-remote command
func New() *cobra.Command {
	return &cobra.Command{
		Use:     "list-remote",
		Aliases: []string{"ls-remote"},
		Short:   "List the available remote versions.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listRemoteVersions(&http.Client{})
		},
	}
}

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

func listRemoteVersions(client *http.Client) error {
	url := "https://storage.googleapis.com/golang/?prefix=go1"
	w := url
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	var contents []Content
	go func(ctx context.Context) {
		s := spin.New()
	Loop:
		for {
			select {
			case <-ctx.Done():
				break Loop
			default:
				fmt.Printf("\rfetching remote versions... %s", s.Next())
				// Pause current goroutine to reduce cpu workload
				time.Sleep(50 * time.Millisecond)
			}
		}
	}(ctx)

	for {
		listBucketResult, err := getBinaryReleases(url, client)
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

	versions := mapToVersion(contents)
	// sort in-place comparing version numbers
	sort.Sort(version.Collection(versions))

	cancelFunc()
	fmt.Println()

	for _, v := range versions {
		fmt.Println(v.Original())
	}

	return nil
}

func getBinaryReleases(url string, c *http.Client) (*ListBucketResult, error) {
	res, err := c.Get(url)
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
		return nil, fmt.Errorf("Error decoding xml: %v", err)
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

func mapToVersion(contents []Content) []*version.Version {
	versions := make([]*version.Version, len(contents))
	for i, c := range contents {
		versions[i] = godlutil.GetVersion(c.Key)
	}
	return versions
}
