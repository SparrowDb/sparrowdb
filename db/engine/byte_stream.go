package engine

import (
	"encoding/binary"
	"unsafe"
)

const (
	uint16Size = 2
	uint32Size = 4
	uint64Size = 8
)

// LittleEndian byte order
var LittleEndian = binary.LittleEndian

// BigEndian  byte order
var BigEndian = binary.BigEndian

// GetByteOrder returns byte order of running machine
func GetByteOrder() binary.ByteOrder {
	var bo binary.ByteOrder
	var x uint32 = 0x01020304

	switch *(*byte)(unsafe.Pointer(&x)) {
	case 0x01:
		bo = BigEndian
	case 0x04:
		bo = LittleEndian
	}

	return bo
}

// ByteStream holds []byte definition
type ByteStream struct {
	buf   []byte
	cur   uint32
	order binary.ByteOrder
}

// Bytes returns the current []byte
func (bs *ByteStream) Bytes() []byte {
	return bs.buf
}

func (bs *ByteStream) appendBytes(buf []byte) {
	bs.buf = append(bs.buf, buf...)
	bs.cur += uint32(len(buf))
}

// PutUInt16 append uint16 to ByteStream
func (bs *ByteStream) PutUInt16(x uint16) {
	b := make([]byte, uint16Size)
	bs.order.PutUint16(b, x)
	bs.appendBytes(b)
}

// PutUInt32 append uint32 to ByteStream
func (bs *ByteStream) PutUInt32(x uint32) {
	b := make([]byte, uint32Size)
	bs.order.PutUint32(b, x)
	bs.appendBytes(b)
}

// PutUInt64 append uint64 to ByteStream
func (bs *ByteStream) PutUInt64(x uint64) {
	b := make([]byte, uint64Size)
	bs.order.PutUint64(b, x)
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
	y := bs.order.Uint16(x)
	bs.cur += uint16Size
	return y
}

// GetUInt32 returns uint32 from ByteStream
func (bs *ByteStream) GetUInt32() uint32 {
	x := bs.buf[bs.cur : bs.cur+4]
	y := bs.order.Uint32(x)
	bs.cur += uint32Size
	return y
}

// GetUInt64 returns uint64 from ByteStream
func (bs *ByteStream) GetUInt64() uint64 {
	x := bs.buf[bs.cur : bs.cur+uint64Size]
	y := bs.order.Uint64(x)
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
func NewByteStream(byteOrder binary.ByteOrder) *ByteStream {
	return &ByteStream{
		buf:   []byte{},
		cur:   0,
		order: byteOrder,
	}
}

// NewByteStream returns new ByteStream
func NewByteStreamFromBytes(buf []byte, byteOrder binary.ByteOrder) *ByteStream {
	return &ByteStream{
		buf:   buf,
		cur:   0,
		order: byteOrder,
	}
}
