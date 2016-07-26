package cache

import "sync"

type lru struct {
	used     int64
	capacity int64
	count    int64
	lock     sync.RWMutex
	head     *lruNode
	tail     *lruNode
}

type lruNode struct {
	n *Node

	prev, next *lruNode
}

func (c *lru) insertHead(n *lruNode) {
	n.prev = nil
	n.next = c.head
	c.head.prev = n
	c.head = n
}

func (c *lru) moveToFront(n *lruNode) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.count > 1 {
		if n.next == nil {
			c.Ban()
			c.insertHead(n)
		} else {
			n.prev.next = n.next
			n.next.prev = n.prev
			c.insertHead(n)
		}
	}
}

func (c *lru) IncUsed(size int32) {
	c.used += int64(size)
	c.count++
}

func (c *lru) DecUsed(size int32) {
	c.used -= int64(size)
	c.count--
}

func (c *lru) Capactity() int64 {
	return c.capacity
}

func (c *lru) Promote(n *Node) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	new := &lruNode{n: n, prev: nil, next: nil}

	if c.count == 0 {
		c.head = new
		c.tail = c.head
		c.IncUsed(new.n.size)
		return
	}

	if (c.used + int64(n.size)) > c.capacity {
		c.Ban()
		c.DecUsed(n.size)
		c.Promote(n)
	} else {
		c.insertHead(new)
		c.IncUsed(new.n.size)
	}
}

func (c *lru) Ban() {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if c.count > 1 {
		k := c.tail.prev
		c.tail.prev.next = nil
		c.tail = k
	}
}

func (c *lru) Search(key uint32) []byte {
	var cur *lruNode
	cur = c.head

	for {
		if cur == nil {
			break
		}

		if cur.n.key == key {
			v := cur.n.value
			c.moveToFront(cur)
			return v
		}

		cur = cur.next
	}
	return nil
}

func (c *lru) Close() {

}

// NewLRU returns new Cacheable of LRU
func NewLRU(capacity int64) Cacheable {
	c := &lru{
		capacity: capacity,
	}

	return c
}
