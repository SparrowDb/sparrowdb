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
	engine.Storage
	summary *index.Summary
}

// GetIndex returns offset in data file of the key
func (c *Commitlog) GetIndex(key uint32) (*index.Entry, bool) {
	v, ok := c.summary.LookUp(key)
	return v, ok
}

// Add add entry to commitlog
func (c *Commitlog) Add(key uint32, bs *engine.ByteStream) error {
	pos := c.GetSize()
	if err := c.Append(bs); err != nil {
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
	iter, err := iterator.NewDataIterator(c.Filepath)
	if err != nil {
		log.Fatal(err)
	}

	df, h, _ := iterator.Iterate(iter)
	for h == true {
		c.summary.Add(&index.Entry{
			Key:    util.Hash32(df.Key),
			Offset: iter.GetOffset(),
		})
		df, h, _ = iter.Next()
	}
}

// NewCommitLog returns new Commitlog
func NewCommitLog(path string) *Commitlog {
	fpath := filepath.Join(path, commitlogFileFmt)

	if _, err := os.Stat(fpath); err != nil {
		log.Printf("%s not found, creating", path)
		util.CreateEmptyFile(fpath)
	}

	c := Commitlog{}
	c.Filepath = fpath
	c.summary = index.NewSummary()
	return &c
}
