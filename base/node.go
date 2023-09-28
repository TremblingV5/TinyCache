package base

import "unsafe"

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

func (node *Node) Size() int {
	return int(unsafe.Sizeof(*node))
}

func (node *Node) NumBytes() int64 {
	numBytes := 0
	numBytes += node.Value().Len()
	numBytes += node.Key().Len()
	numBytes += int(unsafe.Sizeof(node.last) * 2)
	numBytes += node.Size()
	return int64(numBytes)
}
