package analytics

import (
	"context"
	"time"

	"go.uber.org/fx"
)

// Atomic is its essense just a global store of booleans,
// allowing you to check if a specific event has occured
// within given time and "reserve" it it should primarily
// be used to prevent duplicated events happening within
// a relatively short period of time
type Atomic interface {
	TryAcquire(
		ctx context.Context,
		key string,
		ttl time.Duration,
	) (bool, error)

	Release(ctx context.Context, key string) error
}

type atomicPrefixed struct {
	inner  Atomic
	prefix string
}

// Release implements [Atomic].
func (a *atomicPrefixed) Release(ctx context.Context, key string) error {
	return a.inner.Release(ctx, a.prefix+key)
}

// TryAcquire implements [Atomic].
func (a *atomicPrefixed) TryAcquire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return a.inner.TryAcquire(ctx, a.prefix+key, ttl)
}

func NewPrefixedAtomic(inner Atomic, prefix string) Atomic {
	return &atomicPrefixed{
		inner:  inner,
		prefix: prefix,
	}
}

var atomicModule = fx.Module("ol_analytics_atomic", fx.Provide(
	NewAtomic,
))
