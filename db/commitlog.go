package db

import (
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/sparrowdb/db/index"
	"github.com/sparrowdb/engine"
	"github.com/sparrowdb/slog"
	"github.com/sparrowdb/util"
)

// Commitlog holds commitlog information
type Commitlog struct {
	filepath string
	sto      engine.Storage
	summary  *index.Summary
	mu       sync.RWMutex
	desc     engine.FileDesc
}

// Get returns ByteStream with requested data, nil if not found
func (c *Commitlog) Get(key string) *util.ByteStream {
	// Search in index if found, get from data file
	hKey := util.Hash32(key)
	if idx, ok := c.summary.LookUp(hKey); ok == true {
		freader, _ := c.sto.Open(c.desc)
		r := newReader(freader.(io.ReaderAt))

		b, err := r.Read(idx.Offset)
		if err != nil {
			slog.Fatalf(err.Error())
		}

		bs := util.NewByteStreamFromBytes(b)
		return bs
	}
	return nil
}

// Add add entry to commitlog
func (c *Commitlog) Add(key string, bs *util.ByteStream) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	fwriter, _ := c.sto.Create(c.desc)
	pos, _ := c.sto.Size(c.desc)

	writer := newWriter(fwriter)
	writer.Append(key, bs.Bytes())

	writer.Close()

	hKey := util.Hash32(key)

	c.writeIndex(&index.Entry{
		Key:    hKey,
		Offset: pos,
	})

	return nil
}

func (c *Commitlog) writeIndex(index *index.Entry) error {
	var err error

	fwriter, err := c.sto.Create(engine.FileDesc{Type: engine.FileIndex})

	writer := newBufWriter(fwriter)
	if err = writer.Append(index.Bytes()); err == nil {
		c.summary.Add(index)
		writer.Close()
	}

	return err
}

// LoadData loads commitlog data file
func (c *Commitlog) LoadData() {
	desc := engine.FileDesc{Type: engine.FileIndex}
	var pos int64

	if !c.sto.Exists(desc) {
		return
	}

	size, err := c.sto.Size(desc)
	if err != nil {
		slog.Fatalf(err.Error())
	}

	freader, err := c.sto.Open(desc)
	if err != nil {
		slog.Fatalf(err.Error())
	}

	r := newReader(freader.(io.ReaderAt))

	for pos < size {
		if b, err := r.Read(pos); err == nil {
			bs := util.NewByteStreamFromBytes(b)
			c.summary.Add(index.NewEntryFromByteStream(bs))
			pos += int64(bs.Size()) + 4
		} else {
			slog.Fatalf(err.Error())
		}
	}
}

// Size returns commitlog file size
func (c *Commitlog) Size() (int64, error) {
	return c.sto.Size(c.desc)
}

// GetSummary returns commitlog index
func (c *Commitlog) GetSummary() index.Summary {
	return *c.summary
}

// RenameTo rename commitlog file
func (c *Commitlog) RenameTo(newpath string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	os.Rename(c.filepath, newpath)
}

// NewCommitLog returns new Commitlog
func NewCommitLog(path string) *Commitlog {
	var err error

	c := Commitlog{}
	c.filepath = filepath.Join(path, "commitlog")
	c.summary = index.NewSummary()
	c.desc = engine.FileDesc{Type: engine.FileCommitlog}

	c.sto, err = engine.OpenFile(c.filepath)
	if err != nil {
		slog.Fatalf(err.Error())
	}

	return &c
}
