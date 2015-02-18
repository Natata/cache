package lru

import (
	"sync"
)

type key string
type value interface{}

type node struct {
	k    key
	v    value
	prev *node
	next *node
}

type LRUCache struct {
	lookup     map[key]*node
	head       *node
	currentNum int
	maxNum     int
	sync.RWMutex
}

func NewLRUCache(maxNum int) *LRUCache {
	return &LRUCache{
		lookup:     make(map[key]*node),
		head:       nil,
		currentNum: 0,
		maxNum:     maxNum,
	}
}

func (c *LRUCache) Get(k key) (value, error) {
	c.Lock()
	defer c.Unlock()
	v, err := c.get(k)
	if err == nil {
		c.moveToHead(k)
	}
	return v, err
}

func (c *LRUCache) get(k key) (value, error) {
	n, ok := c.lookup[k]
	if ok {
		return n.v, nil
	}
	return nil, NoValueError
}

func (c *LRUCache) moveToHead(k key) {
	n := c.lookup[k]
	if n == c.head {
		return
	}
	n.prev.next = n.next
	n.next.prev = n.prev

	n.prev, n.next = c.head.prev, c.head
	c.head.prev.next, c.head.prev = n, n

	c.head = n
}

func (c *LRUCache) Set(k key, v value) value {
	c.Lock()
	defer c.Unlock()

	if c.currentNum == c.maxNum {
		c.del()
	}
	return c.set(k, v)
}

func (c *LRUCache) del() {
	switch c.currentNum {
	case 0:
		return
	case 1:
		delete(c.lookup, c.head.prev.k)
		c.head = nil
	default:
		delete(c.lookup, c.head.prev.k)
		c.head.prev = c.head.prev.prev
		c.head.prev.next = c.head
	}
	c.currentNum--
	return
}

func (c *LRUCache) set(k key, newValue value) value {
	//the key is exist
	oldValue, err := c.get(k)
	if err == nil {
		c.lookup[k].v = newValue
		c.moveToHead(k)
		return oldValue
	}

	c.currentNum++
	newNode := &node{
		k: k,
		v: newValue,
	}
	c.lookup[k] = newNode

	// there's no node in the cache
	if c.head == nil {
		newNode.prev = newNode
		newNode.next = newNode
		c.head = newNode
		return nil
	}

	newNode.prev = c.head.prev
	newNode.next = c.head
	c.head.prev.next = newNode
	c.head.prev = newNode
	c.head = newNode
	return nil
}
