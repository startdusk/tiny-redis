package tinyredis

import (
	"sync"
	"time"
)

const (
	DefaultExpiration time.Duration = 0
	NoExpiration      time.Duration = -1
)

type Item struct {
	object     interface{}
	expiration int64
	hasExpired bool
}

func (i Item) Expired() bool {
	if i.hasExpired {
		return time.Now().UnixNano() > i.expiration
	}

	return false
}

func (i Item) Value() interface{} {
	return i.object
}

type Cache struct {
	defaultExpiration time.Duration
	items             map[string]Item
	mu                sync.RWMutex
	gcInterval        time.Duration
	stopGC            chan bool
}

func (c *Cache) Set(k string, v interface{}, d time.Duration) {
	c.items[k] = Item{
		object:     v,
		hasExpired: d > 0,
		expiration: time.Now().Add(d).UnixNano(),
	}
}

func (c *Cache) Get(k string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[k]
	if !found {
		return nil, false
	}

	if item.Expired() {
		return nil, false
	}

	return item.object, true
}

func (c *Cache) Delete(k string) {
	c.mu.Lock()
	delete(c.items, k)
	c.mu.Unlock()
}

func (c *Cache) All() map[string]Item {
	c.mu.RLock()
	defer c.mu.RUnlock()

	all := make(map[string]Item, len(c.items))
	for k, v := range c.items {
		all[k] = v
	}
	return all
}

func NewCache(defaultExpiration, gcInterval time.Duration) *Cache {
	return &Cache{
		defaultExpiration: defaultExpiration,
		gcInterval:        gcInterval,
		items:             make(map[string]Item),
		stopGC:            make(chan bool),
	}
}
