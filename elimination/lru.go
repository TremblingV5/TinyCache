package elimination

import (
	"github.com/TremblingV5/TinyCache/ds/linkedlist"
)

type LRU struct{}

func (l *LRU) Add(list *linkedlist.List, node *linkedlist.Node) {
	list.PushFront(node)
}

func (l *LRU) OpAfterView(list *linkedlist.List, node *linkedlist.Node) {
	list.Remove(node)
	list.PushFront(node)
}

func (l *LRU) Remove(list *linkedlist.List, node *linkedlist.Node) {
	list.Remove(node)
}
