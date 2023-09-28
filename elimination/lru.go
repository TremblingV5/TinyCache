package elimination

import "github.com/TremblingV5/TinyCache/base"

type LRU struct{}

func (l *LRU) Add(list *base.List, node *base.Node) {
	list.PushFront(node)
}

func (l *LRU) OpAfterView(list *base.List, node *base.Node) {
	list.Remove(node)
	list.PushFront(node)
}

func (l *LRU) Remove(list *base.List, node *base.Node) {
	list.Remove(node)
}
