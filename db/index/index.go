package index

import "github.com/SparrowDb/sparrowdb/util"

// Entry holds index entry
type Entry struct {
	Key      uint32
	Offset   int64
	Status   uint16
	Revision uint32
	Version  []uint32
}

// Bytes returns byte array with index entry data
func (e *Entry) Bytes() []byte {
	bs := util.NewByteStream()
	bs.PutUInt32(e.Key)
	bs.PutUInt64(uint64(e.Offset))
	bs.PutUInt16(e.Status)
	bs.PutUInt32(e.Revision)

	var idx uint32
	for idx = 0; idx < e.Revision; idx++ {
		bs.PutUInt32(e.Version[idx])
	}

	return bs.Bytes()
}

// NewEntryFromByteStream convert ByteStream to Entry
func NewEntryFromByteStream(bs *util.ByteStream) *Entry {
	df := Entry{}
	df.Key = bs.GetUInt32()
	df.Offset = int64(bs.GetUInt64())
	df.Status = bs.GetUInt16()
	df.Revision = bs.GetUInt32()

	df.Version = make([]uint32, 0)
	for idx := 0; idx < int(df.Revision); idx++ {
		df.Version = append(df.Version, bs.GetUInt32())
	}

	return &df
}
