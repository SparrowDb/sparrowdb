package db

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/SparrowDb/sparrowdb/db/index"
	"github.com/SparrowDb/sparrowdb/engine"
	"github.com/SparrowDb/sparrowdb/errors"
	"github.com/SparrowDb/sparrowdb/slog"
	"github.com/SparrowDb/sparrowdb/util"
)

// DataHolder definitive data file after commitlog flush
type DataHolder struct {
	path        string
	sto         engine.Storage
	summary     index.Summary
	bloomfilter util.BloomFilter
}

// Get get ByteStream from dataholder for a given position in data file
func (d *DataHolder) Get(position int64) (*util.ByteStream, error) {
	// Search in index if found, get from data file
	freader, err := d.sto.Open(engine.FileDesc{Type: engine.FileData})
	if err != nil {
		slog.Errorf(errors.ErrFileCorrupted.Error(), d.path)
		return nil, nil
	}

	r := newReader(freader.(io.ReaderAt))

	// If found key but can't load it from file, it will return nil to avoid
	// db crash. Returning nil will send to user empty query result
	b, err := r.Read(position)
	if err != nil {
		slog.Errorf(errors.ErrFileCorrupted.Error(), d.path)
		return nil, nil
	}

	bs := util.NewByteStreamFromBytes(b)
	return bs, nil
}

// GetSummary get index summary of current data gile
func (d *DataHolder) GetSummary() index.Summary {
	return d.summary
}

// NewDataHolder returns new DataHolder pointer
func NewDataHolder(sto *engine.Storage, dbPath string, bloomFilterFp float32) (*DataHolder, error) {
	var err error

	// commitlog full path
	cPath := filepath.Join(dbPath, FolderCommitlog)

	// new name for commitlog folder
	uTime := fmt.Sprintf("%v", time.Now().UnixNano())
	newPath := filepath.Join(dbPath, uTime)

	// Rename commitlog file to data file
	if err := (*sto).Rename(engine.FileDesc{Type: engine.FileCommitlog}, engine.FileDesc{Type: engine.FileData}); err != nil {
		return nil, err
	}

	// Rename directory to unix time
	if err := os.Rename(cPath, newPath); err != nil {
		return nil, err
	}

	// Load dataholder
	dh := DataHolder{path: newPath}
	if dh.sto, err = engine.OpenFile(newPath); err != nil {
		return nil, err
	}

	// Load index from dataholder
	ir := newIndexReader(&dh.sto)
	dh.summary, err = ir.LoadIndex()
	if err != nil {
		return nil, err
	}

	// Create and populate bloomfilter
	table := dh.summary.GetTable()
	dh.bloomfilter = util.NewBloomFilter(dh.summary.Count(), bloomFilterFp)
	for _, v := range table {
		dh.bloomfilter.Add(strconv.Itoa(int(v.Key)))
	}

	bfw, err := dh.sto.Create(engine.FileDesc{Type: engine.FileBloomFilter})
	if err != nil {
		return nil, err
	}

	writer := newBufWriter(bfw)
	b, errbs := dh.bloomfilter.ByteStream()
	if errbs != nil {
		return nil, errbs
	}

	if err = writer.Append(b.Bytes()); err == nil {
		writer.Close()
	}

	return &dh, nil
}

// OpenDataHolder opens data holder for a given path
func OpenDataHolder(path string) (*DataHolder, error) {
	var err error

	dh := DataHolder{path: path}

	dh.sto, err = engine.OpenFile(path)
	if err != nil {
		return nil, err
	}

	// Loads index
	ir := newIndexReader(&dh.sto)
	dh.summary, err = ir.LoadIndex()
	if err != nil {
		return nil, err
	}

	// Loads bloomfilter
	var pos int64
	var bfreader io.Reader

	bfreader, err = dh.sto.Open(engine.FileDesc{Type: engine.FileBloomFilter})
	if err != nil {
		return nil, err
	}

	r := newReader(bfreader.(io.ReaderAt))

	if b, err := r.Read(pos); err == nil {
		bs := util.NewByteStreamFromBytes(b)
		dh.bloomfilter, err = util.NewBloomFilterFromByteStream(bs)
		if err != nil {
			return nil, err
		}
	}

	return &dh, nil
}
