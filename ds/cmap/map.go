package cmap

import (
	"sync"
)

type Shard[V any] struct {
	data map[string]V
	mu   sync.RWMutex
}

type CMap[V any] struct {
	shards      []*Shard[V]
	sharding    func(key string) uint32
	shard_count int
}

func create[V any](shard_count int, sharding func(key string) uint32) CMap[V] {
	cmap := CMap[V]{
		shards:      make([]*Shard[V], shard_count),
		sharding:    sharding,
		shard_count: shard_count,
	}

	for i := 0; i < shard_count; i++ {
		cmap.shards[i] = &Shard[V]{
			data: make(map[string]V),
		}
	}

	return cmap
}

func New[V any](shard_count int) CMap[V] {
	return create[V](shard_count, func(key string) uint32 {
		hash := uint32(2166136261)
		const prime32 = uint32(16777619)
		keyLength := len(key)
		for i := 0; i < keyLength; i++ {
			hash *= prime32
			hash ^= uint32(key[i])
		}
		return hash
	})
}

func (cmap CMap[V]) GetShard(key string) *Shard[V] {
	return cmap.shards[cmap.sharding(key)%uint32(cmap.shard_count)]
}

func (cmap CMap[V]) Set(key string, value V) {
	shard := cmap.GetShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	shard.data[key] = value
}

func (cmap CMap[V]) Get(key string) (V, bool) {
	shard := cmap.GetShard(key)
	shard.mu.RLock()
	defer shard.mu.RUnlock()
	val, ok := shard.data[key]
	return val, ok
}

func (cmap CMap[V]) Count() int {
	count := 0
	for i := 0; i < cmap.shard_count; i++ {
		shard := cmap.shards[i]
		shard.mu.RLock()
		count += len(shard.data)
		shard.mu.RUnlock()
	}
	return count
}

func (cmap CMap[V]) Remove(key string) {
	shard := cmap.GetShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	delete(shard.data, key)
}

func (cmap CMap[V]) Clear() {
	for i := 0; i < cmap.shard_count; i++ {
		shard := cmap.shards[i]
		shard.mu.RLock()
		shard = &Shard[V]{}
		shard.mu.RUnlock()
	}
}
