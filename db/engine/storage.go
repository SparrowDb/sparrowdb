package engine

import (
	"errors"
	"log"
	"os"
	"sync"
)

var (
	ErrReadDir = errors.New("Could not list directory")
)

// Storage holds storage information
type Storage struct {
	Filepath   string
	dataHeader *DataHeader
	lock       sync.RWMutex
}

// Append appends ByteStream to file
func (s *Storage) Append(bs *ByteStream) error {
	s.lock.RLock()
	defer s.lock.RUnlock()

	w, err := OpenWriter(s.Filepath)
	if err != nil {
		return err
	}

	buf := bs.Bytes()

	bout := NewByteStream(LittleEndian)
	bout.PutUInt32(uint32(len(buf)))

	if _, err := w.Append(bout.Bytes()); err != nil {
		return err
	}
	if _, err := w.Append(buf); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}
	return nil
}

// GetSize return the file size
func (s *Storage) GetSize() int64 {
	stat, _ := os.Stat(s.Filepath)
	return stat.Size()
}

// Get returns ByteStream from offset
func (s *Storage) Get(offset int64) (*ByteStream, error) {
	r, err := OpenReader(s.Filepath)
	if err != nil {
		return nil, err
	}

	// Reads first 4 bytes to know the DataDefinition size
	bSize := make([]byte, 4)
	if err := r.ReadAt(offset, bSize); err != nil {
		log.Fatal(err)
	}

	bs := NewByteStreamFromBytes(bSize, LittleEndian)
	size := bs.GetUInt32()

	// Skip 4 bytes of the size mark
	offset += 4

	// Reads data
	bufData := make([]byte, size)
	if err = r.ReadAt(offset, bufData); err != nil {
		log.Fatal(err)
	}

	if err := r.Close(); err != nil {
		return nil, err
	}

	return NewByteStreamFromBytes(bufData, LittleEndian), nil
}

// CheckHeader checks if the file has the header
// if not, write it
func (s *Storage) CheckHeader() {
	if f, err := os.Stat(s.Filepath); err == nil {
		if f.Size() == 0 {
			dh := &DataHeader{
				Index:       0,
				BloomFilter: 0,
			}
			s.Append(dh.ToByteStream())
		}
	}
}

// NewStorage returns new Storage passing full
// file path
func NewStorage(filepath string) *Storage {
	sto := &Storage{
		Filepath: filepath,
	}
	return sto
}
