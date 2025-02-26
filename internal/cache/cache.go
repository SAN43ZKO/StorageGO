package cache

import (
	"sync"
)

type Cache struct {
	items map[string]cacheItem
	mu    sync.Mutex
}

type cacheItem struct {
	value interface{}
}

func NewCache() *Cache {
	return &Cache{
		items: make(map[string]cacheItem),
	}
}	

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = cacheItem{
		value: value,
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	return item.value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}
