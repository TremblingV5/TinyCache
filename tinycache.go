package tinycache

import (
	"fmt"
	"github.com/TremblingV5/TinyCache/pb"
	"log"
	"sync"

	"github.com/TremblingV5/TinyCache/base"
	"github.com/TremblingV5/TinyCache/ds/cmap"
	"github.com/TremblingV5/TinyCache/singleflight"
)

// A Bucket is a cache namespace and associated data loaded spread over
type Bucket struct {
	name      string
	getter    Getter
	mainCache cache
	peers     PeerPicker
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

	secondaryServerLock sync.RWMutex
	secondaryServerList = cmap.New[string](32)
)

// NewBucket create a new instance of Bucket
func NewBucket(name string, cacheBytes int64, getter Getter) *Bucket {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Bucket{
		name:      name,
		getter:    getter,
		mainCache: newCache(cacheBytes),
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

// Get value for a key from cache
func (g *Bucket) Get(key string) (base.ByteView, error) {
	if key == "" {
		return base.ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	return g.load(key)
}

func (b *Bucket) GetLocally(key string) (base.ByteView, error) {
	if v, ok := b.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	return base.ByteView{}, ErrKeyNotFound
}

func (g *Bucket) Set(key string, value base.ByteView) {
	g.mainCache.set(key, value)
}

func (b *Bucket) SetLocally(key string, value base.ByteView) {
	b.mainCache.set(key, value)
}

func (g *Bucket) Del(key string) {
	g.mainCache.del(key)
}

func (g *Bucket) DelLocally(key string) {
	g.mainCache.del(key)
}

// RegisterPeers registers a PeerPicker for choosing remote peer
func (g *Bucket) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

func (g *Bucket) load(key string) (value base.ByteView, err error) {
	// each key is only fetched once (either locally or remotely)
	// regardless of the number of concurrent callers.
	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err = g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}

		return g.getLocally(key)
	})

	if err == nil {
		return viewi.(base.ByteView), nil
	}
	return
}

func (b *Bucket) upload(key string, value base.ByteView) {

}

func (g *Bucket) populateCache(key string, value base.ByteView) {
	g.mainCache.set(key, value)
}

func (g *Bucket) getLocally(key string) (base.ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return base.ByteView{}, err

	}
	value := base.ByteView{B: base.CloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Bucket) getFromPeer(peer PeerHandle, key string) (base.ByteView, error) {
	req := &pb.GetKeyRequest{
		Bucket: g.name,
		Key:    key,
	}
	res := &pb.GetKeyResponse{}
	err := peer.Get(req, res)
	if err != nil {
		return base.ByteView{}, err
	}
	return base.ByteView{B: res.Value}, nil
}

func AddSecondaryServer(name string, addr string) {
	secondaryServerLock.Lock()
	defer secondaryServerLock.Unlock()

	secondaryServerList.Set(name, addr)
}
