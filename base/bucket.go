package base

import (
	"log"
	"sync"

	"github.com/TremblingV5/TinyCache/ds/cmap"
	"github.com/TremblingV5/TinyCache/ds/typedef"
	"github.com/TremblingV5/TinyCache/elimination"
	"github.com/TremblingV5/TinyCache/singleflight"
)

// A Bucket is a cache namespace and associated data loaded spread over
type Bucket struct {
	name      string
	mainCache *Cache
	// use singleflight.Group to make sure that
	// each key is only fetched once
	loader *singleflight.Group
}

var (
	mu      sync.RWMutex
	buckets = cmap.New[*Bucket](32)
)

// NewBucket create a new instance of Bucket
func NewBucket(name string, cacheBytes int64) *Bucket {
	mu.Lock()
	defer mu.Unlock()

	cache := NewCache(cacheBytes)
	cache.SetElimination(&elimination.LRU{})

	g := &Bucket{
		name:      name,
		mainCache: cache,
		loader:    &singleflight.Group{},
	}
	buckets.Set(name, g)
	return g
}

func AddBucketLocally(name string, cacheBytes int64) {
	buckets.Set(name, NewBucket(name, cacheBytes))
}

func RemoveBucketLocally(name string) {
	buckets.Remove(name)
}

// GetBucket returns the named group previously created with NewBucket, or
// nil if there's no such group.
func GetBucket(name string) *Bucket {
	mu.RLock()
	defer mu.RUnlock()
	if g, ok := buckets.Get(name); ok {
		return g
	}
	return nil
}
func (b *Bucket) GetLocally(key string) (typedef.DataBytes, error) {
	v, err := b.loader.Do(key, func() (interface{}, error) {
		if v, ok := b.mainCache.Get(key); ok {
			log.Println("[GeeCache] hit")
			return v, nil
		}

		return typedef.DataBytes{}, ErrKeyNotFound
	})
	return v.(typedef.DataBytes), err
}

func (b *Bucket) SetLocally(key string, value typedef.DataBytes) {
	b.mainCache.Set(key, value)
}

func (g *Bucket) DelLocally(key string) {
	g.mainCache.Del(key)
}
