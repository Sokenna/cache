package cache

import (
	"sync"
	"time"
)

const (
	//
	NoExpiration      time.Duration = -1
	DefaultExpiration time.Duration = 0
)

type Item struct {
	Data       interface{}
	Expiration int64
}

func (i Item) Expired() bool {
	if i.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > i.Expiration
}

type Cache struct {
	*cache
}

type janitor struct {
	Interval time.Duration
	stop     chan bool
}
type cache struct {
	defaultExpiration time.Duration
	items             map[string]Item
	mu                sync.RWMutex
	onEvicted         func(string, interface{})
	janitor           *janitor
}

func (c *cache) Set(k string, x interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	c.mu.Lock()
	c.items[k] = Item{
		Data:       x,
		Expiration: e,
	}
	c.mu.Unlock()
}
