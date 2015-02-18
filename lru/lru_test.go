package lru

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrivateGet(t *testing.T) {
	lruc := LRUCache{
		lookup:     make(map[key]*node),
		head:       nil,
		currentNum: 0,
		maxNum:     3,
	}
	v, err := lruc.get("test")
	assert.Nil(t, v)
	assert.Equal(t, NoValueError, err)

	n := &node{
		k: "test",
		v: "hi",
	}
	n.prev = n
	n.next = n
	lruc.lookup[n.k] = n
	lruc.head = n
	lruc.currentNum = 1

	v, err = lruc.get("test")
	assert.NoError(t, err)
	assert.Equal(t, "hi", v)

	v, err = lruc.get("tes")
	assert.Nil(t, v)
	assert.Equal(t, NoValueError, err)
}

func TestDel(t *testing.T) {
	lruc := LRUCache{
		lookup:     make(map[key]*node),
		head:       nil,
		currentNum: 0,
		maxNum:     3,
	}
	n1 := &node{
		k: "test1",
		v: 1.618,
	}
	n2 := &node{
		k: "test2",
		v: 2.718,
	}
	n3 := &node{
		k: "test3",
		v: 3.1415,
	}
	lruc.head = n1
	n1.prev = n3
	n1.next = n2
	n2.prev = n1
	n2.next = n3
	n3.prev = n2
	n3.next = n1
	lruc.lookup[n1.k] = n1
	lruc.lookup[n2.k] = n2
	lruc.lookup[n3.k] = n3
	lruc.currentNum = 3

	lruc.del()
	assert.Equal(t, 2, lruc.currentNum)
	assert.Equal(t, n2.v, lruc.head.prev.v)
	assert.Equal(t, 2, len(lruc.lookup))
	fmt.Println(lruc.lookup)
	_, ok := lruc.lookup["test3"]
	assert.False(t, ok)
	assert.Equal(t, lruc.head.prev.prev, lruc.head)

	lruc.del()
	assert.Equal(t, 1, lruc.currentNum)
	assert.Equal(t, n1.v, lruc.head.prev.v)
	assert.Equal(t, 1, len(lruc.lookup))
	_, ok = lruc.lookup["test2"]
	assert.False(t, ok)
	assert.Equal(t, lruc.head.prev.prev, lruc.head)

	lruc.del()
	assert.Equal(t, 0, lruc.currentNum)
	assert.Nil(t, lruc.head)
	assert.Equal(t, 0, len(lruc.lookup))

	lruc.del()
	assert.Equal(t, 0, lruc.currentNum)
	assert.Nil(t, lruc.head)
	assert.Equal(t, 0, len(lruc.lookup))
}

func TestMoveToHead(t *testing.T) {
	lruc := LRUCache{
		lookup:     make(map[key]*node),
		head:       nil,
		currentNum: 0,
		maxNum:     3,
	}
	n1 := &node{
		k: "test1",
		v: 1.618,
	}
	n2 := &node{
		k: "test2",
		v: 2.718,
	}
	n3 := &node{
		k: "test3",
		v: 3.1415,
	}
	lruc.head = n1
	n1.prev = n3
	n1.next = n2
	n2.prev = n1
	n2.next = n3
	n3.prev = n2
	n3.next = n1
	lruc.lookup[n1.k] = n1
	lruc.lookup[n2.k] = n2
	lruc.lookup[n3.k] = n3
	lruc.currentNum = 3

	lruc.moveToHead("test2")
	assert.Equal(t, "test2", lruc.head.k)
	assert.Equal(t, "test1", lruc.head.next.k)
	assert.Equal(t, "test3", lruc.head.next.next.k)
	assert.Equal(t, "test2", lruc.head.next.next.next.k)
	assert.Equal(t, "test3", lruc.head.prev.k)
	assert.Equal(t, "test1", lruc.head.prev.prev.k)
	assert.Equal(t, "test2", lruc.head.prev.prev.prev.k)

	lruc.moveToHead("test3")
	assert.Equal(t, "test3", lruc.head.k)
	assert.Equal(t, "test2", lruc.head.next.k)
	assert.Equal(t, "test1", lruc.head.next.next.k)
	assert.Equal(t, "test3", lruc.head.next.next.next.k)
	assert.Equal(t, "test1", lruc.head.prev.k)
	assert.Equal(t, "test2", lruc.head.prev.prev.k)
	assert.Equal(t, "test3", lruc.head.prev.prev.prev.k)

	lruc.moveToHead("test3")
	assert.Equal(t, "test3", lruc.head.k)
	assert.Equal(t, "test2", lruc.head.next.k)
	assert.Equal(t, "test1", lruc.head.next.next.k)
	assert.Equal(t, "test3", lruc.head.next.next.next.k)
	assert.Equal(t, "test1", lruc.head.prev.k)
	assert.Equal(t, "test2", lruc.head.prev.prev.k)
	assert.Equal(t, "test3", lruc.head.prev.prev.prev.k)
}

func TestInnerSet(t *testing.T) {
	lruc := LRUCache{
		lookup:     make(map[key]*node),
		head:       nil,
		currentNum: 0,
		maxNum:     3,
	}

	lruc.set("test1", 1.618)
	assert.Equal(t, "test1", lruc.head.k)
	assert.Equal(t, 1.618, lruc.head.v)
	assert.Equal(t, 1, len(lruc.lookup))
	assert.Equal(t, 1, lruc.currentNum)

	lruc.set("test1", 2.718)
	assert.Equal(t, "test1", lruc.head.k)
	assert.Equal(t, 2.718, lruc.head.v)
	assert.Equal(t, 1, len(lruc.lookup))
	assert.Equal(t, 1, lruc.currentNum)

	lruc.set("test2", 1.618)
	assert.Equal(t, "test2", lruc.head.k)
	assert.Equal(t, 1.618, lruc.head.v)
	assert.Equal(t, 2, len(lruc.lookup))
	assert.Equal(t, 2, lruc.currentNum)

}

func TestPublicGet(t *testing.T) {
	lruc := LRUCache{
		lookup:     make(map[key]*node),
		head:       nil,
		currentNum: 0,
		maxNum:     3,
	}
	n1 := &node{
		k: "test1",
		v: 1.618,
	}
	n2 := &node{
		k: "test2",
		v: 2.718,
	}
	lruc.head = n1
	n1.prev = n2
	n1.next = n2
	n2.prev = n1
	n2.next = n1
	lruc.lookup[n1.k] = n1
	lruc.lookup[n2.k] = n2
	lruc.currentNum = 2

	v, err := lruc.Get("test2")
	assert.NoError(t, err)
	assert.Equal(t, lruc.head.v, v)

	v, err = lruc.Get("test9")
	assert.Error(t, err)
	assert.Nil(t, v)
}

func TestPublicSet(t *testing.T) {
	lruc := LRUCache{
		lookup:     make(map[key]*node),
		head:       nil,
		currentNum: 0,
		maxNum:     3,
	}
	n1 := &node{
		k: "test1",
		v: 1.618,
	}
	n2 := &node{
		k: "test2",
		v: 2.718,
	}
	n3 := &node{
		k: "test3",
		v: 3.1415,
	}
	lruc.head = n1
	n1.prev = n3
	n1.next = n2
	n2.prev = n1
	n2.next = n3
	n3.prev = n2
	n3.next = n1
	lruc.lookup[n1.k] = n1
	lruc.lookup[n2.k] = n2
	lruc.lookup[n3.k] = n3
	lruc.currentNum = 3

	lruc.Set("test4", 4321)
	assert.Equal(t, 3, len(lruc.lookup))
	assert.Equal(t, 3, lruc.currentNum)
	assert.Equal(t, "test4", lruc.head.k)
	assert.Equal(t, "test2", lruc.head.prev.k)
	assert.Equal(t, "test1", lruc.head.next.k)
	_, ok := lruc.lookup["test3"]
	assert.False(t, ok)

}
