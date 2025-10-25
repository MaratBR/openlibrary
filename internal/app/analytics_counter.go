package app

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type AnalyticsCounters interface {
	Incr(ctx context.Context, key, uniqueId string, incrBy int64, expire time.Duration) error
	Get(ctx context.Context, key string) (int64, error)
	Delete(ctx context.Context, key string) error
}

type RedisAnalyticsCounters struct {
	redisClient *redis.Client
}

func NewAnalyticsCounters(redisClient *redis.Client) AnalyticsCounters {
	return &RedisAnalyticsCounters{
		redisClient: redisClient,
	}
}

func (c *RedisAnalyticsCounters) Incr(ctx context.Context, key, uniqueId string, incrBy int64, expire time.Duration) error {
	set, err := c.redisClient.SetNX(ctx, fmt.Sprintf("viewed:%s:%s", key, uniqueId), 1, expire).Result()
	if err != nil {
		return err
	}

	if set {
		_, err = c.redisClient.IncrBy(ctx, fmt.Sprintf("views:%s", key), incrBy).Result()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *RedisAnalyticsCounters) Get(ctx context.Context, key string) (int64, error) {
	v, err := c.redisClient.Get(ctx, fmt.Sprintf("views:%s", key)).Result()
	if err != nil {
		return 0, err
	}

	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return i, nil
	}
	return 0, nil
}

func (c *RedisAnalyticsCounters) Delete(ctx context.Context, key string) error {
	_, err := c.redisClient.Del(ctx, fmt.Sprintf("views:%s", key)).Result()
	return err
}
