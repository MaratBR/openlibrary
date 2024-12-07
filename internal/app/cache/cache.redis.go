package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	fallback CacheBackend
	redis    *redis.Client
}

// Get implements CacheBackend.
func (r *redisCache) Get(ctx context.Context, key string) ([]byte, error) {
	cmd := r.redis.Get(ctx, key)
	bytes, err := cmd.Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrCacheMiss
		} else {
			return nil, err
		}
	}
	return bytes, nil
}

// Put implements CacheBackend.
func (r *redisCache) Put(ctx context.Context, entry CacheEntry) error {
	return r.redis.Set(ctx, entry.Key, entry.Value, entry.Expiration.Sub(time.Now())).Err()
}

// Remove implements CacheBackend.
func (r *redisCache) Remove(ctx context.Context, key string) (bool, error) {
	affected, err := r.redis.Del(ctx, key).Result()
	return affected == 1, err
}

// Start implements CacheBackend.
func (r *redisCache) Start(ctx context.Context) error {
	cmd := r.redis.Ping(ctx)
	return cmd.Err()
}

// Stop implements CacheBackend.
func (r *redisCache) Stop(ctx context.Context) {
	r.redis.Close()
}

func NewRedisCacheBackend(
	url string,
	fallback CacheBackend,
) CacheBackend {
	return &redisCache{
		fallback: fallback,
		redis:    redis.NewClient(&redis.Options{Addr: url, Password: "", DB: 0}),
	}
}
