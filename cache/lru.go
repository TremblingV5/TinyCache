package cache

import "github.com/TremblingV5/TinyCache/base"

type LRUCache struct {
	base.Base
}

func add(base *base.Base, node *base.Node, isNew bool) {
	if isNew {
		base.AddBytes(node)
	}
	base.PushFront(node)
}

func remove(base *base.Base, node *base.Node) {
	base.Remove(node)
}

func opAfterView(base *base.Base, node *base.Node) {
	base.Rm(base, node)
	base.Add(base, node, false)
}

func reduce(base *base.Base) {
	for base.IsFull() {
		least := base.End()
		base.Rm(base, least)
		base.Del(string(least.Key()))
		base.ReduceBytes(least)
	}
}

func NewLRUCache(maxBytes int64, shardCount int) *LRUCache {
	base := base.New(maxBytes, shardCount, func(s string, v base.Value) {})

	base.Add = add
	base.OpAfterView = opAfterView
	base.Reduce = reduce
	base.Rm = remove

	return &LRUCache{*base}
}
