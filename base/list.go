package base

type List struct {
	head *Node
	tail *Node
}

func NewList(shardCount int) *List {
	head := NewNode("", nil)
	tail := NewNode("", nil)

	head.SetNext(tail)
	tail.SetLast(head)

	return &List{
		head: head,
		tail: tail,
	}
}

func (l *List) Begin() *Node {
	return l.head.Next()
}

func (l *List) End() *Node {
	return l.tail
}

func (l *List) PushFront(node *Node) {
	node.SetLast(l.head)
	node.SetNext(l.head.Next())
	l.head.Next().SetLast(node)
	l.head.SetNext(node)
}

func (l *List) Remove(node *Node) {
	node.Last().SetNext(node.Next())
	node.Next().SetLast(node.Last())
}
