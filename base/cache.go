package base

import (
	"unsafe"

	"github.com/TremblingV5/TinyCache/ds/cmap"
)

type ICache interface {
	Set(string, Value) error
	Get(string) (Value, bool)
	Del(string) error
	Len() int

	IsFull() bool
	Reduce()
	AddBytes(*Node)
	ReduceBytes(*Node)
}

type ICacheElimination interface {
	Add(list *List, node *Node)
	OpAfterView(list *List, node *Node)
	Remove(list *List, node *Node)
}

type BaseCache struct {
	maxBytes int64
	numBytes int64
	list     *List
	dict     cmap.CMap[*Node]

	handle ICacheElimination

	// optional and executed when an entry is purged.
	onEvicted func(key string, value Value)
}

func NewBaseCache(maxBytes int64, shardCount int, onEvicted func(string, Value)) *BaseCache {
	return &BaseCache{
		maxBytes:  maxBytes,
		numBytes:  0,
		list:      NewList(shardCount),
		dict:      cmap.New[*Node](shardCount),
		onEvicted: onEvicted,
	}
}

func (c *BaseCache) SetHandle(handle ICacheElimination) {
	c.handle = handle
}

func (c *BaseCache) Set(key string, value Value) error {
	element, ok := c.dict.Get(key)
	if ok {
		element.value = value
		c.handle.OpAfterView(c.list, element)
		return nil
	}

	newNode := NewNode(key, value)
	c.handle.Add(c.list, newNode)
	c.AddBytes(newNode)
	c.Reduce()
	c.dict.Set(key, newNode)
	return nil
}

func (c *BaseCache) Get(key string) (Value, bool) {
	element, ok := c.dict.Get(key)
	if !ok {
		return nil, ok
	}

	c.handle.OpAfterView(c.list, element)
	return element.value, ok
}

func (c *BaseCache) Del(key string) error {
	element, ok := c.dict.Get(key)
	if !ok {
		return nil
	}

	c.handle.Remove(c.list, element)
	c.dict.Remove(key)
	return nil
}

func (c *BaseCache) Reduce() {
	for c.IsFull() {
		least := c.list.End()
		c.handle.Remove(c.list, least)
		c.dict.Remove(string(least.Key()))
		c.ReduceBytes(least)
	}
}

func (c *BaseCache) AddBytes(element *Node) {
	c.numBytes += int64(element.Value().Len())
	c.numBytes += int64(element.Key().Len())
	c.numBytes += int64(2 * unsafe.Sizeof(element.last))
}

func (c *BaseCache) ReduceBytes(element *Node) {
	c.numBytes -= int64(element.Value().Len())
	c.numBytes -= int64(element.Key().Len())
	c.numBytes -= int64(2 * unsafe.Sizeof(element.last))
}

func (c *BaseCache) IsFull() bool {
	return c.numBytes > c.maxBytes
}

func (c *BaseCache) Len() int {
	return c.dict.Count()
}
