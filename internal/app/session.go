package app

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

var (
	ErrSessionNotFound = AppErrors.NewType("session_not_found").New("session not found")
)

type SessionInfo struct {
	SessionID    string
	CreatedAt    time.Time
	ExpiresAt    time.Time
	UserID       uuid.UUID
	UserAgent    string
	IpAddress    string
	UserName     string
	UserJoinedAt time.Time
	UserRole     UserRole
}

type CreateSessionCommand struct {
	UserID    uuid.UUID
	UserAgent string
	IpAddress string
	ExpiresAt time.Time
}

type RenewSessionCommand struct {
	SessionID string
	UserAgent string
	IpAddress string
	ExpiresAt time.Time
}

type SessionService interface {
	GetBySID(ctx context.Context, sessionID string) (*SessionInfo, error)
	Create(ctx context.Context, command CreateSessionCommand) (*SessionInfo, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]SessionInfo, error)
	TerminateBySID(ctx context.Context, sessionID string) error
	TerminateAllByUserID(ctx context.Context, userID uuid.UUID) error
	Renew(ctx context.Context, command RenewSessionCommand) (*SessionInfo, error)
}
