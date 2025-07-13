package session

import (
	"context"

	"github.com/joomcode/errorx"
)

type Session interface {
	Get(ctx context.Context, key string) (string, bool)
	Put(ctx context.Context, key string, value string)
	Save(ctx context.Context) error
}

type Store interface {
	Get(ctx context.Context, id string) (Session, error)
}

var (
	errNs        = errorx.NewNamespace("_store")
	ErrNoSession = errNs.NewType("no_session").New("session not found")
)
