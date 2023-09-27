package base

import (
	"unsafe"

	"github.com/TremblingV5/TinyCache/ds/cmap"
)

type Value interface {
	Len() int
}

type Node struct {
	last, next *Node
	key        string
	value      Value
}

func NewNode(key string, value Value) *Node {
	return &Node{
		key:   key,
		value: value,
	}
}

func (node *Node) Next() *Node {
	return node.next
}

func (node *Node) SetNext(nextNode *Node) *Node {
	node.next = nextNode
	return nextNode
}

func (node *Node) Last() *Node {
	return node.last
}

func (node *Node) SetLast(lastNode *Node) *Node {
	node.last = lastNode
	return lastNode
}

func (node *Node) Key() String {
	return String(node.key)
}

func (node *Node) Value() Value {
	return node.value
}

func (node *Node) Size() int64 {
	return int64(unsafe.Sizeof(*node))
}

type Cache interface {
	Set(string, Value) error
	Get(string) (Value, bool)
	Del(string) error
	Len() int
}

type Base struct {
	maxBytes int64
	numBytes int64
	head     *Node
	tail     *Node
	cache    cmap.CMap[*Node]

	Add         func(base *Base, node *Node, isNew bool)
	Rm          func(base *Base, node *Node)
	OpAfterView func(base *Base, node *Node)
	Reduce      func(base *Base)

	// optional and executed when an entry is purged.
	onEvicted func(key string, value Value)
}

func newElement(key string, data Value) *Node {
	return &Node{
		key:   key,
		value: data,
		last:  nil,
		next:  nil,
	}
}

func (element *Node) Len() int {
	return element.value.Len()
}

func New(maxBytes int64, shardCount int, onEvicted func(string, Value)) *Base {
	head := NewNode("", nil)
	tail := NewNode("", nil)

	head.SetNext(tail)
	tail.SetLast(head)

	return &Base{
		maxBytes:  maxBytes,
		numBytes:  0,
		head:      head,
		tail:      tail,
		cache:     cmap.New[*Node](shardCount),
		onEvicted: onEvicted,
	}
}

func (c *Base) Begin() *Node {
	return c.head.Next()
}

func (c *Base) End() *Node {
	return c.tail.Last()
}

func (c *Base) Set(key string, value Value) error {
	element, ok := c.cache.Get(key)
	if ok {
		element.value = value
		c.OpAfterView(c, element)
		return nil
	}

	newElement := newElement(key, value)
	c.Add(c, newElement, true)
	c.Reduce(c)
	c.cache.Set(key, newElement)
	return nil
}

func (c *Base) Get(key string) (Value, bool) {
	element, ok := c.cache.Get(key)
	if !ok {
		return nil, ok
	}

	c.OpAfterView(c, element)
	return element.value, ok
}

func (c *Base) Del(key string) error {
	element, ok := c.cache.Get(key)
	if !ok {
		return nil
	}

	c.Rm(c, element)
	c.cache.Remove(key)
	return nil
}

func (c *Base) PushFront(element *Node) {
	element.SetLast(c.head)
	element.SetNext(c.head.Next())
	c.head.Next().SetLast(element)
	c.head.SetNext(element)
}

func (c *Base) Remove(element *Node) {
	element.Last().SetNext(element.Next())
	element.Next().SetLast(element.Last())
}

func (c *Base) AddBytes(element *Node) {
	c.numBytes += int64(element.Value().Len())
	c.numBytes += int64(element.Key().Len())
	c.numBytes += int64(2 * unsafe.Sizeof(element.last))
}

func (c *Base) ReduceBytes(element *Node) {
	c.numBytes -= int64(element.Value().Len())
	c.numBytes -= int64(element.Key().Len())
	c.numBytes -= int64(2 * unsafe.Sizeof(element.last))
}

func (c *Base) IsFull() bool {
	return c.numBytes > c.maxBytes
}

func (c *Base) Len() int {
	return c.cache.Count()
}
