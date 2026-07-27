package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type atomicStore struct {
	redisClient *redis.Client
}

// Release implements [Atomic].
func (a *atomicStore) Release(ctx context.Context, key string) error {
	_, err := a.redisClient.Del(ctx, a.key(key)).Result()
	return err
}

// TryAcquire implements [Atomic].
func (a *atomicStore) TryAcquire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	set, err := a.redisClient.SetNX(ctx, a.key(key), 1, ttl).Result()
	return set, err
}

func (a *atomicStore) key(key string) string {
	return fmt.Sprintf("atom:%s", key)
}

func NewAtomic(redisClient *redis.Client) Atomic {
	return &atomicStore{redisClient: redisClient}
}
