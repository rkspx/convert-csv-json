package path

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"errors"
)

// Verify verify the safety of the given path and wether it is rooted in the given trustedRoot directory.
func Verify(path, trustedRoot string) (string, error) {
	r, err := filepath.EvalSymlinks(path)
	if err != nil {
		return path, fmt.Errorf("eval failed: %s (%T)", err, err)
	}

	if err := IsInTrustedRoot(r, trustedRoot); err != nil {
		return path, err
	}

	return r, nil
}

var (
	// ErrNotInTrusted is error returned when a path is not rooted in a trusted directory.
	ErrNotInTrusted = errors.New("not in trusted directory")
)

// IsInTrustedRoot checks wether the given path is rooted in a trusted directory,
// returns error if it doesn't.
func IsInTrustedRoot(path, trustedRoot string) error {
	for path != "/" && path != "." {
		path = filepath.Dir(path)
		if path == trustedRoot {
			return nil
		}
	}

	return ErrNotInTrusted
}

// SafelyOpenFile safely opens a file. SafelyOpenFile do a sanitization and trustedRoot checking against the given
// path first before opening the file.
func SafelyOpenFile(path, trustedRoot string) (io.ReadCloser, error) {
	var err error
	path, err = Verify(path, trustedRoot)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("can not open file, %s", err)
	}

	return f, nil
}

// SafelyCreateFile safely create a file. SafelyCreateFile do a sanitization and trustedRoot checking against the given
// path first before opening the file. WARNING: this overrides the file if it already exists!.
func SafelyCreateFile(path, trustedRoot string) (io.WriteCloser, error) {
	path = filepath.Clean(path)

	err := IsInTrustedRoot(path, trustedRoot)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("can not open file, %s", err)
	}

	return f, nil
}
