package hw04_lru_cache //nolint:golint,stylecheck

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)
		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)

		c.Clear()

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)

		_, ok = c.Get("ccc")
		require.False(t, ok)

		wasInCache := c.Set("ccc", 100)
		require.False(t, wasInCache)

		val, ok := c.Get("ccc")
		require.True(t, ok)
		require.Equal(t, 100, val)
	})

	t.Run("ejection", func(t *testing.T) {
		c := NewCache(3)
		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)
		c.Set("ddd", 400)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.True(t, ok)

		_, ok = c.Get("ccc")
		require.True(t, ok)

		_, ok = c.Get("ddd")
		require.True(t, ok)
	})

	t.Run("lru", func(t *testing.T) {
		c := NewCache(3)
		c.Set("0", 100)
		c.Set("1", 200)
		c.Set("2", 300)

		for i := 0; i < 100; i++ {
			c.Set(Key(strconv.Itoa(rand.Intn(3))), rand.Intn(42))
			c.Get(Key(strconv.Itoa(rand.Intn(3))))
		}
		c.Set("strike", 777)

		_, ok0 := c.Get("0")
		_, ok1 := c.Get("1")
		_, ok2 := c.Get("2")
		require.ElementsMatch(t, []bool{true, true, false}, []bool{ok0, ok1, ok2})
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
