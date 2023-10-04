package elimination

import (
	"github.com/TremblingV5/TinyCache/ds/typedef"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/TremblingV5/TinyCache/base"
)

func runLRUCacheTest(t *testing.T, maxBytes int64, shardCount int, test func(*testing.T, *base.BaseCache)) {
	cache := base.NewBaseCache(maxBytes, shardCount, func(s string, v typedef.Value) {})
	cache.SetElimination(&LRU{})
}

func lruSet(t *testing.T, cache *base.BaseCache, key string, value string, expectErr error) {
	err := cache.Set(key, typedef.String(value))
	require.Equal(t, err, expectErr)
}

func lruGet(t *testing.T, cache *base.BaseCache, key string, expectOk bool, expectValue string) {
	value, ok := cache.Get(key)
	require.Equal(t, ok, expectOk)
	if ok {
		require.Equal(t, string(value.(typedef.String)), expectValue)
	}
}

func TestLRUCache(t *testing.T) {
	t.Run("Normal set and get", func(t *testing.T) {
		runLRUCacheTest(t, 64, 32, func(t *testing.T, cache *base.BaseCache) {
			lruSet(t, cache, "1", "1", nil)
			lruGet(t, cache, "1", true, "1")
		})
	})

	t.Run("Delete the lru key", func(t *testing.T) {
		runLRUCacheTest(t, 64, 32, func(t *testing.T, cache *base.BaseCache) {
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
		runLRUCacheTest(t, 64, 32, func(t *testing.T, cache *base.BaseCache) {
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
