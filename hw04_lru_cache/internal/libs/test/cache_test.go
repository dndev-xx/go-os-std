package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/dndev-xx/go-os-std/hw04_lru_cache/internal/libs" //nolint:depguard
	"github.com/stretchr/testify/require"                        //nolint:depguard
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := libs.NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := libs.NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val.(*libs.CacheItem).Value)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val.(*libs.CacheItem).Value)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val.(*libs.CacheItem).Value)

		_, ok = c.Get("ccc")
		require.False(t, ok)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := libs.NewCache(10)
		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val.(*libs.CacheItem).Value)
		c.Clear()
		_, ok = c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})
}

func TestCacheMultithreading(_ *testing.T) {
	c := libs.NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(libs.Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(libs.Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()
	wg.Wait()
}
