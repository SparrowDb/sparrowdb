package index

import (
	"strconv"

	"github.com/sparrowdb/db/engine"
)

// Entry holds index entry
type Entry struct {
	Key    uint32
	Offset int64
	Active bool
}

// Bytes returns byte array with index entry data
func (e *Entry) Bytes() []byte {
	bs := engine.NewByteStream(engine.LittleEndian)
	bs.PutUInt32(e.Key)
	bs.PutUInt64(uint64(e.Offset))
	bs.PutString(strconv.FormatBool(e.Active))
	return bs.Bytes()
}

// NewEntryFromByteStream convert ByteStream to Entry
func NewEntryFromByteStream(bs *engine.ByteStream) *Entry {
	df := Entry{}
	df.Key = bs.GetUInt32()
	df.Offset = int64(bs.GetUInt64())

	active, _ := strconv.ParseBool(bs.GetString())
	df.Active = active
	return &df
}
