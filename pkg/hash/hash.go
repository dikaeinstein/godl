package hash

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// RemoteHash generates a hash from a remote source
type RemoteHasher struct {
	client *http.Client
}

// Hash fetches the hash of the given URL and returns it as a string.
func (r RemoteHasher) Hash(url string) (string, error) {
	res, err := r.client.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s: %v", url, res.Status)
	}

	urlHash, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("reading %s: %v", url, err)
	}

	return string(urlHash), nil
}

type FakeHasher struct{}

func (FakeHasher) Hash(path string) (string, error) {
	return "fakehash", nil
}
