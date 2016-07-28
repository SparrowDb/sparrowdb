package iterator

import (
	"os"

	"github.com/sparrowdb/db/engine"
	"github.com/sparrowdb/model"
)

// DataIterator holds iteration data
type DataIterator struct {
	storage  *engine.Storage
	filepath string

	offset  int64 // holds data offset in datafile
	current int64 // holds the position of the cursor
	fsize   int64
}

// Next returns next data and if has more data
func (di *DataIterator) Next() (*model.DataDefinition, bool, error) {
	return Iterate(di)
}

// GetOffset returns the current offset
func (di *DataIterator) GetOffset() int64 {
	return di.offset
}

// NewDataIterator returns new DataIterator
func NewDataIterator(filepath string) (*DataIterator, error) {
	stat, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}

	// When data is written to disk, it always skip 4 bytes
	// to keep the length of next data block
	start := engine.DataHeaderSize() + 4

	return &DataIterator{
		filepath: filepath,
		current:  start,
		offset:   start,
		storage:  engine.NewStorage(filepath),
		fsize:    stat.Size(),
	}, nil
}
