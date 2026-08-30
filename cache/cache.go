package cache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	entries  map[string]Entry
	mutex    *sync.RWMutex
	duration time.Duration
}

type Entry struct {
	data      []byte
	createdAt time.Time
}

func NewCache(interval time.Duration) *Cache {
	c := Cache{
		entries:  make(map[string]Entry),
		mutex:    new(sync.RWMutex),
		duration: interval,
	}

	go c.ReapLoop()

	return &c
}

func (c *Cache) ReapLoop() {
	ticker := time.NewTicker(c.duration)

	for range ticker.C {
		for k, v := range c.entries {
			c.mutex.Lock()
			if time.Now().Sub(v.createdAt).Seconds() > c.duration.Seconds() {
				delete(c.entries, k)
			}
			c.mutex.Unlock()
		}
	}
}

func (c *Cache) Add(key string, data []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries[key] = Entry{
		data:      data,
		createdAt: time.Now(),
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	e, ok := c.entries[key]

	if !ok {
		return nil, ok
	}

	return e.data, ok
}

func (c *Cache) PrintEntries() {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for k, _ := range c.entries {
		fmt.Println(k)
	}
}
