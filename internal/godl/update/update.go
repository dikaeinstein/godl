package update

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-version"
	"github.com/spf13/viper"
)

type Asset struct {
	Name string `json:"name"`
}
type Release struct {
	Assets  []Asset `json:"assets"`
	TagName string  `json:"tag_name"`
}
type ListReleasesResult []Release
type ListReleasesErrorResp struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}

// Update checks for if there are updates available for Godl
type Update struct {
	Client *http.Client
	Output io.Writer
}

func (u *Update) Run(ctx context.Context, currentVersion string) error {
	exists, latest, err := u.CheckForUpdate(ctx, currentVersion)
	if err != nil {
		return err
	}

	if exists {
		fmt.Fprintf(u.Output, `Your version of Godl is out of date! The latest version
 is %s. You can update by downloading from https://github.com/dikaeinstein/godl/releases
`, latest.TagName)
	} else {
		fmt.Fprintln(u.Output, "No update available.")
	}

	return nil
}

func (u *Update) CheckForUpdate(ctx context.Context, currentVersion string) (bool, *Release, error) {
	url := "https://api.github.com/repos/dikaeinstein/godl/releases?per_page=10"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "godl")

	ghToken := viper.GetString("gh_token")
	if ghToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", ghToken))
	}

	res, err := u.Client.Do(req)
	if err != nil {
		return false, nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotFound {
		var errResp ListReleasesErrorResp
		// using io.ReadAll because res.Body is very small
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return false, nil, err
		}
		err = json.Unmarshal(data, &errResp)
		if err != nil {
			return false, nil, err
		}
		return false, nil, fmt.Errorf("%s: %v %v", url, res.StatusCode, errResp.Message)
	}

	var releases ListReleasesResult
	err = json.NewDecoder(res.Body).Decode(&releases)
	if err != nil {
		return false, nil, err
	}

	// Only a single version exists or no version :)
	const minNumOfRelease = 2
	if len(releases) < minNumOfRelease {
		return false, nil, nil
	}

	// pick latest release
	r := releases[0]
	latest := version.Must(version.NewVersion(r.TagName))
	current := version.Must(version.NewVersion(currentVersion))
	return latest.GreaterThan(current), &r, nil
}