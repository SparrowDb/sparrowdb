package util

// Hash32 hashes string into uint32
func Hash32(s string) uint32 {
	b := []byte(s)
	return Murmurhash3_x86_32(b, len(b), 0)
}
