package app

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	testCases := []struct {
		name           string
		description    string
		info           BuildInfo
		expectedOutput string
	}{
		{
			name:        "defaultVersion",
			description: "prints default version info",
			info:        BuildInfo{GoVersion: "go1.20.3"},
			expectedOutput: `Version: dev
Go version: go1.20.3
Git hash: none
Built: unknown
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := NewVersion(tc.info)
			out := new(bytes.Buffer)

			err := v.Run(out)
			require.NoError(t, err)
			require.Equal(t, tc.expectedOutput, out.String())
		})
	}
}
