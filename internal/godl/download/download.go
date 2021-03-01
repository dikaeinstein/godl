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

package download

import (
	"context"
	"fmt"
	"time"

	"github.com/dikaeinstein/godl/internal/pkg/downloader"
)

// Download downloads go binaries
type Download struct {
	Dl      *downloader.Downloader
	Timeout time.Duration
}

// Run downloads the specified go version
func (d *Download) Run(ctx context.Context, version string) error {
	fmt.Printf("Downloading go archive %v\n", version)

	ctx, cancel := context.WithTimeout(ctx, d.Timeout)
	defer cancel()
	err := d.Dl.Download(ctx, version)
	if err != nil {
		return fmt.Errorf("error downloading %v: %v", version, err)
	}

	fmt.Println("\nDownload complete")
	return nil
}
