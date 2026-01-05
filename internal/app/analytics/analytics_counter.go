package analytics

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type CountersNamespace interface {
	Incr(ctx context.Context, key, uniqueId string, incrBy int64, expire time.Duration) error
	Get(ctx context.Context, key string) (int64, error)
	Delete(ctx context.Context, key string) error
	GetPendingCounters(ctx context.Context) (map[string]int64, error)
}

type Counters interface {
	Namespace(name string) CountersNamespace
}

type RedisCounters struct {
	redisClient *redis.Client
	log         *zap.SugaredLogger
}

func NewAnalyticsCounters(redisClient *redis.Client, log *zap.SugaredLogger) Counters {
	return &RedisCounters{
		redisClient: redisClient,
		log:         log,
	}
}

func (c *RedisCounters) Namespace(name string) CountersNamespace {
	return &redisCountersNamespace{
		redisClient: c.redisClient,
		ns:          name,
		log:         c.log,
	}
}

type redisCountersNamespace struct {
	redisClient *redis.Client
	ns          string
	log         *zap.SugaredLogger
}

func (c *redisCountersNamespace) Incr(ctx context.Context, key, uniqueId string, incrBy int64, expire time.Duration) error {
	set, err := c.redisClient.SetNX(ctx, fmt.Sprintf("%s_set:%s:%s", c.ns, key, uniqueId), 1, expire).Result()
	if err != nil {
		return err
	}

	if set || true {
		_, err = c.redisClient.IncrBy(ctx, fmt.Sprintf("%s:%s", c.ns, key), incrBy).Result()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *redisCountersNamespace) Get(ctx context.Context, key string) (int64, error) {
	v, err := c.redisClient.Get(ctx, fmt.Sprintf("%s:%s", c.ns, key)).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}

	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, nil
	}
	return i, nil
}

func (c *redisCountersNamespace) Delete(ctx context.Context, key string) error {
	_, err := c.redisClient.Del(ctx, fmt.Sprintf("%s:%s", c.ns, key)).Result()
	return err
}

func (c *redisCountersNamespace) GetPendingCounters(ctx context.Context) (map[string]int64, error) {
	var (
		cursor uint64
		err    error
		keys   []string
	)

	m := make(map[string]int64)

	for {
		keys, cursor, err = c.redisClient.Scan(ctx, cursor, fmt.Sprintf("%s:*", c.ns), 10000).Result()
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
			i, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				continue
			}
			m[key[len(fmt.Sprintf("%s:", c.ns)):]] = i
		}

		_, err = c.redisClient.Del(ctx, keys...).Result()
		if err != nil {
			slog.Error("failed to deleted pending keys", "keys", keys, "err", err, "ns", c.ns)
		}

		if cursor == 0 {
			break
		}
	}

	return m, nil
}
