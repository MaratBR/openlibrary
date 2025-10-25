package cache

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf/v2"
)

func CacheBackendFromConfig(cfg *koanf.Koanf) (CacheBackend, error) {
	backend := cfg.String("cache.type")
	switch backend {
	case "disabled":
		return NewDisabledCacheBackend(), nil
	case "memory":
		return NewMemoryCacheBackend(), nil
	case "redis":
		url := cfg.String("redis.url")
		if strings.Trim(url, " \n\t") == "" {
			return nil, fmt.Errorf("redis.url is empty")
		}

		return NewRedisCacheBackend(
			url,
			NewMemoryCacheBackend(),
		), nil
	default:
		return nil, fmt.Errorf("unknown cache backend: %s", backend)
	}
}
