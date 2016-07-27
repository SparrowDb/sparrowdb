package util

import "hash/fnv"

// Hash32 hashes string into uint32
func Hash32(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
