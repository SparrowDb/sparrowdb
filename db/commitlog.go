package db

import (
	"log"
	"os"
	"path/filepath"

	"github.com/sparrowdb/db/engine"
	"github.com/sparrowdb/db/index"
	"github.com/sparrowdb/db/iterator"
	"github.com/sparrowdb/util"
)

var (
	commitlogFileFmt = "commitlog.spw"
)

// Commitlog holds commitlog information
type Commitlog struct {
	filepath string
	sto      *engine.Storage
	summary  *index.Summary
}

// Get returns ByteStream with requested data, nil if not found
func (c *Commitlog) Get(key uint32) *engine.ByteStream {
	// Search in index if found, get from data file
	if idx, ok := c.summary.LookUp(key); ok == true {
		if bs, err := c.sto.Get(idx.Offset); err == nil {
			return bs
		}
	}
	return nil
}

// Add add entry to commitlog
func (c *Commitlog) Add(key uint32, bs *engine.ByteStream) error {
	pos := c.sto.GetSize()
	if err := c.sto.Append(bs); err != nil {
		return err
	}
	c.summary.Add(&index.Entry{
		Key:    key,
		Offset: pos,
	})
	return nil
}

// LoadData loads commitlog data file
func (c *Commitlog) LoadData() {
	iter, err := iterator.NewDataIterator(c.filepath)
	if err != nil {
		log.Fatal(err)
	}

	for df, h, _ := iterator.Iterate(iter); h == true; df, h, _ = iter.Next() {
		c.summary.Add(&index.Entry{
			Key:    util.Hash32(df.Key),
			Offset: iter.GetOffset(),
		})
	}
}

// NewCommitLog returns new Commitlog
func NewCommitLog(path string) *Commitlog {
	fpath := filepath.Join(path, commitlogFileFmt)

	if _, err := os.Stat(fpath); err != nil {
		util.CreateEmptyFile(fpath)
	}

	c := Commitlog{}
	c.filepath = fpath
	c.sto = engine.NewStorage(fpath)
	c.summary = index.NewSummary()
	return &c
}
