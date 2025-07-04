package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	interval time.Duration
	Data     map[string]cacheEntry
	Mu       sync.Mutex
}

func NewCache(ttl time.Duration) *Cache {
	result := Cache{
		interval: ttl,
		Data:     make(map[string]cacheEntry),
	}

	go result.readLoop()
	return &result
}

func (c *Cache) Add(key string, val []byte) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.Data[key] = cacheEntry{
		val:       val,
		createdAt: time.Now(),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	result, ok := c.Data[key]
	if ok {
		return result.val, ok
	} else {
		return []byte{}, ok
	}
}

func (c *Cache) readLoop() {
	ticker := time.NewTicker(c.interval)
	for {
		time := <-ticker.C
		for key, val := range c.Data {
			if delta := time.Sub(val.createdAt); delta >= c.interval {
				c.Mu.Lock()
				delete(c.Data, key)
				c.Mu.Unlock()
			}
		}
	}
}
