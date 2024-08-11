package pokecache

import (
	"time"
	"sync"
)


type CacheEntry[T any] struct {
	CreatedAt time.Time
	Val T
}

type Cache[T any] struct {
	CacheMap map[string]CacheEntry[T]
	mu sync.RWMutex
	interval time.Duration
}

func NewCache[T any](interval time.Duration) *Cache[T] {
	/*
	* @param interval (int): maximum time of persistance of the elements
	* 	in the cache 
	*
	* @return: the new pointer to cache struct 
	*/
	c := new(Cache[T])
	c.interval = interval
	c.CacheMap = make(map[string]CacheEntry[T], 10)
	go c.ReapLoop()

	return c
}

func (c *Cache[T]) Add(key string, val T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := CacheEntry[T]{CreatedAt: time.Now(), Val : val}
	c.CacheMap[key] = entry
}

func (c *Cache[T]) Get(key string) (val T, found bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, found := c.CacheMap[key]
	val = entry.Val

	return val, found
}

func (c *Cache[T]) ReapLoop() {
	ticker := time.NewTicker(c.interval)
	for {
		<-ticker.C
		
		// for the delete operation
    	c.mu.Lock() 

		for key, entry := range c.CacheMap {
			if time.Since(entry.CreatedAt) > c.interval {
				delete(c.CacheMap, key)
			}
		}
		c.mu.Unlock()
	}
}
