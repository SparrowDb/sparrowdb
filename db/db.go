package db

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/sparrowdb/cache"
	"github.com/sparrowdb/db/index"
	"github.com/sparrowdb/engine"
	"github.com/sparrowdb/model"
	"github.com/sparrowdb/slog"
	"github.com/sparrowdb/util"
)

// Database holds database definitions
type Database struct {
	Descriptor *DatabaseDescriptor
	commitlog  *Commitlog
	dhList     []dataHolder
	cache      *cache.Cache
	mu         sync.RWMutex
}

type dataHolder struct {
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
	dh := dataHolder{}
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

	dh := dataHolder{}

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

func (d *dataHolder) Get(position int64) (*util.ByteStream, error) {
	// Search in index if found, get from data file
	freader, _ := d.sto.Open(engine.FileDesc{Type: engine.FileData})
	r := newReader(freader.(io.ReaderAt))

	b, err := r.Read(position)
	if err != nil {
		slog.Fatalf(err.Error())
	}

	bs := util.NewByteStreamFromBytes(b)
	return bs, nil
}

// InsertData insert data into database
func (db *Database) InsertData(df *model.DataDefinition) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	hKey := util.Hash32(df.Key)
	bs := df.ToByteStream()

	// Put in cache
	db.cache.Put(hKey, bs.Bytes())

	// Get last position in commitlog
	size, err := db.commitlog.Size()
	if err != nil {
		return err
	}

	// Check if commitlog has the max file size
	if size+int64(df.Size) > int64(db.Descriptor.MaxDataLogSize) {
		ndh, err := newDataHolder(&db.commitlog.sto, db.Descriptor.Path, db.Descriptor.BloomFilterFp)
		if err != nil {
			return err
		}

		db.dhList = append(db.dhList, *ndh)
		db.commitlog = NewCommitLog(db.Descriptor.Path)
	}

	if err = db.commitlog.Add(df.Key, bs); err != nil {
		return err
	}

	return nil
}

// GetDataByKey returns pointer to DataDefinition and bool if found the data
func (db *Database) GetDataByKey(key string) (*model.DataDefinition, bool) {
	hkey := util.Hash32(key)

	// Search for given key in cache
	if c := db.cache.Get(hkey); c != nil {
		bs := util.NewByteStreamFromBytes(c)
		return model.NewDataDefinitionFromByteStream(bs), true
	}

	// Search in commitlog
	if bs := db.commitlog.Get(key); bs != nil {
		db.cache.Put(hkey, bs.Bytes())
		return model.NewDataDefinitionFromByteStream(bs), true
	}

	// Search in data files
	strKey := strconv.Itoa(int(hkey))
	for _, d := range db.dhList {
		if d.bloomfilter.Contains(strKey) {
			if e, eIdx := d.summary.LookUp(hkey); eIdx == true {
				bs, _ := d.Get(e.Offset)
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
		if m, _ := regexp.MatchString("^([0-9]{19})$", v.Name()); m == true {
			dh, err := openDataHolder(filepath.Join(db.Descriptor.Path, v.Name()))
			if err != nil {
				slog.Fatalf(err.Error())
			}
			db.dhList = append(db.dhList, *dh)
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
