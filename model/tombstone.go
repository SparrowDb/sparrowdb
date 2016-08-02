package model

// NewTombstone returns new DataDefinition
// Tombstones are DataDefinition with Status = DataDefinitionRemoved
// and empty byte buffer containing the image data
func NewTombstone(df *DataDefinition) *DataDefinition {
	df.Status = DataDefinitionRemoved
	df.Buf = []byte("")
	return df
}
