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

	iter.offset = iter.current
	iter.current = iter.current + int64(len(bs.Bytes())) + 4

	return model.NewDataDefinitionFromByteStream(bs), true, nil
}
