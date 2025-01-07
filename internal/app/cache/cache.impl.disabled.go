package cache

import "context"

type disabledCacheImpl struct {
}

// Get implements CacheBackend.
func (d *disabledCacheImpl) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, ErrCacheMiss
}

// Put implements CacheBackend.
func (d *disabledCacheImpl) Put(ctx context.Context, entry CacheEntry) error {
	return nil
}

// Remove implements CacheBackend.
func (d *disabledCacheImpl) Remove(ctx context.Context, key string) (bool, error) {
	return false, nil
}

// Start implements CacheBackend.
func (d *disabledCacheImpl) Start(ctx context.Context) error {
	return nil
}

// Stop implements CacheBackend.
func (d *disabledCacheImpl) Stop(ctx context.Context) {
}

func NewDisabledCacheBackend() CacheBackend {
	return &disabledCacheImpl{}
}
