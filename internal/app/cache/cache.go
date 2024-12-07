package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrCacheMiss = errors.New("cache miss")
)

func (c CacheEntry) copyInto(dest *CacheEntry) {
	dest.Key = c.Key
	dest.Value = c.Value // shallow copy, probably fine in this case
	dest.Expiration = c.Expiration
}

type CacheBackend interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context)
	Get(ctx context.Context, key string) ([]byte, error)
	Put(ctx context.Context, entry CacheEntry) error
	Remove(ctx context.Context, key string) (bool, error)
}

type Cache struct {
	backend CacheBackend
}

func New(backend CacheBackend) *Cache {
	return &Cache{backend: backend}
}

func (c Cache) PutRaw(ctx context.Context, entry CacheEntry) error {
	return c.backend.Put(ctx, entry)
}

func (c Cache) Get(ctx context.Context, key string) ([]byte, error) {
	return c.backend.Get(ctx, key)
}

func (c Cache) PutJSON(ctx context.Context, key string, value any, expiration time.Time) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.backend.Put(ctx, CacheEntry{
		Key:        key,
		Expiration: expiration,
		Value:      b,
	})
}

func (c Cache) GetJSON(ctx context.Context, key string, dest any) error {
	b, err := c.backend.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dest)
}

func (c *Cache) Remove(ctx context.Context, key string) (bool, error) {
	return c.backend.Remove(ctx, key)
}

func GetOrSet[T any](ctx context.Context, c *Cache, key string, getOrSetFn func(entry *CacheEntry) (T, error)) (T, error) {
	var value T
	err := c.GetJSON(ctx, key, &value)
	if err != nil {
		if errors.Is(err, ErrCacheMiss) {
			entry := CacheEntry{
				Key: key,
			}
			value, err = getOrSetFn(&entry)
			if err != nil {
				return value, err
			}
			err = c.PutJSON(ctx, key, value, entry.Expiration)
			if err != nil {
				return value, err
			}
			return value, nil
		} else {
			return value, err

		}
	} else {
		return value, nil
	}
}
