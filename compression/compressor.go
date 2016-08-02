package compression

import "github.com/golang/snappy"

// Compressor compression interface
type Compressor interface {
	Compress(src []byte) []byte
	Decompress(src []byte) ([]byte, error)
}

var (
	compressor Compressor = snappyCompressor{}
)

// SetCompressor sets compressor
func SetCompressor(c Compressor) {
	compressor = c
}

// Compress compress []byte and returns []byte
// with compressed data
func Compress(src []byte) []byte {
	return compressor.Compress(src)
}

// Decompress []byte and returns []byte
// with decompressed data
func Decompress(src []byte) ([]byte, error) {
	return compressor.Decompress(src)
}

type snappyCompressor struct{}

func (snappyCompressor) Compress(src []byte) []byte {
	return snappy.Encode(nil, src)
}

func (snappyCompressor) Decompress(src []byte) ([]byte, error) {
	return snappy.Decode(nil, src)
}

// NewSnappyCompressor returns new snappyCompressor
func NewSnappyCompressor() snappyCompressor {
	return snappyCompressor{}
}
