//go:build !solution

package lrucache

import (
	"container/list"
)

type cache struct {
	access *list.List

	mp       map[int]*list.Element
	values   map[int]int
	capacity int
}

func (c *cache) touch(key int) bool {
	if c.capacity < 1 {
		return false
	}

	if node, ok := c.mp[key]; ok {
		c.access.MoveToFront(node)
		return true
	}

	if len(c.mp) >= c.capacity {
		eraseKey := c.access.Back().Value.(int)
		if c.mp[eraseKey] != c.access.Back() {
			panic("something wrong")
		}

		delete(c.values, eraseKey)
		delete(c.mp, eraseKey)
		c.access.Remove(c.access.Back())
	}

	c.mp[key] = c.access.PushFront(key)

	return true
}

func (c *cache) Get(key int) (int, bool) {
	val, ok := c.values[key]
	ok = ok && c.touch(key)
	return val, ok
}

func (c *cache) Set(key, value int) {
	if c.touch(key) {
		c.values[key] = value
	}
}

func (c cache) Range(f func(key int, value int) bool) {
	for e := c.access.Back(); e != nil; e = e.Prev() {
		key := e.Value.(int)
		if !f(key, c.values[key]) {
			break
		}
	}
}

func (c *cache) Clear() {
	c.access = list.New()
	c.values = make(map[int]int, c.capacity)
	c.mp = make(map[int]*list.Element, c.capacity)

}

func New(cap int) Cache {
	res := &cache{capacity: cap}
	res.Clear()
	return res
}
