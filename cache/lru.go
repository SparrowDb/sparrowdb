package cache

import "sync"

type lru struct {
	used     int64 // Used size of cache in bytes
	capacity int64 // Max size of cache in bytes
	count    int64 // Itens in cache
	mu       sync.RWMutex
	kv       map[uint32]**lruNode
	head     *lruNode
}

type lruNode struct {
	n          *Node
	refs       uint32 // Keeps the reference count
	prev, next *lruNode
}

func (c *lru) insertHead(n *lruNode) {
	n.next = c.head
	n.prev = c.head.prev
	n.prev.next = n
	n.next.prev = n
}

func (c *lru) removeNode(n *lruNode) {
	n.next.prev = n.prev
	n.prev.next = n.next
}

func (c *lru) unref(n *lruNode) {
	if n.refs > 0 {
		n.refs--
	} else if n.refs <= 0 {
		c.decUsed(n.n.size)
		c.removeNode(n)
	}
}

func (c *lru) incUsed(size int32) {
	c.used += int64(size)
	c.count++
}

func (c *lru) decUsed(size int32) {
	c.used -= int64(size)
	c.count--
}

func (c *lru) Usage() (int64, int64, int64) {
	return c.capacity, c.used, c.count
}

func (c *lru) Insert(n *Node) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ln := &lruNode{n: n, refs: 2}
	c.insertHead(ln)

	if _, ok := c.kv[n.key]; !ok {
		c.incUsed(n.size)
	}

	c.kv[n.key] = &ln

	for c.used > c.capacity && c.head.next != c.head {
		old := c.head.next
		c.decUsed(old.n.size)
		c.removeNode(old)
	}
}

func (c *lru) LookUp(key uint32) *Node {
	c.mu.Lock()
	defer c.mu.Unlock()

	var n *Node
	if vaddr, ok := c.kv[key]; ok == true {
		cur := *vaddr
		cur.refs++
		c.removeNode(cur)
		c.insertHead(cur)
		n = cur.n
	}

	return n
}

// NewLRU returns new Cacheable of LRU
func NewLRU(capacity int64) Cacheable {
	c := &lru{
		capacity: capacity,
		kv:       make(map[uint32]**lruNode),
	}

	// Make empty node
	n := &lruNode{}
	n.next = n
	n.prev = n
	c.head = n

	return c
}
