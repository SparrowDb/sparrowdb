package util

import "unsafe"

const (
	c1 uint32 = 0xcc9e2d51
	c2 uint32 = 0x1b873593
	c3 uint32 = 0x85ebca6b
	c4 uint32 = 0xc2b2ae35
	r1 uint32 = 15
	r2 uint32 = 13
	m  uint32 = 5
	n  uint32 = 0xe6546b64
)

func rot32(x, y uint32) uint32 {
	return (x << y) | (x >> (32 - y))
}

func getBlock32(data []byte, i int) uint32 {
	return *(*uint32)(unsafe.Pointer(&data[i*4]))
}

func finalization32(hash uint32) uint32 {
	hash ^= (hash >> 16)
	hash *= c3
	hash ^= (hash >> 13)
	hash *= c4
	hash ^= (hash >> 16)
	return hash
}

func Murmurhash3_x86_32(data []byte, length int, seed int) uint32 {
	var (
		hash    = uint32(seed)
		nblocks = length / 4
	)

	var k uint32
	for i := 0; i < nblocks; i++ {
		k = getBlock32(data, i)
		k *= c1
		k = rot32(k, r1)
		k *= c2

		hash ^= k
		hash = rot32(hash, r2)
		hash = hash*m + n
	}

	tail := data[nblocks*4:]

	var k1 uint32
	switch length & 3 {
	case 3:
		k1 ^= uint32(tail[2]) << 16
	case 2:
		k1 ^= uint32(tail[1]) << 8
	case 1:
		k1 ^= uint32(tail[0])
		k1 *= c1
		k1 = rot32(k1, r1)
		k1 *= c2
		hash ^= k1
	}

	hash ^= uint32(length)
	hash = finalization32(hash)

	return hash
}
