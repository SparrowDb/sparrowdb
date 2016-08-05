package util

import (
	"math"

	"github.com/sparrowdb/db/engine"
)

// BloomFilter host data about  probabilistic data structure
type BloomFilter struct {
	size      uint32
	hashCount uint32
	array     []uint8
}

func (bf *BloomFilter) calculateBitSetSize(elements uint32, falsePositive float64) uint32 {
	r := math.Ceil((float64(elements) * math.Log(falsePositive)) / (math.Pow(math.Log(2), 2)))
	return uint32(r * -1)
}

func (bf *BloomFilter) calculateHashCount(elements uint32, bitSetSize float64) uint32 {
	return uint32(((bitSetSize / float64(elements)) * math.Log(2)))
}

func (bf *BloomFilter) getHashes(key string) []uint32 {
	r := make([]uint32, bf.hashCount)

	h1 := Hash32Seed(key, 0)
	h2 := Hash32Seed(key, h1)

	var i uint32
	for i = 0; i < bf.hashCount; i++ {
		rs := (h1 + uint32(i)*h2) % bf.size
		r[i] = uint32(math.Abs(float64(rs)))
	}

	return r
}

// Add adds key to BloomFilter
func (bf *BloomFilter) Add(key string) {
	for _, v := range bf.getHashes(key) {
		bf.array[v] = 1
	}
}

// Contains checks if perhaps BloomFilter contains the key
func (bf *BloomFilter) Contains(key string) bool {
	for _, v := range bf.getHashes(key) {
		if bf.array[v] == 0 {
			return false
		}
	}
	return true
}

// ByteStream returns byte stream of bloom filter data
func (bf *BloomFilter) ByteStream() *engine.ByteStream {
	bs := engine.NewByteStream(engine.LittleEndian)
	bs.PutUInt32(bf.size)
	bs.PutUInt32(bf.hashCount)
	for _, v := range bf.array {
		bs.PutUInt16(uint16(v))
	}
	return bs
}

// NewBloomFilterFromByteStream convert ByteStream to BloomFilter
func NewBloomFilterFromByteStream(bs *engine.ByteStream) *BloomFilter {
	bf := BloomFilter{}
	bf.size = bs.GetUInt32()
	bf.hashCount = bs.GetUInt32()
	var i uint32
	bf.array = make([]uint8, bf.size)
	for i = 0; i < bf.size; i++ {
		bf.array[i] = uint8(bs.GetUInt16())
	}
	return &bf
}

// NewBloomFilter returns new BloomFilter
func NewBloomFilter(elements uint32, falsePositive float32) BloomFilter {
	bf := BloomFilter{}
	bf.size = bf.calculateBitSetSize(elements, float64(falsePositive))
	bf.array = make([]uint8, bf.size)
	bf.hashCount = bf.calculateHashCount(elements, float64(bf.size))
	return bf
}
