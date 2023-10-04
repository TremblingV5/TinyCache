package linkedlist

import (
	"github.com/TremblingV5/TinyCache/ds/typedef"
	"unsafe"
)

type Node struct {
	last, next *Node
	key        string
	Val        typedef.Value
}

func NewNode(key string, value typedef.Value) *Node {
	return &Node{
		key: key,
		Val: value,
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

func (node *Node) Key() typedef.String {
	return typedef.String(node.key)
}

func (node *Node) Value() typedef.Value {
	return node.Val
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
