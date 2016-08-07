package db

import (
	"io"

	"github.com/sparrowdb/db/index"
	"github.com/sparrowdb/engine"
	"github.com/sparrowdb/slog"
	"github.com/sparrowdb/util"
)

type dbReader struct {
	reader io.ReaderAt
	offset uint64
}

func (r *dbReader) Read(offset int64) ([]byte, error) {
	bSize := make([]byte, 4)
	if _, err := r.reader.ReadAt(bSize, offset); err != nil {
		slog.Fatalf(err.Error())
	}

	bs := util.NewByteStreamFromBytes(bSize)
	size := bs.GetUInt32()

	// Skip 4 bytes of the size mark
	offset += 4

	// Reads data
	bufData := make([]byte, size)
	if _, err := r.reader.ReadAt(bufData, offset); err != nil {
		slog.Fatalf(err.Error())
	}

	if err := r.Close(); err != nil {
		return nil, err
	}

	return bufData, nil
}

func (r *dbReader) Close() error {
	return nil
}

func newReader(f io.ReaderAt) *dbReader {
	return &dbReader{f, 0}
}

type indexReader struct {
	sto *engine.Storage
}

func (ir *indexReader) LoadIndex() (index.Summary, error) {
	summary := index.NewSummary()

	desc := engine.FileDesc{Type: engine.FileIndex}
	var pos int64
	var s = (*ir.sto)

	size, err := s.Size(desc)
	if err != nil {
		slog.Fatalf(err.Error())
	}

	freader, err := s.Open(desc)
	if err != nil {
		slog.Fatalf(err.Error())
	}

	r := newReader(freader.(io.ReaderAt))

	for pos < size {
		if b, err := r.Read(pos); err == nil {
			bs := util.NewByteStreamFromBytes(b)
			summary.Add(index.NewEntryFromByteStream(bs))
			pos += int64(bs.Size()) + 4
		} else {
			slog.Fatalf(err.Error())
		}
	}

	return *summary, nil
}

func newIndexReader(sto *engine.Storage) *indexReader {
	return &indexReader{sto}
}
