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
		c := NewCache(5)

		c.Set("aaa", 500)
		c.Set("bbb", 1000)
		c.Set("ccc", 1500)
		c.Set("ddd", 2000)
		c.Set("eee", 2500)
		c.Set("hhh", 3000) // replacing aaa [hhh, eee, ddd, ccc, bbb]

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb") // bbb not replaced, moving bbb to front, ccc is back now [bbb, hhh, eee, ddd, ccc]
		require.True(t, ok)

		c.Set("iii", 3500) // replacing ccc [iii, bbb, hhh, eee, ddd]

		_, ok = c.Get("ccc")
		require.False(t, ok)

		c.Set("jjj", 4.3)
		c.Set("kkk", []string{"hi", "hello"}) // replacing ccc [kkk, jjj, iii, bbb, hhh]

		_, ok = c.Get("eee")
		require.False(t, ok)

		_, ok = c.Get("ddd")
		require.False(t, ok)

		_, ok = c.Get("hhh")
		require.True(t, ok)
	})

	t.Run("unset logic", func(t *testing.T) {
		c := NewCache(100)

		c.Set("aaa", 500)
		c.Set("aaa", nil)

		val, ok := c.Unset("aaa")
		require.True(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)
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
