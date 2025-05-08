package engine

import (
	"log"
	"sync"
	"time"
)

type jokerCache struct {
	cacheMap         map[string]*cacheItem
	checkTime_second int64
	mu               sync.RWMutex
}

type cacheItem struct {
	expiresAt int64
	Value     interface{}
}

func (c *jokerCache) init() {
	c.cacheMap = make(map[string]*cacheItem)
	c.checkTime_second = 30
	log.Println("[info] Cache initialized")
	go c.checkCache()
}

func (c *jokerCache) checkCache() {
	ticker := time.NewTicker(time.Duration(c.checkTime_second) * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now().Unix()
		var keysToDelete []string
		c.mu.RLock()
		for key, value := range c.cacheMap {
			if value.expiresAt < now {
				keysToDelete = append(keysToDelete, key)
			}
		}
		c.mu.RUnlock()
		if len(keysToDelete) > 0 {
			c.mu.Lock()
			for _, key := range keysToDelete {
				delete(c.cacheMap, key)
			}
			c.mu.Unlock()
		}
	}
}

func (c *jokerCache) Set(key string, value interface{}, expiresAt int64) {
	if expiresAt <= 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cacheMap[key] = &cacheItem{
		expiresAt: expiresAt,
		Value:     value,
	}
}

func (c *jokerCache) TryGet(key string) (interface{}, bool) {
	c.mu.RLock()
	value, ok := c.cacheMap[key]
	if ok && value.expiresAt < time.Now().Unix() {
		c.mu.RUnlock()
		c.mu.Lock()
		if value, ok = c.cacheMap[key]; ok && value.expiresAt < time.Now().Unix() {
			delete(c.cacheMap, key)
			c.mu.Unlock()
			return nil, false
		}
		c.mu.Unlock()
		return nil, false
	}
	c.mu.RUnlock()
	if ok {
		return value.Value, true
	} else {
		return nil, false
	}
}

func (c *jokerCache) Remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cacheMap, key)
}

func (c *jokerCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key := range c.cacheMap {
		delete(c.cacheMap, key)
	}
}
