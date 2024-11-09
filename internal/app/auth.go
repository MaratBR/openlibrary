package app

import (
	"context"
	"time"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

type AuthService struct {
	queries *store.Queries
	db      DB
}

func NewAuthService(db DB) *AuthService {
	return &AuthService{
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
		if err == pgx.ErrNoRows {
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
		CreatedAt: timeToTimestamptz(time.Now()),
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
		rollbackTx(ctx, tx)
		return SignUpResult{}, err
	}

	if exists {
		rollbackTx(ctx, tx)
		return SignUpResult{IsSuccess: false}, nil
	}

	userID, err := createUser(ctx, queries, input.Username, input.Password)
	if err != nil {
		rollbackTx(ctx, tx)
		return SignUpResult{}, err
	}

	sessionID, err := s.createNewSession(ctx, queries, userID, input.UserAgent, input.IpAddress)
	if err != nil {
		rollbackTx(ctx, tx)
		return SignUpResult{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return SignUpResult{}, err
	}
	return SignUpResult{IsSuccess: true, SessionID: sessionID}, nil
}

func createUser(ctx context.Context, queries *store.Queries, username, password string) (id uuid.UUID, err error) {
	userID := uuidV4()
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return
	}
	err = queries.InsertUser(ctx, store.InsertUserParams{
		ID:           uuidDomainToDb(userID),
		PasswordHash: hashedPassword,
		Name:         username,
		JoinedAt:     timeToTimestamptz(time.Now()),
	})
	if err == nil {
		id = userID
	}
	return
}

func (s AuthService) EnsureAdminUserExists(ctx context.Context) error {
	exists, err := s.queries.UserExistsByUsername(ctx, "admin")
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

	_, err = createUser(ctx, queries, "admin", "admin")
	if err != nil {
		rollbackTx(ctx, tx)
		return err
	}

	err = tx.Commit(ctx)
	return err
}

type SessionInfo struct {
	SessionID    string
	CreatedAt    time.Time
	ExpiresAt    time.Time
	UserID       uuid.UUID
	UserAgent    string
	IpAddress    string
	UserName     string
	UserJoinedAt time.Time
}

func (s AuthService) GetUserBySessionID(ctx context.Context, sessionID string) (*SessionInfo, error) {
	result, err := s.queries.GetSessionInfo(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return &SessionInfo{
		SessionID:    result.ID,
		CreatedAt:    timeDbToDomain(result.CreatedAt),
		ExpiresAt:    timeDbToDomain(result.ExpiresAt),
		UserID:       uuidDbToDomain(result.UserID),
		UserAgent:    result.UserAgent,
		IpAddress:    result.IpAddress,
		UserName:     result.UserName,
		UserJoinedAt: timeDbToDomain(result.UserJoinedAt),
	}, nil
}
