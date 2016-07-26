package engine

import (
	"fmt"
	"os"
	"sync"
)

var (
	// ErrReadFile error message when read file
	ErrReadFile = "Error trying to read file"
)

// Reader holds writer
type Reader struct {
	fp   *os.File
	lock sync.RWMutex
}

// ReadAt reads file at offset position
func (r *Reader) ReadAt(offset int64, b []byte) error {
	read, err := r.fp.ReadAt(b, offset)
	if err != nil {
		return err
	}

	if read != len(b) {
		return fmt.Errorf("%s: %s", ErrReadFile, r.fp.Name())
	}

	return nil
}

// Read reads file at offset position
func (r *Reader) Read(b []byte) error {
	read, err := r.fp.Read(b)
	if err != nil {
		return err
	}

	if read != len(b) {
		return fmt.Errorf("%s: %s", ErrReadFile, r.fp.Name())
	}

	return nil
}

// Close closes current file
func (r *Reader) Close() error {
	if err := r.fp.Close(); err != nil {
		return err
	}
	return nil
}

// OpenReader opens file reader
func OpenReader(filepath string) (*Reader, error) {
	r := &Reader{}
	f, err := OpenFile(filepath, os.O_RDONLY|os.O_CREATE)
	if err != nil {
		return nil, err
	}
	r.fp = f
	return r, nil
}
