package hash

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// RemoteHasher generates a hash from a remote source
type RemoteHasher struct {
	client *http.Client
}

// Hash fetches the hash of the given URL and returns it as a string.
func (r RemoteHasher) Hash(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return "", err
	}

	res, err := r.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s: %v", url, res.Status)
	}

	urlHash, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("reading %s: %v", url, err)
	}

	return string(urlHash), nil
}

// NewRemoteHasher returns an initialized RemoteHasher
func NewRemoteHasher(client *http.Client) RemoteHasher {
	return RemoteHasher{client}
}

type FakeHasher struct{}

func (FakeHasher) Hash(ctx context.Context, path string) (string, error) {
	return "fakehash", nil
}
