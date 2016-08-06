package engine

import (
	"fmt"
	"os"
)

var (
	ErrNotAFile = "Not a file"
)

type fileLock interface {
	release() error
}

// OpenFile opens file and returns file
func OpenFile(filepath string, flag int) (*os.File, error) {
	if fi, err := os.Stat(filepath); err == nil {
		if fi.IsDir() {
			return nil, fmt.Errorf("%s: %s", ErrNotAFile, filepath)
		}
	}

	f, err := os.OpenFile(filepath, flag, 0644)
	if err != nil {
		return nil, err
	}

	return f, nil
}
