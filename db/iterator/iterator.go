package iterator

import "github.com/sparrowdb/model"

// Iterate iterate over data file
func Iterate(iter *DataIterator) (*model.DataDefinition, bool, error) {
	if iter.current == iter.fsize {
		return nil, false, nil
	}

	bs, err := iter.storage.Get(iter.current)
	if err != nil {
		return nil, false, err
	}

	df := model.DataDefinition{}
	df.Key = bs.GetString()
	df.Size = bs.GetUInt32()
	df.Ext = bs.GetString()
	df.Buf = bs.GetBytes()

	iter.offset = iter.current
	iter.current = iter.current + int64(len(bs.Bytes())) + 4

	return &df, true, nil
}
