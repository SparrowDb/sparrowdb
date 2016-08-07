package model

import (
	"github.com/sparrowdb/compression"
	"github.com/sparrowdb/util"
	"github.com/sparrowdb/util/uuid"
)

const (
	// DataDefinitionActive active status
	DataDefinitionActive = 1

	// DataDefinitionRemoved removed status
	DataDefinitionRemoved = 2
)

// DataDefinition holds the stored item
type DataDefinition struct {
	Key    string
	Size   uint32
	Token  string
	Ext    string
	Status uint16
	Buf    []byte
}

// DataDefinitionResult holds DataDefinition query result
type DataDefinitionResult struct {
	Key       string
	Size      uint32
	Token     string
	TImestamp string
	Ext       string
	Status    string
}

// QueryResult convert DataDefinition to DataDefinitionResult
func (df *DataDefinition) QueryResult() *DataDefinitionResult {
	dfr := DataDefinitionResult{
		Key:   df.Key,
		Size:  df.Size,
		Token: df.Token,
		Ext:   df.Ext,
	}

	switch df.Status {
	case 1:
		dfr.Status = "Active"
	case 2:
		dfr.Status = "Removed"
	}

	uuid, _ := uuid.ParseUUID(df.Token)
	dfr.TImestamp = uuid.Time().String()

	return &dfr
}

// ToByteStream convert DataDefinition to ByteStream
func (df *DataDefinition) ToByteStream() *util.ByteStream {
	byteStream := util.NewByteStream()
	byteStream.PutString(df.Key)
	byteStream.PutString(df.Token)
	byteStream.PutUInt32(df.Size)
	byteStream.PutString(df.Ext)
	byteStream.PutUInt16(df.Status)

	encoded := compression.Compress(df.Buf)
	byteStream.PutBytes(encoded)

	return byteStream
}

// NewDataDefinitionFromByteStream convert ByteStream to DataDefinition
func NewDataDefinitionFromByteStream(bs *util.ByteStream) *DataDefinition {
	df := DataDefinition{}
	df.Key = bs.GetString()
	df.Token = bs.GetString()
	df.Size = bs.GetUInt32()
	df.Ext = bs.GetString()
	df.Status = bs.GetUInt16()

	buf := bs.GetBytes()
	if decoded, err := compression.Decompress(buf); err == nil {
		df.Buf = decoded
	}

	return &df
}
