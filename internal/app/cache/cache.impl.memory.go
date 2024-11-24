package cache

import (
	"time"
)

type memoryCacheImpl struct {
	cache map[string]CacheEntry
}

// Remove implements CacheBackend.
func (m *memoryCacheImpl) Remove(key string) (bool, error) {
	if _, ok := m.cache[key]; ok {
		delete(m.cache, key)
		return true, nil
	}
	return false, nil
}

// Get implements CacheBackend.
func (m *memoryCacheImpl) Get(key string) (CacheEntry, error) {
	entry, ok := m.cache[key]
	if !ok {
		return CacheEntry{}, ErrCacheMiss
	}
	if entry.Expiration.Before(time.Now()) {
		delete(m.cache, key)
		return CacheEntry{}, nil
	}
	return entry, nil
}

// Set implements CacheBackend.
func (m *memoryCacheImpl) Put(entry CacheEntry) error {

	m.cache[entry.Key] = entry
	return nil
}

// Start implements CacheBackend.
func (m *memoryCacheImpl) Start() error {
	m.cache = make(map[string]CacheEntry)
	return nil
}

// Stop implements CacheBackend.
func (m *memoryCacheImpl) Stop() {
	m.cache = nil
}

func NewMemoryCacheBackend() CacheBackend {
	return &memoryCacheImpl{}
}
