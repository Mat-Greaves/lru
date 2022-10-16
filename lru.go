// TODO: can we make this entry[K,V] everywhere cleaner?

// Package lru implements a thread safe lru cache
package lru

import (
	"errors"
	"sync"

	"github.com/Mat-Greaves/lru/list"
)

type cache[K comparable, V any] struct {
	l   *list.List[entry[K, V]]
	m   map[K]*list.Element[entry[K, V]]
	cap int
	mu  sync.Mutex
}

// entry holds a cache entry
type entry[K comparable, V any] struct {
	key   K
	value V
}

// New created a new empty cache with capacity cap
func New[K comparable, V any](cap int) (*cache[K, V], error) {
	if cap <= 0 {
		return nil, errors.New("cap must be greater than 0")
	}
	l := list.New[entry[K, V]]()
	m := make(map[K]*list.Element[entry[K, V]], cap)
	return &cache[K, V]{
		l:   l,
		m:   m,
		cap: cap,
	}, nil
}

// Len returns the number of entries in c
func (c *cache[K, V]) Len() int {
	return c.l.Len()
}

// Keys returns a slice of the keys in cache c, from newest to oldest
func (c *cache[K, V]) Keys() []K {
	keys := make([]K, 0, c.Len())
	for entry := c.l.Front(); entry != nil; entry = entry.Next() {
		keys = append(keys, entry.Value.key)
	}
	return keys
}

// Add element value to with key key to cache c
func (c *cache[K, V]) Add(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.l.Len() == c.cap {
		e := c.l.Remove(c.l.Back())
		delete(c.m, e.key)
	}
	e := c.l.PushFront(entry[K, V]{key: key, value: value})
	c.m[key] = e
}

// Get looks up a key's value from cache c.
// ok indicates whether the value was persent in the cache
func (c *cache[K, V]) Get(key K) (value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e, ok := c.m[key]; ok {
		c.l.MoveToFront(e)
		return e.Value.value, true
	}
	return value, false
}

// Purge clears c
func (c *cache[K, V]) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m = make(map[K]*list.Element[entry[K, V]], c.cap)
	c.l.Init()
}
