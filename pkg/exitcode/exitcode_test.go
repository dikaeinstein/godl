package exitcode

import "testing"

func TestExitCodeErrorImplementsError(t *testing.T) {
	var err interface{} = exitcodeError{}

	if _, ok := err.(Error); !ok {
		t.Fatalf("expected %t to implement Error", err)
	}
}
