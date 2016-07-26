package cache

// Cacheable interface
type Cacheable interface {
	Capactity() int64

	Promote(n *Node)

	Ban()

	Search(key uint32) []byte

	Close()
}

// Cache holds cache operations
type Cache struct {
	cacheable Cacheable
}

// Get gets data from cache
func (c *Cache) Get(key uint32) []byte {
	return c.cacheable.Search(key)
}

// Put puts data in cache
func (c *Cache) Put(key uint32, value []byte) {
	c.cacheable.Promote(&Node{
		key:   key,
		value: value,
		size:  int32(len(value)),
	})
}

// Capactity returns cache capacity
func (c *Cache) Capactity() int64 {
	return c.cacheable.Capactity()
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
