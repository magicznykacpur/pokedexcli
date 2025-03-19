package internal

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	entries  map[string]cacheEntry
	interval time.Duration
	mu       sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{entries: map[string]cacheEntry{}, interval: interval}
	go cache.reapLoop()

	return &cache
}

func (c *Cache) Add(key string, val []byte) {
	newEntry := cacheEntry{createdAt: time.Now(), val: val}

	c.mu.Lock()
	c.entries[key] = newEntry
	c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	entry, ok := c.entries[key]
	c.mu.Unlock()

	return entry.val, ok
}

func (c *Cache) reapLoop() {
	interval := c.interval
	ticker := time.NewTicker(interval)

	for {
		<-ticker.C
		for key, entry := range c.entries {
			if time.Since(entry.createdAt) > interval {
				fmt.Printf("deleting: %s\n", key)
				delete(c.entries, key)
			}
		}

	}
}
