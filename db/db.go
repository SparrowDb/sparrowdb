package db

import (
	"log"
	"strconv"

	"github.com/sparrowdb/db/cache"
	"github.com/sparrowdb/db/engine"
	"github.com/sparrowdb/db/index"
	"github.com/sparrowdb/model"
	"github.com/sparrowdb/util"
)

// Database holds database definitions
type Database struct {
	Descriptor *DatabaseDescriptor
	commitlog  *Commitlog
	dh         []*DataHolder
	cache      *cache.Cache
}

type DataHolder struct {
	storage *engine.Storage
	summary *index.Summary
}

// InsertData insert data into database
func (db *Database) InsertData(df *model.DataDefinition) {
	key := util.Hash32(df.Key)
	bs := df.ToByteStream()

	err := db.commitlog.Add(key, bs)
	if err != nil {
		log.Fatalf(err.Error())
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

	if bs := db.commitlog.Get(key); bs != nil {
		db.cache.Put(key, bs.Bytes())
		return model.NewDataDefinitionFromByteStream(bs), true
	}

	return nil, false
}

// NewDatabase returns new Database
func NewDatabase(descriptor *DatabaseDescriptor) *Database {
	cacheSize, _ := strconv.ParseInt(descriptor.MaxCacheSize, 10, 64)

	db := Database{
		Descriptor: descriptor,
		commitlog:  NewCommitLog(descriptor.Path),
		cache:      cache.NewCache(cache.NewLRU(int64(cacheSize))),
	}

	return &db
}

// OpenDatabase returns oppened Database
func OpenDatabase(descriptor *DatabaseDescriptor) *Database {
	db := NewDatabase(descriptor)
	db.commitlog.LoadData()
	return db
}
