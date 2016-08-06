package engine

import (
	"fmt"
	"os"
	"sync"
)

var (
	// ErrWriteFile holds error message when write file
	ErrWriteFile = "Error trying to write file"
)

// Writer holds writer
type Writer struct {
	fp *os.File
	mu sync.RWMutex
}

// Append append bytes to current file
// returns written bytes and error message
func (w *Writer) Append(b []byte) (int, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	written, err := w.fp.Write(b)
	if err != nil {
		return 0, err
	}

	if written != len(b) {
		return 0, fmt.Errorf("%s: %s", ErrWriteFile, w.fp.Name())
	}

	return written, nil
}

// Truncate truncates current file
func (w *Writer) Truncate(size int64) error {
	return w.fp.Truncate(size)
}

// Close closes current file
func (w *Writer) Close() error {
	if err := w.fp.Sync(); err != nil {
		return err
	}
	if err := w.fp.Close(); err != nil {
		return err
	}
	return nil
}

// OpenWriter opens file writer
func OpenWriter(filepath string) (*Writer, error) {
	w := &Writer{}
	f, err := OpenFile(filepath, os.O_WRONLY|os.O_APPEND|os.O_CREATE)
	if err != nil {
		return nil, err
	}
	w.fp = f
	return w, nil
}
