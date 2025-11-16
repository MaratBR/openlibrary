package app

import (
	"context"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

type authService struct {
	queries  *store.Queries
	sessions SessionService
	db       DB
}

// SignOut implements AuthService.
func (s *authService) SignOut(ctx context.Context, sessionID string) error {
	return s.sessions.TerminateBySID(ctx, sessionID)
}

func createUser(ctx context.Context, queries *store.Queries, username, email, password string, role UserRole) (id uuid.UUID, err error) {
	userID := uuidV4()
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return
	}
	err = queries.User_Insert(ctx, store.User_InsertParams{
		ID:           uuidDomainToDb(userID),
		PasswordHash: hashedPassword,
		Name:         username,
		JoinedAt:     timeToTimestamptz(time.Now()),
		Email:        email,
	})
	if err == nil {
		id = userID
	} else {
		err = wrapUnexpectedDBError(err)
	}

	if role != RoleUser {
		err = queries.UpdateUserRole(ctx, store.UpdateUserRoleParams{
			Role: store.UserRole(role),
			ID:   uuidDomainToDb(userID),
		})
		if err != nil {
			err = wrapUnexpectedDBError(err)
		}
	}

	return
}

func (s *authService) EnsureAdminUserExists(ctx context.Context) error {
	exists, err := s.queries.User_ExistsByUsername(ctx, "admin")
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	queries := s.queries.WithTx(tx)
	if err != nil {
		rollbackTx(ctx, tx)
		return err
	}

	_, err = createUser(ctx, queries, "admin", "", "admin", RoleAdmin)
	if err != nil {
		rollbackTx(ctx, tx)
		return err
	}

	err = tx.Commit(ctx)
	return err
}

// SignIn signs in user and returns a session ID.
//
// If the user is not found by the provided username, or the password doesn't match,
// the function returns SignInResult with IsSuccess set to false, and the error is nil.
//
// If the user is found and the password matches, the function creates a new session
// and returns SignInResult with IsSuccess set to true, and the session ID.
func (s *authService) SignIn(ctx context.Context, input SignInCommand) (SignInResult, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return SignInResult{}, err
	}
	queries := s.queries.WithTx(tx)
	user, err := queries.User_FindByLogin(ctx, input.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return SignInResult{}, ErrInvalidCredentials
		}
		return SignInResult{}, err
	}

	match, err := verifyPassword(input.Password, user.PasswordHash)
	if err != nil {
		return SignInResult{}, err
	}

	if !match {
		return SignInResult{}, ErrInvalidCredentials
	}

	sessionID, err := s.createNewSession(ctx, queries, uuidDbToDomain(user.ID), input.UserAgent, input.IpAddress)
	if err != nil {
		return SignInResult{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return SignInResult{}, err
	}

	return SignInResult{SessionID: sessionID}, nil
}

func (s *authService) CreateSessionForUser(ctx context.Context, userID uuid.UUID, userAgent, ip string) (string, error) {
	return s.createNewSession(ctx, s.queries, userID, userAgent, ip)
}

func (s *authService) createNewSession(ctx context.Context, queries *store.Queries, userID uuid.UUID, userAgent, ip string) (string, error) {
	session, err := s.sessions.Create(ctx, CreateSessionCommand{
		UserID:    userID,
		UserAgent: userAgent,
		IpAddress: ip,
		ExpiresAt: time.Now().Add(90 * 24 * time.Hour),
	})
	if err != nil {
		return "", err
	}
	return session.SID, nil
}

func NewAuthService(db DB, sessions SessionService) AuthService {
	return &authService{
		queries:  store.New(db),
		db:       db,
		sessions: sessions,
	}
}
