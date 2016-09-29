package util

import "github.com/SparrowDb/sparrowdb/util/murmurhash3"

// Hash32 hashes string with default seed = 0 into uint32
func Hash32(s string) uint32 {
	return Hash32Seed(s, 0)
}

// Hash32Seed hashes string with seed into uint32
func Hash32Seed(s string, seed uint32) uint32 {
	b := []byte(s)
	return murmurhash3.Murmurhash3X86_32(b, seed)
}
