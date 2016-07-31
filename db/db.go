package db

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/sparrowdb/db/cache"
	"github.com/sparrowdb/db/engine"
	"github.com/sparrowdb/db/index"
	"github.com/sparrowdb/model"
	"github.com/sparrowdb/util"
)

var (
	dataFileFmt  = "db-%d.spw"
	indexFileFmt = "db-%d.idx"
)

// Database holds database definitions
type Database struct {
	Descriptor *DatabaseDescriptor
	commitlog  *Commitlog
	dh         []*dataHolder
	cache      *cache.Cache
	lock       sync.RWMutex
}

type dataHolder struct {
	storage *engine.Storage
	summary index.Summary
}

func newDataHolder(filepath string, summary index.Summary) {
	dh := dataHolder{}
	dh.storage = engine.NewStorage(filepath)
	dh.summary = summary

	idxOffset := uint64(dh.storage.GetSize())

	for _, v := range summary.GetTable() {
		dh.storage.Append(engine.NewByteStreamFromBytes(v.Bytes(), engine.LittleEndian))
	}

	bfOffset := uint64(dh.storage.GetSize())

	engine.UpdateDataHeaderFile(dh.storage, &engine.DataHeader{
		Index:       idxOffset,
		BloomFilter: bfOffset,
	})
}

func openDataHolder(filepath string) *dataHolder {
	dh := dataHolder{}
	dh.storage = engine.NewStorage(filepath)
	dh.summary = *index.NewSummary()

	dhHeader := engine.GetDataHeaderFromFile(dh.storage)
	offset := dhHeader.Index

	for offset < dhHeader.BloomFilter {
		bs, _ := dh.storage.Get(int64(offset))
		dh.summary.Add(index.NewEntryFromByteStream(bs))
		offset += uint64(len(bs.Bytes())) + 4
	}

	return &dh
}

// nextDataHolderFile returns the next file name
func nextDataHolderFile(filepath string) string {
	p, err := ioutil.ReadDir(filepath)
	if err != nil {
		log.Fatal(err)
	}

	last := 0

	if len(p) > 0 {
		t := len(p) - 1
		for i := t; i >= 0; i-- {
			if !p[i].IsDir() {
				last = i + 1
				break
			}
		}
	}

	return fmt.Sprintf(dataFileFmt, last)
}

// InsertData insert data into database
func (db *Database) InsertData(df *model.DataDefinition) {
	key := util.Hash32(df.Key)
	bs := df.ToByteStream()
	_ = key
	_ = bs

	if db.commitlog.Size()+uint64(df.Size) > db.Descriptor.MaxDataLogSize {
		db.lock.Lock()
		defer db.lock.Unlock()

		next := nextDataHolderFile(db.Descriptor.Path)
		ndh := filepath.Join(db.Descriptor.Path, next)

		cIndex := db.commitlog.GetSummary()

		db.commitlog.RenameTo(ndh)
		db.commitlog = NewCommitLog(db.Descriptor.Path)

		newDataHolder(ndh, cIndex)
	} else {
		err := db.commitlog.Add(key, bs)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	// Put in cache
	db.cache.Put(key, bs.Bytes())
}

// GetDataByKey returns pointer to DataDefinition and bool if found the data
func (db *Database) GetDataByKey(key uint32) (*model.DataDefinition, bool) {
	// Search for given key in cache
	if c := db.cache.Get(key); c != nil {
		bs := engine.NewByteStreamFromBytes(c, engine.LittleEndian)
		return model.NewDataDefinitionFromByteStream(bs), true
	}

	// Search in commitlog
	if bs := db.commitlog.Get(key); bs != nil {
		db.cache.Put(key, bs.Bytes())
		return model.NewDataDefinitionFromByteStream(bs), true
	}

	// Search in data files
	for _, d := range db.dh {
		if e, ok := d.summary.LookUp(key); ok == true {
			bs, _ := d.storage.Get(e.Offset)
			return model.NewDataDefinitionFromByteStream(bs), true
		}
	}

	return nil, false
}

// LoadData loads index and bloom filter from each data file
func (db *Database) LoadData() {
	flist, _ := ioutil.ReadDir(db.Descriptor.Path)
	for _, v := range flist {
		if m, _ := regexp.MatchString("^db\\-[0-9]+.spw$", v.Name()); m == true {
			fpath := filepath.Join(db.Descriptor.Path, v.Name())
			db.dh = append(db.dh, openDataHolder(fpath))
		}
	}
}

// NewDatabase returns new Database
func NewDatabase(descriptor *DatabaseDescriptor) *Database {
	db := Database{
		Descriptor: descriptor,
		commitlog:  NewCommitLog(descriptor.Path),
		cache:      cache.NewCache(cache.NewLRU(int64(descriptor.MaxCacheSize))),
	}

	return &db
}

// OpenDatabase returns oppened Database
func OpenDatabase(descriptor *DatabaseDescriptor) *Database {
	db := NewDatabase(descriptor)
	db.commitlog.LoadData()
	db.LoadData()
	return db
}
