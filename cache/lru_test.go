package cache

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/TremblingV5/TinyCache/base"
)

func runLRUCacheTest(t *testing.T, maxBytes int64, shardCount int, test func(*testing.T, *LRUCache)) {
	lruCache := NewLRUCache(maxBytes, shardCount)
	test(t, lruCache)
}

func lruSet(t *testing.T, cache *LRUCache, key string, value string, expectErr error) {
	err := cache.Set(key, base.String(value))
	require.Equal(t, err, expectErr)
}

func lruGet(t *testing.T, cache *LRUCache, key string, expectOk bool, expectValue string) {
	value, ok := cache.Get(key)
	require.Equal(t, ok, expectOk)
	if ok {
		require.Equal(t, string(value.(base.String)), expectValue)
	}
}

func lruIsFull(t *testing.T, cache *LRUCache, expectValue bool) {
	require.Equal(t, cache.IsFull(), expectValue)
}

func TestLRUCache(t *testing.T) {
	t.Run("Normal set and get", func(t *testing.T) {
		runLRUCacheTest(t, 64, 32, func(t *testing.T, cache *LRUCache) {
			lruSet(t, cache, "1", "1", nil)
			lruGet(t, cache, "1", true, "1")
		})
	})

	t.Run("Delete the lru key", func(t *testing.T) {
		runLRUCacheTest(t, 64, 32, func(t *testing.T, cache *LRUCache) {
			lruSet(t, cache, "1", "1", nil)
			lruSet(t, cache, "2", "2", nil)
			lruSet(t, cache, "3", "3", nil)
			lruSet(t, cache, "4", "4", nil)

			lruGet(t, cache, "1", false, "")
			lruGet(t, cache, "2", true, "2")
			lruGet(t, cache, "3", true, "3")
			lruGet(t, cache, "4", true, "4")
		})
	})

	t.Run("Delete the lru key with double view for the first key", func(t *testing.T) {
		runLRUCacheTest(t, 64, 32, func(t *testing.T, cache *LRUCache) {
			lruSet(t, cache, "1", "1", nil)
			lruSet(t, cache, "2", "2", nil)
			lruSet(t, cache, "3", "3", nil)
			lruGet(t, cache, "1", true, "1")
			lruSet(t, cache, "4", "4", nil)

			lruGet(t, cache, "1", true, "1")
			lruGet(t, cache, "2", false, "")
			lruGet(t, cache, "3", true, "3")
			lruGet(t, cache, "4", true, "4")
		})
	})
}
