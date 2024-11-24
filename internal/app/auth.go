package app

import (
	"context"
)

var (
	ErrTypeSignUp    = AppErrors.NewType("signup")
	ErrUsernameTaken = ErrTypeSignUp.New("username is already taken")

	ErrTypeAuthenticationError = AppErrors.NewType("authentication", ErrTraitAuthorizationIssue)
	ErrInvalidCredentials      = ErrTypeAuthenticationError.New("invalid credentials")
)

type AuthService interface {
	SignIn(ctx context.Context, input SignInCommand) (SignInResult, error)
	SignUp(ctx context.Context, input SignUpCommand) (SignUpResult, error)
	SignOut(ctx context.Context, sessionID string) error
	EnsureAdminUserExists(ctx context.Context) error
}

type SignInCommand struct {
	Username string
	Password string
}

type SignInResult struct {
	SessionID string
}

type SignUpCommand struct {
	Username  string
	Password  string
	UserAgent string
	IpAddress string
}

type SignUpResult struct {
	SessionID string
}
