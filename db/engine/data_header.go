package engine

import "unsafe"

// DataHeader holds the header of data file
// Index the offset of index
// BloomFilter the offset of bloomfilter
type DataHeader struct {
	Index       uint64
	BloomFilter uint64
}

// DataHeaderSize returns the total size
// of DataHeader object
func DataHeaderSize() int64 {
	var info DataHeader
	return int64(unsafe.Sizeof(info))
}

// ToByteStream convert DataHeader to ByteStream
func (dh *DataHeader) ToByteStream() *ByteStream {
	byteStream := NewByteStream(LittleEndian)
	byteStream.PutUInt64(dh.Index)
	byteStream.PutUInt64(dh.BloomFilter)
	return byteStream
}

// GetDataHeaderFromFile convert ByteStream to DataHeader
func GetDataHeaderFromFile(s *Storage) *DataHeader {
	bs, _ := s.Get(0)
	return &DataHeader{
		Index:       bs.GetUInt64(),
		BloomFilter: bs.GetUInt64(),
	}
}

// UpdateDataHeaderFile updates header in data file
func UpdateDataHeaderFile(s *Storage, d *DataHeader) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rw, _ := OpenRandomWriter(s.Filepath)

	bs := d.ToByteStream()
	buf := bs.Bytes()

	bout := NewByteStream(LittleEndian)
	bout.PutUInt32(uint32(len(buf)))

	if _, err := rw.AppendAt(bout.Bytes(), 0); err != nil {
		return err
	}
	if _, err := rw.AppendAt(buf, 4); err != nil {
		return err
	}

	if err := rw.Close(); err != nil {
		return err
	}

	return nil
}
