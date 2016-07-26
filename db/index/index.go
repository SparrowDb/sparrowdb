package index

import "github.com/sparrowdb/db/engine"

// Entry holds index entry
type Entry struct {
	Key    uint32
	Offset int64
	Active byte
}

// Bytes returns byte array with index entry data
func (e *Entry) Bytes() []byte {
	bs := engine.NewByteStream(engine.LittleEndian)
	bs.PutUInt32(e.Key)
	return bs.Bytes()
}
