package app

import "context"

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
