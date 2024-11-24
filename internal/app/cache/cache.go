package cache

import (
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
	Start() error
	Stop()
	Get(key string) (CacheEntry, error)
	Put(entry CacheEntry) error
	Remove(key string) (bool, error)
}

type Cache struct {
	backend CacheBackend
}

func New(backend CacheBackend) Cache {
	return Cache{backend: backend}
}

func (c Cache) PutRaw(entry CacheEntry) error {
	return c.backend.Put(entry)
}

func (c Cache) GetRaw(key string) (CacheEntry, error) {
	return c.backend.Get(key)
}

func (c Cache) PutJSON(key string, value any, expiration time.Time) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.backend.Put(CacheEntry{
		Key:        key,
		Expiration: expiration,
		Value:      b,
	})
}

func (c Cache) GetJSON(key string, dest any) error {
	b, err := c.backend.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(b.Value, dest)
}

func (c *Cache) Remove(key string) (bool, error) {
	return c.backend.Remove(key)
}

func GetOrSet[T any](c Cache, key string, getOrSetFn func(entry *CacheEntry) (T, error)) (T, error) {
	var value T
	err := c.GetJSON(key, &value)
	if err != nil {
		if errors.Is(err, ErrCacheMiss) {
			entry := CacheEntry{
				Key: key,
			}
			value, err = getOrSetFn(&entry)
			if err != nil {
				return value, err
			}
			err = c.PutJSON(key, value, entry.Expiration)
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
