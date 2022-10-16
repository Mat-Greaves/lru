package lru

import (
	"math"
	"sync"
	"testing"
)

// check cache validates the cache is in the expected state
func checkCache[K comparable, V any](t *testing.T, c *Cache[K, V], es ...entry[K, V]) {
	if l := len(es); c.Len() != l {
		t.Errorf("c.Len() = %d, want %d", c.Len(), l)
		return
	}

	keys := c.Keys()
	for i, v := range es {
		key := keys[i]
		if key != v.key {
			t.Errorf("c.Keys[%d] = %v, want %v", i, key, v.key)
			return
		}

		// TODO: how should we check this? Not all value types are comparable
		// note: this will mutate the cache
		if _, ok := c.Get(key); !ok {
			t.Errorf("c.Get(%v).ok = %v, expected %v", key, ok, true)
		}
	}
}

func getEntry(key string, value int) entry[string, int] {
	return entry[string, int]{
		key:   key,
		value: value,
	}
}

func Test_New(t *testing.T) {
	t.Run("error below miniumum", func(t *testing.T) {
		_, err := New[string, int](0)
		if err == nil {
			t.Error("expected err != nil")
			return
		}
		expected := "cap must be greater than 0"
		if err.Error() != expected {
			t.Errorf("err = %s expected %s", err, expected)
		}
	})

	t.Run("cap 1", func(t *testing.T) {
		cap := 1
		c, err := New[string, int](cap)

		if err != nil {
			t.Errorf("err = %s, expected nil", err)
		}

		if c.cap != cap {
			t.Errorf("c.cap = %d, expected %d", c.cap, cap)
		}
	})

	t.Run("cap max int", func(t *testing.T) {
		cap := math.MaxInt
		c, err := New[string, int](cap)

		if err != nil {
			t.Errorf("err = %s, expected nil", err)
		}

		if c.cap != cap {
			t.Errorf("c.cap = %d, expected %d", c.cap, cap)
		}
	})
}

func Test_Add(t *testing.T) {
	t.Run("add element to empty cache", func(t *testing.T) {
		c, _ := New[string, int](1)
		c.Add("one", 1)
		checkCache(t, c, getEntry("one", 1))
	})

	t.Run("add multiple elements to cache, past cap", func(t *testing.T) {
		c, _ := New[string, int](1)
		c.Add("one", 1)
		c.Add("two", 2)
		checkCache(t, c, getEntry("two", 2))
	})

	t.Run("add multiple elements to cache, under cap", func(t *testing.T) {
		c, _ := New[string, int](2)
		c.Add("one", 1)
		c.Add("two", 2)
		checkCache(t, c, getEntry("two", 2), getEntry("one", 1))
	})
}

func Test_Get(t *testing.T) {
	t.Run("get from empty cache", func(t *testing.T) {
		c, _ := New[string, int](1)
		key := "key"
		v, ok := c.Get(key)
		if ok {
			t.Errorf("c.Get[%v].ok = %v, expected %v", key, ok, false)
		}
		if v != 0 {
			t.Errorf("c.Get[%v] = %d, expected %d", key, v, 0)
		}
	})

	t.Run("get, bring element to front", func(t *testing.T) {
		c, _ := New[string, int](2)
		c.Add("one", 1)
		c.Add("two", 2)
		// bring element at back to front
		c.Get("one")
		checkCache(t, c, getEntry("one", 1), getEntry("two", 2))
	})
}

func Test_Purge(t *testing.T) {
	t.Run("purge empty cache", func(t *testing.T) {
		c, _ := New[string, int](1)
		c.Purge()
		checkCache(t, c)
	})

	t.Run("purge populated cache", func(t *testing.T) {
		c, _ := New[string, int](1)
		c.Add("one", 1)
		c.Purge()
		checkCache(t, c)
	})
}

// run with -race to detect data races on read/write
func Test_Concurrent(t *testing.T) {
	t.Run("concurrent read/write", func(t *testing.T) {
		cap := 100
		c, _ := New[int, int](cap)
		wg := sync.WaitGroup{}
		wg.Add(100)
		for i := 0; i < cap; i++ {
			i := i
			go func() {
				c.Add(i, i)
				wg.Done()
			}()
			go c.Get(i)
		}
		wg.Wait()
		if c.Len() != cap {
			t.Errorf("c.Len() = %d, expected %d", c.Len(), cap)
		}
	})
}
