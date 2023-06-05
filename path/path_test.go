package path

import (
	"testing"
)

func TestPath(t *testing.T) {
	const trustedRoot = "testdata"
	c, err := Verify("/etc/hosts", trustedRoot)
	if err == nil {
		t.Fatal("expecting error")
	}

	c, err = Verify("data.csv", trustedRoot)
	if err == nil {
		t.Fatal("expecting error")
	}

	c, err = Verify("testdata/foo.txt", trustedRoot)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	t.Logf("safe path: %s", c)
}
