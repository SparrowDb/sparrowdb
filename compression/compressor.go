package compression

import (
	"github.com/bkaradzic/go-lz4"
	"github.com/golang/snappy"
)

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

// SNAPPY COMPRESSOR
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

// LZ4 COMPRESSOR
type lz4Compressor struct{}

func (lz4Compressor) Compress(src []byte) []byte {
	b, _ := lz4.Encode(nil, src)
	return b
}

func (lz4Compressor) Decompress(src []byte) ([]byte, error) {
	return lz4.Decode(nil, src)
}

// NewLZ4Compressor returns new lz4Compressor
func NewLZ4Compressor() lz4Compressor {
	return lz4Compressor{}
}
