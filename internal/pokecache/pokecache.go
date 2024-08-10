package pokecache

import (
	"time"
	"sync"
)


type CacheEntry struct {
	CreatedAt time.Time
	Val []byte
}

type Cache struct {
	CacheMap map[string]CacheEntry
	mu sync.RWMutex
	interval time.Duration
}

func NewCache(interval time.Duration) *Cache {
	/*
	* @param interval (int): maximum time of persistance of the elements
	* 	in the cache 
	*
	* @return: the new pointer to cache struct 
	*/
	c := new(Cache)
	//c.interval = time.Duration(interval)
	c.interval = interval
	c.CacheMap = make(map[string]CacheEntry, 10)
	go c.ReapLoop()

	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := CacheEntry{CreatedAt: time.Now(), Val : val}
	c.CacheMap[key] = entry
}

func (c *Cache) Get(key string) (val []byte, found bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, found := c.CacheMap[key]
	val = entry.Val

	return val, found
}

func (c *Cache) ReapLoop() {
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
