package base

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Run("Get default config items", func(t *testing.T) {
		cfg := LoadConfig()
		require.Equal(t, "TinyCache", cfg.Name)
		require.Equal(t, "8001", cfg.Port)
		require.Equal(t, "9999", cfg.ApiPort)
		require.Equal(t, false, cfg.StartApi)
		require.Equal(t, "localhost:8001", cfg.Master)
		require.Equal(t, int64(20480), cfg.MaxBytes)
		require.Equal(t, "LRU", cfg.EliminationMethod)
	})
}
