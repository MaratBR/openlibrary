package app

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type AnalyticsCounters interface {
	Incr(ctx context.Context, key, uniqueId string, incrBy int64, expire time.Duration) error
	Get(ctx context.Context, key string) (int64, error)
	Delete(ctx context.Context, key string) error
	GetPendingCounters(ctx context.Context) (map[string]int64, error)
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
		return 0, nil
	}
	return i, nil
}

func (c *RedisAnalyticsCounters) Delete(ctx context.Context, key string) error {
	_, err := c.redisClient.Del(ctx, fmt.Sprintf("views:%s", key)).Result()
	return err
}

func (c *RedisAnalyticsCounters) GetPendingCounters(ctx context.Context) (map[string]int64, error) {
	var (
		cursor uint64
		err    error
		keys   []string
	)

	m := make(map[string]int64)

	for {
		keys, cursor, err = c.redisClient.Scan(ctx, cursor, "views:*", 10000).Result()
		if err != nil {
			return nil, err
		}

		if len(keys) == 0 {
			break
		}

		values, err := c.redisClient.MGet(ctx, keys...).Result()
		if err != nil {
			return nil, err
		}

		for i, key := range keys {
			value := values[i]
			str, ok := value.(string)
			if !ok {
				continue
			}
			i, err := strconv.ParseInt(str, 10, 62)
			if err != nil {
				continue
			}
			m[key[len("views:"):]] = i
		}

		_, err = c.redisClient.Del(ctx, keys...).Result()
		if err != nil {
			slog.Error("failed to deleted pending views keys", "keys", keys, "err", err)
		}

		if cursor == 0 {
			break
		}
	}

	return m, nil
}
