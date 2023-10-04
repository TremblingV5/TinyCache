package base

import (
	"github.com/TremblingV5/TinyCache/ds/linkedlist"
	"github.com/TremblingV5/TinyCache/ds/typedef"
	"unsafe"

	"github.com/TremblingV5/TinyCache/ds/cmap"
)

type ICache interface {
	Set(string, typedef.Value) error
	Get(string) (typedef.Value, bool)
	Del(string) error
	Len() int

	IsFull() bool
	Reduce()
	AddBytes(*linkedlist.Node)
	ReduceBytes(*linkedlist.Node)

	SetElimination(ICacheElimination)
}

type ICacheElimination interface {
	Add(list *linkedlist.List, node *linkedlist.Node)
	OpAfterView(list *linkedlist.List, node *linkedlist.Node)
	Remove(list *linkedlist.List, node *linkedlist.Node)
}

type BaseCache struct {
	maxBytes int64
	numBytes int64
	list     *linkedlist.List
	dict     cmap.CMap[*linkedlist.Node]

	handle ICacheElimination

	// optional and executed when an entry is purged.
	onEvicted func(key string, value typedef.Value)
}

func NewBaseCache(maxBytes int64, shardCount int, onEvicted func(string, typedef.Value)) *BaseCache {
	return &BaseCache{
		maxBytes:  maxBytes,
		numBytes:  0,
		list:      linkedlist.NewList(shardCount),
		dict:      cmap.New[*linkedlist.Node](shardCount),
		onEvicted: onEvicted,
	}
}

func (c *BaseCache) SetElimination(handle ICacheElimination) {
	c.handle = handle
}

func (c *BaseCache) Set(key string, value typedef.Value) error {
	element, ok := c.dict.Get(key)
	if ok {
		element.Val = value
		c.handle.OpAfterView(c.list, element)
		return nil
	}

	newNode := linkedlist.NewNode(key, value)
	c.handle.Add(c.list, newNode)
	c.AddBytes(newNode)
	c.Reduce()
	c.dict.Set(key, newNode)
	return nil
}

func (c *BaseCache) Get(key string) (typedef.Value, bool) {
	element, ok := c.dict.Get(key)
	if !ok {
		return nil, ok
	}

	c.handle.OpAfterView(c.list, element)
	return element.Val, ok
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

func (c *BaseCache) AddBytes(element *linkedlist.Node) {
	c.numBytes += int64(element.Value().Len())
	c.numBytes += int64(element.Key().Len())
	c.numBytes += int64(2 * unsafe.Sizeof(element.Last()))
}

func (c *BaseCache) ReduceBytes(element *linkedlist.Node) {
	c.numBytes -= int64(element.Value().Len())
	c.numBytes -= int64(element.Key().Len())
	c.numBytes -= int64(2 * unsafe.Sizeof(element.Last()))
}

func (c *BaseCache) IsFull() bool {
	return c.numBytes > c.maxBytes
}

func (c *BaseCache) Len() int {
	return c.dict.Count()
}
