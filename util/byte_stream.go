package util

import (
	"encoding/binary"
	"unsafe"
)

const (
	uint16Size = 2
	uint32Size = 4
	uint64Size = 8
)

var currentEndianess = binary.LittleEndian

// GetByteOrder returns byte order of running machine
func GetByteOrder() binary.ByteOrder {
	var bo binary.ByteOrder
	var x uint32 = 0x01020304

	switch *(*byte)(unsafe.Pointer(&x)) {
	case 0x01:
		//currentEndianess = binary.BigEndian
	case 0x04:
		//currentEndianess = binary.LittleEndian
	}

	return bo
}

// ByteStream holds []byte definition
type ByteStream struct {
	buf []byte
	cur uint32
}

// Bytes returns the current []byte
func (bs *ByteStream) Bytes() []byte {
	return bs.buf
}

// Size returns ByteStream size
func (bs *ByteStream) Size() int {
	return len(bs.buf)
}

func (bs *ByteStream) appendBytes(buf []byte) {
	bs.buf = append(bs.buf, buf...)
	bs.cur += uint32(len(buf))
}

// PutUInt16 append uint16 to ByteStream
func (bs *ByteStream) PutUInt16(x uint16) {
	b := make([]byte, uint16Size)
	currentEndianess.PutUint16(b, x)
	bs.appendBytes(b)
}

// PutUInt32 append uint32 to ByteStream
func (bs *ByteStream) PutUInt32(x uint32) {
	b := make([]byte, uint32Size)
	currentEndianess.PutUint32(b, x)
	bs.appendBytes(b)
}

// PutUInt64 append uint64 to ByteStream
func (bs *ByteStream) PutUInt64(x uint64) {
	b := make([]byte, uint64Size)
	currentEndianess.PutUint64(b, x)
	bs.appendBytes(b)
}

// PutString append string to ByteStream
func (bs *ByteStream) PutString(x string) {
	bs.PutUInt32(uint32(len(x)))
	bs.appendBytes([]byte(x))
}

// PutBytes append []byte to ByteStream
func (bs *ByteStream) PutBytes(x []byte) {
	bs.PutUInt32(uint32(len(x)))
	bs.appendBytes(x)
}

// GetUInt16 returns uint16 from ByteStream
func (bs *ByteStream) GetUInt16() uint16 {
	x := bs.buf[bs.cur : bs.cur+uint16Size]
	y := currentEndianess.Uint16(x)
	bs.cur += uint16Size
	return y
}

// GetUInt32 returns uint32 from ByteStream
func (bs *ByteStream) GetUInt32() uint32 {
	x := bs.buf[bs.cur : bs.cur+4]
	y := currentEndianess.Uint32(x)
	bs.cur += uint32Size
	return y
}

// GetUInt64 returns uint64 from ByteStream
func (bs *ByteStream) GetUInt64() uint64 {
	x := bs.buf[bs.cur : bs.cur+uint64Size]
	y := currentEndianess.Uint64(x)
	bs.cur += uint64Size
	return y
}

// GetString returns string from ByteStream
func (bs *ByteStream) GetString() string {
	len := bs.GetUInt32()
	x := bs.buf[bs.cur : bs.cur+len]
	bs.cur += len
	return string(x)
}

// GetBytes returns []byte from ByteStream
func (bs *ByteStream) GetBytes() []byte {
	len := bs.GetUInt32()
	x := bs.buf[bs.cur : bs.cur+len]
	bs.cur += len
	return x
}

// Reset put the current cur at the beginning of []byte
func (bs *ByteStream) Reset() {
	bs.cur = 0
}

// NewByteStream returns new ByteStream
func NewByteStream() *ByteStream {
	return &ByteStream{
		buf: []byte{},
		cur: 0,
	}
}

// NewByteStreamFromBytes returns new ByteStream
func NewByteStreamFromBytes(buf []byte) *ByteStream {
	return &ByteStream{
		buf: buf,
		cur: 0,
	}
}
