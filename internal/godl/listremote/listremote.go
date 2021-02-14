package listremote

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/dikaeinstein/godl/internal/pkg/gv"
	"github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
	"github.com/tj/go-spin"
)

var sortDirection string

// New returns the list-remote command
func New() *cobra.Command {
	lsRemoteExAsc := "ls-remote -s asc or ls-remote -s=asc"
	lsRemoteExDesc := "ls-remote -s desc or ls-remote -s=desc"

	listRemote := &cobra.Command{
		Use:     "list-remote",
		Aliases: []string{"ls-remote"},
		Example: fmt.Sprintf("%11s\n%38s\n%40s", "ls-remote", lsRemoteExAsc, lsRemoteExDesc),
		Short:   "List the available remote versions.",
		RunE: func(cmd *cobra.Command, args []string) error {
			lsRemote := listRemoteCmd{http.DefaultClient}
			return lsRemote.Run(cmd.Context())
		},
	}

	listRemote.Flags().StringVarP(&sortDirection, "sortDirection", "s", string(gv.Asc),
		"Specify the sort direction of the output of `list-remote`. It sorts in ascending order by default.")

	return listRemote
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

type listRemoteCmd struct {
	c *http.Client
}

func (lsRemote *listRemoteCmd) Run(ctx context.Context) error {
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
				fmt.Printf("\rfetching remote versions... %s", s.Next())
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

	versions := mapToVersion(contents)
	// sort in-place comparing version numbers
	sort.Slice(versions, func(i, j int) bool {
		return gv.CompareVersions(versions[i], versions[j], gv.SortDirection(sortDirection))
	})

	cancel()
	fmt.Println()

	for _, v := range versions {
		fmt.Println(v.Original())
	}

	return nil
}

func (lsRemote *listRemoteCmd) getBinaryReleases(url string) (*ListBucketResult, error) {
	var timeout time.Duration
	if lsRemote.c.Timeout != 0 {
		timeout = lsRemote.c.Timeout
	} else {
		timeout = 5000 * time.Millisecond
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	res, err := lsRemote.c.Do(req)
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

func mapToVersion(contents []Content) []*version.Version {
	versions := make([]*version.Version, len(contents))
	for i, c := range contents {
		versions[i] = gv.GetVersion(c.Key)
	}
	return versions
}
