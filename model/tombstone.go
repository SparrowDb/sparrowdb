package model

// NewTombstone returns new DataDefinition
// Tombstones are DataDefinition with Status = 2
// and empty byte buffer containing the image data
func NewTombstone(df *DataDefinition) *DataDefinition {
	df.Status = 2
	df.Buf = []byte("")
	return df
}
