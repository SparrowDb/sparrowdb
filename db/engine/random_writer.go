package engine

import (
	"fmt"
	"os"
	"sync"
)

// RandomWriter holds writer
type RandomWriter struct {
	fp *os.File
	mu sync.RWMutex
}

// AppendAt append bytes to current file
// starting from offset position
// returns written bytes and error message
func (w *RandomWriter) AppendAt(b []byte, offset int64) (int, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	written, err := w.fp.WriteAt(b, offset)
	if err != nil {
		return 0, err
	}

	if written != len(b) {
		return 0, fmt.Errorf("%s: %s", ErrWriteFile, w.fp.Name())
	}

	return written, nil
}

// Close closes current file
func (w *RandomWriter) Close() error {
	if err := w.fp.Sync(); err != nil {
		return err
	}
	if err := w.fp.Close(); err != nil {
		return err
	}
	return nil
}

// RandomWriter opens file writer
func OpenRandomWriter(filepath string) (*RandomWriter, error) {
	w := &RandomWriter{}
	f, err := OpenFile(filepath, os.O_WRONLY)
	if err != nil {
		return nil, err
	}
	w.fp = f
	return w, nil
}
