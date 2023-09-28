package tinycache

import (
	"sync"

	"github.com/TremblingV5/TinyCache/base"
	"github.com/TremblingV5/TinyCache/elimination"
)

type cache struct {
	mu         sync.Mutex
	cache      base.ICache
	cacheBytes int64
}

func newCache(cacheBytes int64) cache {
	c := base.NewBaseCache(cacheBytes, 32, nil)
	c.SetHandle(&elimination.LRU{})
	return cache{
		cache:      c,
		cacheBytes: cacheBytes,
	}
}

func (c *cache) set(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Set(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if v, ok := c.cache.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}
