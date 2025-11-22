package app

import (
	"context"
	"time"

	"github.com/MaratBR/openlibrary/internal/commonutil"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
)

type sessionService struct {
	db      DB
	queries *store.Queries
}

// GetByUserID implements SessionService.
func (s *sessionService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]SessionInfo, error) {
	sessions, err := s.queries.Session_GetUserSessions(ctx, uuidDomainToDb(userID))
	if err != nil {
		if err == store.ErrNoRows {
			return nil, ErrSessionNotFound
		}
		return nil, wrapUnexpectedDBError(err)
	}
	return MapSlice(sessions, func(s store.Session_GetUserSessionsRow) SessionInfo {
		return SessionInfo{
			ID:           s.ID,
			SID:          s.Sid,
			CreatedAt:    timeDbToDomain(s.CreatedAt),
			ExpiresAt:    timeDbToDomain(s.ExpiresAt),
			UserID:       uuidDbToDomain(s.UserID),
			UserAgent:    s.UserAgent,
			IpAddress:    s.IpAddress,
			UserName:     s.UserName,
			UserJoinedAt: timeDbToDomain(s.UserJoinedAt),
		}
	}), nil
}

// Create implements SessionService.
func (s *sessionService) Create(ctx context.Context, command CreateSessionCommand) (*SessionInfo, error) {
	sessionID, err := commonutil.GenerateRandomStringURLSafe(32)
	if err != nil {
		return nil, err
	}

	err = s.queries.Session_Insert(ctx, store.Session_InsertParams{
		ID:        GenID(),
		Sid:       sessionID,
		UserID:    uuidDomainToDb(command.UserID),
		UserAgent: command.UserAgent,
		IpAddress: command.IpAddress,
		ExpiresAt: timeToTimestamptz(time.Now().Add(90 * 24 * time.Hour)),
		CreatedAt: timeToTimestamptz(time.Now()),
	})
	if err != nil {
		return nil, wrapUnexpectedDBError(err)
	}

	session, err := s.get(ctx, sessionID)
	return &session, err
}

// GetBySID implements SessionService.
func (s *sessionService) GetBySID(ctx context.Context, sessionID string) (*SessionInfo, error) {
	sessionInfo, err := s.get(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return &sessionInfo, nil
}

func (s *sessionService) get(ctx context.Context, sessionID string) (SessionInfo, error) {
	result, err := s.queries.Session_GetInfo(ctx, sessionID)
	if err != nil {
		if err == store.ErrNoRows {
			return SessionInfo{}, ErrSessionNotFound
		}
		return SessionInfo{}, wrapUnexpectedDBError(err)
	}

	return SessionInfo{
		ID:           result.ID,
		SID:          result.Sid,
		CreatedAt:    timeDbToDomain(result.CreatedAt),
		ExpiresAt:    timeDbToDomain(result.ExpiresAt),
		UserID:       uuidDbToDomain(result.UserID),
		UserAgent:    result.UserAgent,
		IpAddress:    result.IpAddress,
		UserName:     result.UserName,
		UserJoinedAt: timeDbToDomain(result.UserJoinedAt),
		UserRole:     UserRole(result.UserRole),
	}, nil
}

// Renew implements SessionService.
func (s *sessionService) Renew(ctx context.Context, command RenewSessionCommand) (*SessionInfo, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	queries := s.queries.WithTx(tx)

	session, err := queries.Session_GetInfo(ctx, command.SessionID)
	if err != nil {
		rollbackTx(ctx, tx)
		if err == store.ErrNoRows {
			return nil, ErrSessionNotFound
		}
		return nil, wrapUnexpectedDBError(err)
	}
	err = queries.Session_Terminate(ctx, command.SessionID)
	if err != nil {
		rollbackTx(ctx, tx)
		return nil, wrapUnexpectedDBError(err)
	}

	err = queries.Session_Insert(ctx, store.Session_InsertParams{
		ID:        GenID(),
		Sid:       command.SessionID,
		UserID:    session.UserID,
		UserAgent: command.UserAgent,
		IpAddress: command.IpAddress,
		ExpiresAt: timeToTimestamptz(time.Now().Add(90 * 24 * time.Hour)),
		CreatedAt: timeToTimestamptz(time.Now()),
	})
	if err != nil {
		rollbackTx(ctx, tx)
		return nil, wrapUnexpectedDBError(err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	si, err := s.get(ctx, command.SessionID)
	if err != nil {
		return nil, err
	}

	return &si, nil
}

// TerminateAllByUserID implements SessionService.
func (s *sessionService) TerminateAllByUserID(ctx context.Context, userID uuid.UUID) error {
	err := s.queries.Session_TerminateAllByUserID(ctx, uuidDomainToDb(userID))
	if err != nil {
		return wrapUnexpectedDBError(err)
	}
	return nil
}

// TerminateBySID implements SessionService.
func (s *sessionService) TerminateBySID(ctx context.Context, sessionID string) error {
	err := s.queries.Session_Terminate(ctx, sessionID)
	if err != nil {
		return wrapUnexpectedDBError(err)
	}
	return nil
}

func NewSessionService(db DB) SessionService {
	return &sessionService{
		db:      db,
		queries: store.New(db),
	}
}
