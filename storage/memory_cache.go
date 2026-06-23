package storage

import (
	"context"
	"sync"
	"time"
)

type entry struct {
	value     []byte
	expiresAt int64 // unix nano, 0 = no expiry
}

type InMemoryCache struct {
	mu       sync.RWMutex
	data     map[string]entry
	stopChan chan struct{}
}

func NewInMemoryCache(cleanupInterval time.Duration) *InMemoryCache {
	c := &InMemoryCache{
		data:     make(map[string]entry),
		stopChan: make(chan struct{}),
	}

	if cleanupInterval > 0 {
		go c.startCleanup(cleanupInterval)
	}

	return c
}

func (c *InMemoryCache) Get(ctx context.Context, key string) ([]byte, bool, error) {
	c.mu.RLock()
	e, ok := c.data[key]
	c.mu.RUnlock()

	if !ok {
		return nil, false, nil
	}

	if e.expiresAt > 0 && time.Now().UnixNano() > e.expiresAt {
		// lazy delete
		c.mu.Lock()
		delete(c.data, key)
		c.mu.Unlock()
		return nil, false, nil
	}

	return e.value, true, nil
}

func (c *InMemoryCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	var expiresAt int64
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl).UnixNano()
	}

	c.mu.Lock()
	c.data[key] = entry{
		value:     value,
		expiresAt: expiresAt,
	}
	c.mu.Unlock()

	return nil
}

func (c *InMemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	delete(c.data, key)
	c.mu.Unlock()
	return nil
}

func (c *InMemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	_, ok, _ := c.Get(ctx, key)
	return ok, nil
}

func (c *InMemoryCache) Clear(ctx context.Context) error {
	c.mu.Lock()
	c.data = make(map[string]entry)
	c.mu.Unlock()
	return nil
}

func (c *InMemoryCache) startCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stopChan:
			return
		}
	}
}

func (c *InMemoryCache) cleanup() {
	now := time.Now().UnixNano()

	c.mu.Lock()
	for k, v := range c.data {
		if v.expiresAt > 0 && now > v.expiresAt {
			delete(c.data, k)
		}
	}
	c.mu.Unlock()
}
