package base

import (
	"log"
	"sync"

	"github.com/TremblingV5/TinyCache/ds/cmap"
	"github.com/TremblingV5/TinyCache/elimination"
	"github.com/TremblingV5/TinyCache/singleflight"
)

// A Bucket is a cache namespace and associated data loaded spread over
type Bucket struct {
	name      string
	getter    Getter
	mainCache *Cache
	// use singleflight.Group to make sure that
	// each key is only fetched once
	loader *singleflight.Group
}

// A Getter loads data for a key.
type Getter interface {
	Get(key string) ([]byte, error)
}

// A GetterFunc implements Getter with a function.
type GetterFunc func(key string) ([]byte, error)

// Get implements Getter interface function
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

var (
	mu      sync.RWMutex
	buckets = cmap.New[*Bucket](32)
)

// NewBucket create a new instance of Bucket
func NewBucket(name string, cacheBytes int64, getter Getter) *Bucket {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()

	cache := NewCache(cacheBytes)
	cache.SetElimination(&elimination.LRU{})

	g := &Bucket{
		name:      name,
		getter:    getter,
		mainCache: cache,
		loader:    &singleflight.Group{},
	}
	buckets.Set(name, g)
	return g
}

func AddBucketLocally(name string, cacheBytes int64) {
	buckets.Set(name, NewBucket(name, cacheBytes, GetterFunc(func(key string) ([]byte, error) {
		return nil, nil
	})))
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
func (b *Bucket) GetLocally(key string) (ByteView, error) {
	if v, ok := b.mainCache.Get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	return ByteView{}, ErrKeyNotFound
}

func (b *Bucket) SetLocally(key string, value ByteView) {
	b.mainCache.Set(key, value)
}

func (g *Bucket) DelLocally(key string) {
	g.mainCache.Del(key)
}
