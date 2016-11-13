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
	"github.com/SparrowDb/sparrowdb/util"
)

type dataHolder struct {
	path        string
	sto         engine.Storage
	summary     index.Summary
	bloomfilter util.BloomFilter
}

func newDataHolder(sto *engine.Storage, dbPath string, bloomFilterFp float32) (*dataHolder, error) {
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
	dh := dataHolder{path: newPath}
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
	b := dh.bloomfilter.ByteStream()
	if err = writer.Append(b.Bytes()); err == nil {
		writer.Close()
	}

	return &dh, nil
}

func openDataHolder(path string) (*dataHolder, error) {
	var err error

	dh := dataHolder{path: path}

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
		dh.bloomfilter = *util.NewBloomFilterFromByteStream(bs)
	}

	return &dh, nil
}
