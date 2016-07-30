package util

import "math"

// BloomFilter host data about  probabilistic data structure
type BloomFilter struct {
	size      uint32
	hashCount uint32
	array     []bool
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
		bf.array[v] = true
	}
}

// Contains checks if perhaps BloomFilter contains the key
func (bf *BloomFilter) Contains(key string) bool {
	for _, v := range bf.getHashes(key) {
		if bf.array[v] == false {
			return false
		}
	}
	return true
}

// NewBloomFilter returns new BloomFilter
func NewBloomFilter(elements uint32, falsePositive float64) *BloomFilter {
	bf := &BloomFilter{}
	bf.size = bf.calculateBitSetSize(elements, falsePositive)
	bf.array = make([]bool, bf.size)
	bf.hashCount = bf.calculateHashCount(elements, float64(bf.size))
	return bf
}
