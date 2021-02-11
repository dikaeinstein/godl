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
	"github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
)

// New returns the list command
func New() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List the downloaded versions.",
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := godlutil.GetDownloadDir()
			if err != nil {
				return err
			}

			return listDownloadedBinaryArchives(d)
		},
	}
}

func listDownloadedBinaryArchives(downloadDir string) error {
	// Create download directory and its parent
	godlutil.Must(os.MkdirAll(downloadDir, os.ModePerm))

	files, err := ioutil.ReadDir(downloadDir)
	if err != nil {
		return err
	}

	versions := mapToVersion(files)
	// sort in-place comparing version numbers
	sort.Sort(version.Collection(versions))

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
			versions = append(versions, godlutil.GetVersion(name))
		}
	}
	return versions
}
