package app

import (
	"context"
	"time"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/gofrs/uuid"
)

var (
	SignUpErrors                       = apperror.AppErrors.NewSubNamespace("signup")
	SignUpInvalidInput                 = SignUpErrors.NewType("invalid")
	SignUpEmailVerificationErrors      = SignUpErrors.NewSubNamespace("email_verification")
	SignUpEmailVerificationNA          = SignUpEmailVerificationErrors.NewType("na")
	SignUpEmailVerificationRateLimit   = SignUpEmailVerificationErrors.NewType("rate_limit")
	SignUpEmailVerificationTimedOut    = SignUpEmailVerificationErrors.NewType("timedout")
	SignUpEmailVerificationInvalidCode = SignUpEmailVerificationErrors.NewType("invalid_code")
)

type SignUpCommand struct {
	Username                  string
	Password                  string
	Email                     string
	UserAgent                 string
	IpAddress                 string
	BypassEmailRequirement    bool
	BypassPasswordRequirement bool
}

type VerifyEmailCommand struct {
	UserID uuid.UUID
	Code   string
}

type SendEmailVerificationCommand struct {
	UserID          uuid.UUID
	BypassRateLimit bool
}

type EmailVerificationStatus struct {
	CanSendAgainAfter Nullable[time.Time]
	WasSent           bool
}

type SendEmailVerificationResult struct {
	CanResendAfter time.Time
}

type SignUpService interface {
	SignUp(ctx context.Context, cmd SignUpCommand) (SignUpResult, error)
	VerifyEmail(ctx context.Context, cmd VerifyEmailCommand) error
	SendEmailVerification(ctx context.Context, cmd SendEmailVerificationCommand) (SendEmailVerificationResult, error)
	GetEmailVerificationStatus(ctx context.Context, userID uuid.UUID) (EmailVerificationStatus, error)
}
