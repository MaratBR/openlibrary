package cache

import (
	"context"
	"time"
)

type memoryCacheImpl struct {
	cache map[string]CacheEntry
}

// Remove implements CacheBackend.
func (m *memoryCacheImpl) Remove(ctx context.Context, key string) (bool, error) {
	if _, ok := m.cache[key]; ok {
		delete(m.cache, key)
		return true, nil
	}
	return false, nil
}

// Get implements CacheBackend.
func (m *memoryCacheImpl) Get(ctx context.Context, key string) ([]byte, error) {
	entry, ok := m.cache[key]
	if !ok {
		return nil, ErrCacheMiss
	}
	if entry.Expiration.Before(time.Now()) {
		delete(m.cache, key)
		return nil, nil
	}
	return entry.Value, nil
}

// Set implements CacheBackend.
func (m *memoryCacheImpl) Put(ctx context.Context, entry CacheEntry) error {

	m.cache[entry.Key] = entry
	return nil
}

// Start implements CacheBackend.
func (m *memoryCacheImpl) Start(ctx context.Context) error {
	m.cache = make(map[string]CacheEntry)
	return nil
}

// Stop implements CacheBackend.
func (m *memoryCacheImpl) Stop(ctx context.Context) {
	m.cache = nil
}

func NewMemoryCacheBackend() CacheBackend {
	return &memoryCacheImpl{}
}
