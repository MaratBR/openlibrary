package cache

type disabledCacheImpl struct {
	cache map[string]CacheEntry
}

// Remove implements CacheBackend.
func (m *disabledCacheImpl) Remove(key string) (bool, error) {
	return false, nil
}

// Get implements CacheBackend.
func (m *disabledCacheImpl) Get(key string) (CacheEntry, error) {
	return CacheEntry{}, ErrCacheMiss
}

// Set implements CacheBackend.
func (m *disabledCacheImpl) Put(entry CacheEntry) error {
	return nil
}

// Start implements CacheBackend.
func (m *disabledCacheImpl) Start() error {
	return nil
}

// Stop implements CacheBackend.
func (m *disabledCacheImpl) Stop() {
}

func NewDisabledCacheBackend() CacheBackend {
	return &memoryCacheImpl{}
}
