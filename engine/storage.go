package engine

import (
	"fmt"
	"io"
)

// FileType Describes wich type is the file
type FileType int

const (
	FileData FileType = iota

	FileIndex

	FileBloomFilter

	FileCommitlog
)

// FileDesc is the file descriptor
type FileDesc struct {
	Type FileType
}

// Name returns the name of the file based on type
func (fd *FileDesc) Name() string {
	switch fd.Type {
	case FileData:
		return fmt.Sprintf("data.spw")
	case FileIndex:
		return fmt.Sprintf("index.spw")
	case FileBloomFilter:
		return fmt.Sprintf("bloom.spw")
	case FileCommitlog:
		return fmt.Sprintf("commitlog.spw")
	default:
		return ""
	}
}

// Reader interface to file reader
type Reader interface {
	io.ReadSeeker
	io.ReaderAt
	io.Closer
}

// Writer interface to file writer
type Writer interface {
	io.WriteCloser
}

// Storage interface to manage the storage system
type Storage interface {
	Open(fd FileDesc) (Reader, error)

	Create(fd FileDesc) (Writer, error)

	Size(fd FileDesc) (int64, error)

	Exists(fd FileDesc) bool

	Remove(fd FileDesc) error

	Rename(ofd, nfd FileDesc) error

	Close() error
}
