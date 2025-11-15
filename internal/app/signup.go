package app

import "context"

var (
	SignUpErrors       = AppErrors.NewSubNamespace("signup")
	SignUpInvalidInput = SignUpErrors.NewType("invalid")
)

type SignUpCommand struct {
	Username  string
	Password  string
	Email     string
	UserAgent string
	IpAddress string
}

type SignUpService interface {
	SignUp(ctx context.Context, input SignUpCommand) (SignUpResult, error)
}
