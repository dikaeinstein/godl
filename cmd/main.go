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

package main

import (
	"os"

	"github.com/dikaeinstein/godl/internal/cli"
)

// injected as ldflags during go build
var (
	buildDate   = "unknown build date"
	gitHash     = "unknown commit"
	godlVersion = "unknown version"
	goVersion   = "unknown go version"
)

func main() {
	opt := cli.Option{}

	opt.Version.BuildDate = buildDate
	opt.Version.GitHash = gitHash
	opt.Version.GodlVersion = godlVersion
	opt.Version.GoVersion = goVersion

	os.Exit(cli.Run(opt))
}
