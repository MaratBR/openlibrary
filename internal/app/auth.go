package app

import (
	"context"

	"github.com/gofrs/uuid"
)

var (
	ErrTypeSignUp    = AppErrors.NewType("signup")
	ErrUsernameTaken = ErrTypeSignUp.New("username is already taken")

	ErrTypeAuthenticationError = AppErrors.NewType("authentication", ErrTraitAuthorizationIssue)
	ErrInvalidCredentials      = ErrTypeAuthenticationError.New("invalid credentials")
)

type AuthService interface {
	SignIn(ctx context.Context, input SignInCommand) (SignInResult, error)
	CreateSessionForUser(ctx context.Context, userID uuid.UUID, userAgent, ip string) (string, error)
	SignOut(ctx context.Context, sessionID string) error
	EnsureAdminUserExists(ctx context.Context) error
}

type SignInCommand struct {
	Username  string
	Password  string
	UserAgent string
	IpAddress string
}

type SignInResult struct {
	SessionID string
}

type SignUpResult struct {
	Created                   bool
	CreatedUserID             uuid.UUID
	EmailTaken                bool
	EmailVerificationRequired bool
}
