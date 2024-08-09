package pokecache

import (
	"time"
	"sync"
)


type cacheEntry struct {
	createdAt time.Time
	val []byte
}

type Cache struct {
	cacheMap map[string]cacheEntry
	mu sync.RWMutex
	interval time.Duration
}

func NewCache(interval int) *Cache {
	/*
	* @param interval (int): maximum time of persistance of the elements
	* 	in the cache expressed in seconds
	*
	* @return: the new cachestruct 
	*/
	c := new(Cache)
	c.interval = time.Duration(interval)
	c.cacheMap = make(map[string]cacheEntry, 10)
	go c.ReapLoop()

	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := cacheEntry{createdAt: time.Now(), val : val}
	c.cacheMap[key] = entry
}

func (c *Cache) Get(key string) (val []byte, found bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, found := c.cacheMap[key]
	val = entry.val

	return val, found
}

func (c *Cache) ReapLoop() {
	ticker := time.NewTicker(c.interval * 1000 * time.Millisecond)
	for i := 0; ;i++ {
		_, ok :=  <- ticker.C
		if  ok {
			for key, entry := range c.cacheMap {
				expirationTime := entry.createdAt.Add(c.interval)
				currentTime := time.Now()

				// check if current time is after eviction time
				if currentTime.After(expirationTime) {
					delete(c.cacheMap, key)
				}
			}	
		}
	}
}