package app

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/mileusna/useragent"
)

var (
	ErrSessionNotFound = AppErrors.NewType("session_not_found").New("session not found")
)

type SessionInfo struct {
	ID           int64
	SID          string
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

type UserAgentInfo struct {
	Value   string `json:"value"`
	Version string `json:"ver"`
	Name    string `json:"name"`
	OS      string `json:"os"`
}

type SessionPublicInfo struct {
	ID        int64         `json:"id,string"`
	UserAgent UserAgentInfo `json:"userAgent"`
	IpAddress string        `json:"ipAddress"`
	CreatedAt time.Time     `json:"createdAt"`
	ExpiresAt time.Time     `json:"expiresAt"`
}

func NewSessionPublicInfo(session SessionInfo) SessionPublicInfo {
	ua := useragent.Parse(session.UserAgent)

	return SessionPublicInfo{
		ID: session.ID,
		UserAgent: UserAgentInfo{
			Value:   session.UserAgent,
			OS:      ua.OS,
			Name:    ua.Name,
			Version: ua.Version,
		},
		IpAddress: session.IpAddress,
		CreatedAt: session.CreatedAt.UTC(),
		ExpiresAt: session.ExpiresAt.UTC(),
	}
}

func SessionPublicInfoArray(sessions []SessionInfo) []SessionPublicInfo {
	result := make([]SessionPublicInfo, len(sessions))
	for i, session := range sessions {
		result[i] = NewSessionPublicInfo(session)
	}
	return result
}

type SessionService interface {
	GetBySID(ctx context.Context, sessionID string) (*SessionInfo, error)
	Create(ctx context.Context, command CreateSessionCommand) (*SessionInfo, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]SessionInfo, error)
	TerminateBySID(ctx context.Context, sessionID string) error
	TerminateAllByUserID(ctx context.Context, userID uuid.UUID) error
	Renew(ctx context.Context, command RenewSessionCommand) (*SessionInfo, error)
}
