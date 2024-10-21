package cache

import (
  "sync"
  "time"
)

type CacheItem struct {
  Value    interface{}
  Expiration    int64
}

type Cache struct {
  items    map[string]CacheItem 
  mu    sync.RWMutex
}

func NewCache() *Cache {
  return &Cache{
    items: make(map[string]CacheItem), 
  }
}

func (c *Cache) Set(key string, value interface{}, duration time.duration) {
  c.mu.lock()
  defer c.mu.Unlock()

  expiration := time.Now().Add(duration).UnixNano()
  c.items[key] = CacheItem{
    Value: value, 
    Expiration: expiration, 
  }
}

func (c *Cache) Get(key string) (interface{}, bool) {
  c.mu.Rlock()
  defer c.mu.RUnlock()

  item, found := c.items[key]
  if !found {
    return nil, false
  }

  if time.Now().UnixNano() > item.Expiration {
    return nil, false 
  }

  return item.Value, true
}

func (c *Cache) Delete(key string) {
  c.mu.Lock()
  defer c.mu.Unlock()

  delete(c.items, key)
}
