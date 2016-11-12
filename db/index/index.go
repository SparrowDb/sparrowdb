package index

import "github.com/SparrowDb/sparrowdb/util"

// Entry holds index entry
type Entry struct {
	Key      uint32
	Offset   int64
	Status   uint16
	Revision uint32
}

// Bytes returns byte array with index entry data
func (e *Entry) Bytes() []byte {
	bs := util.NewByteStream()
	bs.PutUInt32(e.Key)
	bs.PutUInt64(uint64(e.Offset))
	bs.PutUInt16(e.Status)
	bs.PutUInt32(e.Revision)
	return bs.Bytes()
}

// NewEntryFromByteStream convert ByteStream to Entry
func NewEntryFromByteStream(bs *util.ByteStream) *Entry {
	df := Entry{}
	df.Key = bs.GetUInt32()
	df.Offset = int64(bs.GetUInt64())
	df.Status = bs.GetUInt16()
	df.Revision = bs.GetUInt32()
	return &df
}
