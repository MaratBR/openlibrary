package cache

import (
	"time"
)

type CacheEntry struct {
	Key        string
	Value      []byte
	Expiration time.Time
}
