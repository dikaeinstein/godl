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

package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	godlVersion = "unknown version"
	gitHash     = "unknown commit"
	goVersion   = "unknown go version"
	buildDate   = "unknown build date"
)

// New returns the version command
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show the godl version information.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\nGo version: %s\nGit hash: %s\nBuilt: %s\n",
				godlVersion, goVersion, gitHash, buildDate)
		},
	}
}
