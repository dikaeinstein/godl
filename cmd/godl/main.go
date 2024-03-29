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

	"github.com/dikaeinstein/godl/internal/app"
	"github.com/dikaeinstein/godl/internal/cli"
)

// injected as ldflags during go build
var (
	buildTime string
	gitTag    string
)

func main() {
	code := cli.Run(app.BuildInfo{
		BuildTime: buildTime,
		GitTag:    gitTag,
	})

	os.Exit(code)
}
