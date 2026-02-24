package cache

import (
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
)

// Cache wraps bigcache for simple key-value storage.
type Cache struct {
	bc *bigcache.BigCache
}

// New creates a cache with the given TTL.
func New(ttl time.Duration) *Cache {
	cfg := bigcache.DefaultConfig(ttl)
	bc, _ := bigcache.New(context.Background(), cfg)
	return &Cache{bc: bc}
}

func (c *Cache) Set(key string, value []byte) error {
	return c.bc.Set(key, value)
}

func (c *Cache) Get(key string) ([]byte, error) {
	return c.bc.Get(key)
}

func (c *Cache) Delete(key string) error {
	return c.bc.Delete(key)
}
