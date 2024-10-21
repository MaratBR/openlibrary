package app

import (
	"context"
	"database/sql"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

type AuthService struct {
	queries *store.Queries
	db      *pgx.Conn
}

func NewAuthService(db *pgx.Conn) AuthService {
	return AuthService{
		queries: store.New(db),
		db:      db,
	}
}

type SignInCommand struct {
	Username string
	Password string
}

type SignInResult struct {
	IsSuccess bool
	SessionID string
}

// SignIn signs in user and returns a session ID.
//
// If the user is not found by the provided username, or the password doesn't match,
// the function returns SignInResult with IsSuccess set to false, and the error is nil.
//
// If the user is found and the password matches, the function creates a new session
// and returns SignInResult with IsSuccess set to true, and the session ID.
func (s AuthService) SignIn(ctx context.Context, input SignInCommand) (SignInResult, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return SignInResult{}, err
	}
	queries := s.queries.WithTx(tx)
	user, err := queries.FindUserByUsername(ctx, input.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return SignInResult{
				IsSuccess: false}, nil
		}
		return SignInResult{}, err
	}

	match, err := verifyPassword(input.Password, user.PasswordHash)
	if err != nil {
		return SignInResult{}, err
	}

	if !match {
		return SignInResult{IsSuccess: false}, nil
	}

	sessionID, err := s.createNewSession(ctx, queries, uuidDbToDomain(user.ID), "Lol kek, not chrome :)", "123.123.123.123")
	if err != nil {
		return SignInResult{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return SignInResult{}, err
	}

	return SignInResult{IsSuccess: true, SessionID: sessionID}, nil
}

func (s AuthService) createNewSession(ctx context.Context, queries *store.Queries, userID uuid.UUID, userAgent, ip string) (string, error) {
	id := genOpaqueID()
	err := queries.InsertSession(ctx, store.InsertSessionParams{
		ID:        id,
		UserID:    uuidDomainToDb(userID),
		UserAgent: userAgent,
		IpAddress: ip,
		ExpiresAt: timeToTimestamptz(time.Now().Add(90 * 24 * time.Hour)),
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

type SignUpCommand struct {
	Username  string
	Password  string
	UserAgent string
	IpAddress string
}

type SignUpResult struct {
	IsSuccess bool
	SessionID string
}

// SignUp creates a new user and a new session, and returns a session ID.
//
// The function will return SignUpResult with IsSuccess set to false if the user
// with the provided username already exists. The error will be nil in this case.
//
// If the user is created and the session is created, the function returns
// SignUpResult with IsSuccess set to true, and the session ID.
func (s AuthService) SignUp(ctx context.Context, input SignUpCommand) (SignUpResult, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return SignUpResult{}, err
	}
	queries := s.queries.WithTx(tx)
	exists, err := queries.UserExistsByUsername(ctx, input.Username)
	if err != nil {
		return SignUpResult{}, err
	}

	if exists {
		return SignUpResult{IsSuccess: false}, nil
	}

	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		return SignUpResult{}, err
	}
	userID := uuidV4()
	err = queries.InsertUser(ctx, store.InsertUserParams{
		ID:           uuidDomainToDb(userID),
		PasswordHash: hashedPassword,
		Name:         input.Username,
		JoinedAt:     timeToTimestamptz(time.Now()),
	})
	if err != nil {
		return SignUpResult{}, err
	}
	sessionID, err := s.createNewSession(ctx, queries, userID, input.UserAgent, input.IpAddress)
	if err != nil {
		return SignUpResult{}, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return SignUpResult{}, err
	}
	return SignUpResult{IsSuccess: true, SessionID: sessionID}, nil
}
