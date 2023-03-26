package exitcode

import (
	"errors"
	"testing"
)

func TestExitCodeErrorImplementsError(t *testing.T) {
	var err interface{} = exitcodeError{}

	if _, ok := err.(Error); !ok {
		t.Fatalf("expected %t to implement Error", err)
	}
}

var errTest = errors.New("test error")

func TestNewError(t *testing.T) {
	testCases := []struct {
		err     error
		wantErr error
		desc    string
		code    int
	}{
		{
			desc:    "returns <nil> when given err is <nil>",
			err:     nil,
			code:    1,
			wantErr: nil,
		},
		{
			desc:    "returns an exitcodeError that wraps given err",
			err:     errTest,
			code:    1,
			wantErr: exitcodeError{errTest, 1},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			err := NewError(tC.err, tC.code)

			if !errors.Is(err, tC.wantErr) {
				t.Errorf("NewError failed: want %v, got %v", tC.wantErr, err)
			}
		})
	}
}
