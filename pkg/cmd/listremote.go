package cmd

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/tj/go-spin"
)

func init() {
	rootCmd.AddCommand(listRemote)
}

var listRemote = &cobra.Command{
	Use:     "list-remote",
	Aliases: []string{"ls-remote"},
	Short:   "List the available remote versions.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listRemoteVersions(&http.Client{})
	},
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

	var sortedContents []Content
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

		fl := filterBucketList(listBucketResult)
		sortedContents = append(sortedContents, fl.Contents...)

		// if there's nothing left to fetch
		if !listBucketResult.IsTruncated {
			break
		}

		// update url with marker to fetch next list
		url = w + "&marker=" + listBucketResult.NextMarker
	}

	// sort in-place using the LastModified timestamp
	sort.Slice(sortedContents, func(i, j int) bool {
		return sortedContents[i].LastModified.
			Before(sortedContents[j].LastModified)
	})

	cancelFunc()
	fmt.Println()

	for _, c := range sortedContents {
		v := strings.Split(c.Key, ".darwin-amd64")
		fmt.Println(strings.TrimPrefix(v[0], "go"))
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

func filterBucketList(l *ListBucketResult) ListBucketResult {
	var archiveList ListBucketResult

	for _, r := range l.Contents {
		if strings.Contains(r.Key, "darwin-amd64") && strings.HasSuffix(r.Key, "tar.gz") {
			archiveList.Contents = append(archiveList.Contents, r)
		}
	}

	return archiveList
}
