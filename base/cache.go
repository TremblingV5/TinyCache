package base

import (
	"sync"
)

type Cache struct {
	mu         sync.Mutex
	cache      ICache
	cacheBytes int64
}

func NewCache(cacheBytes int64) *Cache {
	c := NewBaseCache(cacheBytes, 32, nil)
	return &Cache{
		mu:         sync.Mutex{},
		cache:      c,
		cacheBytes: cacheBytes,
	}
}

func (c *Cache) SetElimination(e ICacheElimination) {
	c.cache.SetElimination(e)
}

func (c *Cache) Set(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Set(key, value)
}

func (c *Cache) Get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if v, ok := c.cache.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}

func (c *Cache) Del(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache.Del(key)
}
