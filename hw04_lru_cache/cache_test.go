package hw04lrucache

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
		c := NewCache(2)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)
		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)
		_, ok := c.Get("aaa")
		require.True(t, ok)
		_, ok = c.Get("bbb")
		require.True(t, ok)

		c.Clear()

		_, ok = c.Get("aaa")
		require.False(t, ok)
		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("overflow simple", func(t *testing.T) {
		c := NewCache(2)

		// fill whole capacity
		c.Set("aaa", 100)
		c.Set("bbb", 200)
		// add item for overflow
		c.Set("ccc", 300)

		// check 1st added item war discarded
		_, ok := c.Get("aaa")
		require.False(t, ok)
		// check 2nd still there
		_, ok = c.Get("bbb")
		require.True(t, ok)
	})

	t.Run("overflow complex", func(t *testing.T) {
		c := NewCache(3)

		// fill whole capacity
		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)

		// get 1st item to push it to front
		c.Get("aaa")

		// set 3rd item to push it to front
		c.Set("ccc", 301)

		// add item for overflow
		c.Set("ddd", 400)

		// check 2nd added (on start) element is gone
		_, ok := c.Get("bbb")
		require.False(t, ok)
	})
}

func TestCacheMultithreading(_ *testing.T) {
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
