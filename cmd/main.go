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
	"github.com/dikaeinstein/godl/internal/godl"
	"github.com/dikaeinstein/godl/internal/godl/completion"
	"github.com/dikaeinstein/godl/internal/godl/download"
	"github.com/dikaeinstein/godl/internal/godl/install"
	"github.com/dikaeinstein/godl/internal/godl/list"
	"github.com/dikaeinstein/godl/internal/godl/listremote"
	"github.com/dikaeinstein/godl/internal/godl/version"
	"github.com/spf13/cobra"
)

func main() {
	// root command
	godl := godl.New()
	// subcommands
	completion := completion.New(godl)
	download := download.New()
	install := install.New()
	list := list.New()
	listRemote := listremote.New()
	version := version.New()

	godl.RegisterSubCommands([]*cobra.Command{
		completion,
		download,
		install,
		list,
		listRemote,
		version,
	})

	godl.Execute()
}
