package model

import "github.com/sparrowdb/db/engine"

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

// ToByteStream convert DataDefinition to ByteStream
func (df *DataDefinition) ToByteStream() *engine.ByteStream {
	byteStream := engine.NewByteStream(engine.LittleEndian)
	byteStream.PutString(df.Key)
	byteStream.PutString(df.Token)
	byteStream.PutUInt32(df.Size)
	byteStream.PutString(df.Ext)
	byteStream.PutUInt16(df.Status)
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
	df.Status = bs.GetUInt16()
	df.Buf = bs.GetBytes()
	return &df
}
