package gocache

import "time"

const (
	DefaultExpiration = 0
)

type CacheStatus struct {
	Gets        int64
	Hits        int64
	MaxItemSize int
	CurrentSize int
}

type CacheInterface interface {
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
	Delete(key string)
	Status() *CacheStatus
	DeleteExpired()
}

type Item struct {
	value     interface{}
	Expiration int64
	key  interface{}
}


// Returns true if the item has expired.
func (item Item) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > item.Expiration
}
