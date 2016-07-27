package cache

// Cacheable interface
type Cacheable interface {
	// Returns cache capacity, cache used in bytes
	// and total itens in cache
	Usage() (int64, int64, int64)

	Insert(n *Node)

	LookUp(key uint32) *Node
}

// Cache holds cache operations
type Cache struct {
	cacheable Cacheable
}

// Get gets data from cache
func (c *Cache) Get(key uint32) []byte {
	if v := c.cacheable.LookUp(key); v != nil {
		return v.value
	}
	return nil
}

// Put puts data in cache
func (c *Cache) Put(key uint32, value []byte) {
	c.cacheable.Insert(&Node{
		key:   key,
		value: value,
		size:  int32(len(value)),
	})
}

// Usage rseturns cache capacity, cache used in bytes
// and total itens in cache
func (c *Cache) Usage() (int64, int64, int64) {
	return c.cacheable.Usage()
}

// NewCache returns new Cache
func NewCache(c Cacheable) *Cache {
	cache := Cache{
		cacheable: c,
	}
	return &cache
}

// Node holds cache entry
type Node struct {
	key   uint32
	value []byte
	size  int32
}
