package gocache

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"
)

// An AtomicInt is an int64 to be accessed atomically.
type AtomicInt int64

// MemCache is an LRU cache. It is safe for concurrent access.
type Cache struct {
	mutex             sync.RWMutex
	maxItemSize       int
	cacheList         *list.List
	cache             map[interface{}]*list.Element
	hits, gets        AtomicInt
	defaultExpiration time.Duration
}

//NewMemCache If maxItemSize is zero, the cache has no limit.
//if maxItemSize is not zero, when cache's size beyond maxItemSize,start to swap
func NewCache(maxItemSize int, defaultExpiration time.Duration) *Cache {
	cache := &Cache{
		maxItemSize:       maxItemSize,
		cacheList:         list.New(),
		cache:             make(map[interface{}]*list.Element),
		defaultExpiration: defaultExpiration,
	}
	go cache.run()
	return cache
}

//Status return the status of cache
func (c *Cache) Status() *CacheStatus {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return &CacheStatus{
		MaxItemSize: c.maxItemSize,
		CurrentSize: c.cacheList.Len(),
		Gets:        c.gets.Get(),
		Hits:        c.hits.Get(),
	}
}

//Get value with key
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	c.gets.Add(1)
	if ele, hit := c.cache[key]; hit {
		if ele.Value.(*Item).Expiration > 0 {
			//已经过期
			if time.Now().UnixNano() > ele.Value.(*Item).Expiration {
				return nil, false
			}
		}
		c.hits.Add(1)
		c.cacheList.MoveToFront(ele)
		return ele.Value.(*Item).value, true
	}
	return nil, false
}

//Set a value with key
func (c *Cache) Set(key string, value interface{}, d time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.cache == nil {
		c.cache = make(map[interface{}]*list.Element)
		c.cacheList = list.New()
	}
	var e int64
	if d == DefaultExpiration {
		d = c.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	if ele, ok := c.cache[key]; ok {
		ele.Value = &Item{
			key:        key,
			value:      value,
			Expiration: e,
		}
		c.cacheList.MoveToFront(ele)
		return
	}

	ele := c.cacheList.PushFront(&Item{key: key, value: value, Expiration: e})
	c.cache[key] = ele
	if c.maxItemSize != 0 && c.cacheList.Len() > c.maxItemSize {
		c.RemoveOldest()
	}
}

//Delete delete the key
func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.cache == nil {
		return
	}
	if ele, ok := c.cache[key]; ok {
		c.cacheList.Remove(ele)
		key := ele.Value.(*Item).key
		delete(c.cache, key)
		return
	}
}

//RemoveOldest remove the oldest key
func (c *Cache) RemoveOldest() {
	if c.cache == nil {
		return
	}
	ele := c.cacheList.Back()
	if ele != nil {
		c.cacheList.Remove(ele)
		key := ele.Value.(*Item).key
		delete(c.cache, key)
	}
}

func (c *Cache) GetLen() int {
	return c.cacheList.Len()
}

//开启线程 定期清除过期的key
func (c *Cache) run() {
	tick := time.Tick(5 * time.Second)
	for {
		select {
		case <-tick:
			c.DeleteExpired()
		}
	}
}

func (c *Cache) DeleteExpired() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	now := time.Now().UnixNano()
	for _, v := range c.cache {
		if v != nil && v.Value.(*Item).Expiration > 0 && now > v.Value.(*Item).Expiration {
			c.cacheList.Remove(v)
			key := v.Value.(*Item).key
			delete(c.cache, key)
		}
	}
}

// Add atomically adds n to i.
func (i *AtomicInt) Add(n int64) {
	atomic.AddInt64((*int64)(i), n)
}

// Get atomically gets the value of i.
func (i *AtomicInt) Get() int64 {
	return atomic.LoadInt64((*int64)(i))
}
