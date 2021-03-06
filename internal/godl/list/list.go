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
)

// List lists the downloaded go versions
type List struct{}

func (l List) Run(downloadDir string, sortDirection gv.SortDirection) error {
	// Create download directory and its parent
	godlutil.Must(os.MkdirAll(downloadDir, os.ModePerm))

	files, err := ioutil.ReadDir(downloadDir)
	if err != nil {
		return err
	}

	versions := mapToVersion(files)
	// sort in-place comparing version numbers
	sort.Slice(versions, func(i, j int) bool {
		return gv.CompareVersions(versions[i], versions[j], sortDirection)
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
