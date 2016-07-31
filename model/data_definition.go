package model

import "github.com/sparrowdb/db/engine"

// DataDefinition holds the stored item
type DataDefinition struct {
	Key   string
	Size  uint32
	Token string
	Ext   string
	Buf   []byte
}

// ToByteStream convert DataDefinition to ByteStream
func (df *DataDefinition) ToByteStream() *engine.ByteStream {
	byteStream := engine.NewByteStream(engine.LittleEndian)
	byteStream.PutString(df.Key)
	byteStream.PutString(df.Token)
	byteStream.PutUInt32(df.Size)
	byteStream.PutString(df.Ext)
	byteStream.PutBytes(df.Buf)
	return byteStream
}

// NewDataDefinitionFromByteStream convert ByteStream to DataDefinition
func NewDataDefinitionFromByteStream(bs *engine.ByteStream) *DataDefinition {
	df := DataDefinition{}
	df.Key = bs.GetString()
	df.Token = bs.GetString()
	df.Size = bs.GetUInt32()
	df.Ext = bs.GetString()
	df.Buf = bs.GetBytes()
	return &df
}
