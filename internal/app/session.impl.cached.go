package app

import (
	"context"
	"fmt"

	"github.com/MaratBR/openlibrary/internal/app/cache"
	"github.com/gofrs/uuid"
)

type cachedSessionService struct {
	inner SessionService
	cache *cache.Cache
}

// GetByUserID implements SessionService.
func (c *cachedSessionService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]SessionInfo, error) {
	return c.inner.GetByUserID(ctx, userID)
}

func (c *cachedSessionService) Create(ctx context.Context, command CreateSessionCommand) (*SessionInfo, error) {
	return c.inner.Create(ctx, command)
}

// GetBySID implements SessionService.
func (c *cachedSessionService) GetBySID(ctx context.Context, sessionID string) (*SessionInfo, error) {
	return cache.GetOrSet[*SessionInfo](*c.cache, fmt.Sprintf("session:%s", sessionID), func(entry *cache.CacheEntry) (*SessionInfo, error) {
		session, err := c.inner.GetBySID(ctx, sessionID)
		if err != nil {
			return nil, err
		}

		entry.Expiration = session.ExpiresAt
		return session, nil
	})
}

// Renew implements SessionService.
func (c *cachedSessionService) Renew(ctx context.Context, command RenewSessionCommand) (*SessionInfo, error) {
	session, err := c.inner.Renew(ctx, command)
	c.invalidate(command.SessionID)
	return session, err
}

// TerminateAllByUserID implements SessionService.
func (c *cachedSessionService) TerminateAllByUserID(ctx context.Context, userID uuid.UUID) error {
	err := c.TerminateAllByUserID(ctx, userID)
	c.invalidateByUserID(ctx, userID)
	return err
}

// TerminateBySID implements SessionService.
func (c *cachedSessionService) TerminateBySID(ctx context.Context, sessionID string) error {
	err := c.inner.TerminateBySID(ctx, sessionID)
	c.invalidate(sessionID)
	return err
}

func (c *cachedSessionService) invalidate(sid string) {
	c.cache.Remove(fmt.Sprintf("session:%s", sid))
}

func (c *cachedSessionService) invalidateByUserID(ctx context.Context, userID uuid.UUID) {
	sessions, err := c.GetByUserID(ctx, userID)
	if err != nil {
		return
	}
	for _, session := range sessions {
		c.invalidate(session.SessionID)
	}
}

func NewCachedSessionService(inner SessionService, cache *cache.Cache) SessionService {
	return &cachedSessionService{
		inner: inner,
		cache: cache,
	}
}
