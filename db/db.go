package db

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"

	"github.com/sparrowdb/db/cache"
	"github.com/sparrowdb/db/engine"
	"github.com/sparrowdb/db/index"
	"github.com/sparrowdb/model"
	"github.com/sparrowdb/slog"
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
	storage     *engine.Storage
	summary     index.Summary
	bloomfilter util.BloomFilter
}

func newDataHolder(filepath string, bloomFilterFp float32, summary index.Summary) {
	dh := dataHolder{}
	dh.storage = engine.NewStorage(filepath)
	dh.summary = summary

	dh.bloomfilter = util.NewBloomFilter(summary.Count(), bloomFilterFp)

	// Append index
	idxOffset := uint64(dh.storage.GetSize())
	for _, v := range summary.GetTable() {
		dh.bloomfilter.Add(strconv.Itoa(int(v.Key)))
		dh.storage.Append(engine.NewByteStreamFromBytes(v.Bytes(), engine.LittleEndian))
	}

	// Append bloomfilter
	bfOffset := uint64(dh.storage.GetSize())
	dh.storage.Append(dh.bloomfilter.ByteStream())

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

	offset = dhHeader.BloomFilter
	fsize := uint64(dh.storage.GetSize())
	for offset < fsize {
		bs, _ := dh.storage.Get(int64(offset))
		dh.bloomfilter = *util.NewBloomFilterFromByteStream(bs)
		offset += uint64(len(bs.Bytes())) + 4
	}
	return &dh
}

// nextDataHolderFile returns the next file name
func nextDataHolderFile(filepath string) string {
	p, err := ioutil.ReadDir(filepath)
	if err != nil {
		slog.Fatalf(err.Error())
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
func (db *Database) InsertData(df *model.DataDefinition) error {
	key := util.Hash32(df.Key)
	bs := df.ToByteStream()

	// check if DataDefinition will be greater than MaxDataLogSize
	if db.commitlog.Size()+uint64(df.Size) > db.Descriptor.MaxDataLogSize {
		db.lock.Lock()
		defer db.lock.Unlock()

		// get next data file name
		next := nextDataHolderFile(db.Descriptor.Path)
		ndh := filepath.Join(db.Descriptor.Path, next)

		// copy index
		cIndex := db.commitlog.GetSummary()

		// recreate an empty commitlog
		db.commitlog.RenameTo(ndh)
		db.commitlog = NewCommitLog(db.Descriptor.Path)

		// create new data holder file
		newDataHolder(ndh, db.Descriptor.BloomFilterFp, cIndex)
	} else {
		err := db.commitlog.Add(key, bs)
		if err != nil {
			return err
		}
	}

	// Put in cache
	db.cache.Put(key, bs.Bytes())

	return nil
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
	strKey := strconv.Itoa(int(key))
	for _, d := range db.dh {
		if eBF := d.bloomfilter.Contains(strKey); eBF == true {
			if e, eIdx := d.summary.LookUp(key); eIdx == true {
				bs, _ := d.storage.Get(e.Offset)
				db.cache.Put(key, bs.Bytes())
				return model.NewDataDefinitionFromByteStream(bs), true
			}
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
