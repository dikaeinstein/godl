// Copyright © 2019 Onyedikachi Solomon Okwa <solozyokwa@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package list

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/dikaeinstein/godl/internal/pkg/godlutil"
	"github.com/dikaeinstein/godl/internal/pkg/gv"
	"github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
)

var sortDirection string

// New returns the list command
func New() *cobra.Command {
	lsExAsc := "ls -s asc or ls -s=asc"
	lsExDesc := "ls -s desc or ls -s=desc"

	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List the downloaded versions.",
		Example: fmt.Sprintf("%4s\n%24s\n%26s", "ls", lsExAsc, lsExDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := godlutil.GetDownloadDir()
			if err != nil {
				return err
			}

			ls := listCmd{}
			return ls.Run(d)
		},
	}

	list.Flags().StringVarP(&sortDirection, "sortDirection", "s", string(gv.Asc),
		"Specify the sort direction of the output of `list`. It sorts in ascending order by default.")

	return list
}

type listCmd struct{}

func (listCmd) Run(downloadDir string) error {
	// Create download directory and its parent
	godlutil.Must(os.MkdirAll(downloadDir, os.ModePerm))

	files, err := ioutil.ReadDir(downloadDir)
	if err != nil {
		return err
	}

	versions := mapToVersion(files)
	// sort in-place comparing version numbers
	sort.Slice(versions, func(i, j int) bool {
		return gv.CompareVersions(versions[i], versions[j], gv.SortDirection(sortDirection))
	})

	for _, v := range versions {
		fmt.Println(v.Original())
	}

	return nil
}

func mapToVersion(files []os.FileInfo) []*version.Version {
	versions := []*version.Version{}
	for _, file := range files {
		name := file.Name()
		if strings.HasSuffix(name, ".darwin-amd64.tar.gz") {
			versions = append(versions, gv.GetVersion(name))
		}
	}
	return versions
}
