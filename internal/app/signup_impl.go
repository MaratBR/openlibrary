package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/app/email"
	"github.com/MaratBR/openlibrary/internal/commonutil"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
	"github.com/knadh/koanf/v2"
)

type signUpService struct {
	db                     DB
	cfg                    *koanf.Koanf
	siteConfig             *SiteConfig
	emailService           email.Service
	emailCodeValidDuration time.Duration
	emailCodeResendAfter   time.Duration
}

func NewSignUpService(db DB, cfg *koanf.Koanf, siteConfig *SiteConfig, emailService email.Service) SignUpService {
	return &signUpService{
		db:                     db,
		cfg:                    cfg,
		siteConfig:             siteConfig,
		emailService:           emailService,
		emailCodeValidDuration: time.Hour * 12,
		emailCodeResendAfter:   time.Minute * 2,
	}
}

// SignUp creates a new user and a new session, and returns a session ID.
//
// The function will return SignUpResult with IsSuccess set to false if the user
// with the provided username already exists. The error will be nil in this case.
//
// If the user is created and the session is created, the function returns
// SignUpResult with IsSuccess set to true, and the session ID.
func (s *signUpService) SignUp(ctx context.Context, input SignUpCommand) (SignUpResult, error) {
	// VALIDATION
	if s.cfg.Bool("auth.requireEmail") && input.Email == "" && !input.BypassEmailRequirement {
		return SignUpResult{}, SignUpInvalidInput.New("email is required")
	}
	if err := ValidateUserName(input.Username); err != nil {
		return SignUpResult{}, SignUpInvalidInput.Wrap(err, "invalid username")
	}
	if !input.BypassPasswordRequirement {
		passwordRequirements := s.siteConfig.Get().PasswordRequirements
		if err := ValidatePassword(input.Password, passwordRequirements); err != nil {
			return SignUpResult{}, SignUpInvalidInput.Wrap(err, "invalid password")
		}
	}
	if input.Email != "" {
		if err := ValidateEmail(input.Email); err != nil {
			return SignUpResult{}, SignUpInvalidInput.Wrap(err, "invalid email")
		}
	}

	// EMAIL VERIFICATION

	var emailVerificationRequired bool
	if s.cfg.Bool("auth.emailVerification") && input.Email != "" {
		emailVerificationRequired = true
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return SignUpResult{}, err
	}
	queries := store.New(s.db).WithTx(tx)

	if input.Email != "" {
		userWithSameEmailExists, err := queries.User_ExistsByEmail(ctx, input.Username)
		if err != nil {
			rollbackTx(ctx, tx)
			return SignUpResult{}, apperror.WrapUnexpectedDBError(err)
		}

		if userWithSameEmailExists {
			return SignUpResult{EmailTaken: true}, nil
		}
	}

	userWithSameNameExists, err := queries.User_ExistsByUsername(ctx, input.Username)
	if err != nil {
		rollbackTx(ctx, tx)
		return SignUpResult{}, apperror.WrapUnexpectedDBError(err)
	}

	if userWithSameNameExists {
		rollbackTx(ctx, tx)
		return SignUpResult{}, ErrUsernameTaken
	}

	userID, err := createUser(ctx, queries, input.Username, input.Email, input.Password, RoleUser, false)
	if err != nil {
		rollbackTx(ctx, tx)
		return SignUpResult{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return SignUpResult{}, err
	}

	if emailVerificationRequired {
		_, err := s.verifyEmailRequest(ctx, input.Email, userID)
		if err != nil {
			return SignUpResult{}, err
		}
	}

	return SignUpResult{CreatedUserID: userID, EmailVerificationRequired: emailVerificationRequired, Created: true}, nil
}

func (s *signUpService) verifyEmailRequest(ctx context.Context, email string, userID uuid.UUID) (time.Time, error) {
	hash, code, err := newEmailVerificationCode()
	if err != nil {
		return time.Time{}, err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return time.Time{}, apperror.WrapUnexpectedDBError(err)
	}
	queries := store.New(s.db).WithTx(tx)
	err = queries.EmailVerification_Delete(ctx, email)
	if err != nil {
		rollbackTx(ctx, tx)
		return time.Time{}, apperror.WrapUnexpectedDBError(err)
	}
	err = queries.EmailVerification_Insert(ctx, store.EmailVerification_InsertParams{
		UserID:               uuidDomainToDb(userID),
		Email:                email,
		ValidThrough:         timeToTimestamptz(time.Now().Add(s.emailCodeValidDuration)),
		VerificationCodeHash: hash,
	})
	if err != nil {
		rollbackTx(ctx, tx)
		return time.Time{}, apperror.WrapUnexpectedDBError(err)
	}

	// TODO some good looking message or something
	err = s.emailService.Send(ctx, email, "OL verification code", fmt.Sprintf("Your verification code is %s\n\nIf you did not sign up, please just ignore this message and apologies for disturbance :)", code))
	if err != nil {
		rollbackTx(ctx, tx)
		return time.Time{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return time.Time{}, apperror.WrapUnexpectedDBError(err)
	}

	return time.Now().Add(s.emailCodeResendAfter), nil
}

func newEmailVerificationCode() (hash string, code string, err error) {
	code, err = commonutil.GenerateRandomStringAlphabet(6, "0123456789ABCDEFGHJKLMOPRSTUVWXYZ")
	if err != nil {
		return
	}
	h := sha256.New()
	h.Write([]byte(code))
	hashBytes := h.Sum(nil)
	hash = hex.Dump(hashBytes)
	return
}

func verifyEmailCode(hash, code string) bool {
	h := sha256.New()
	h.Write([]byte(code))
	hashBytes := h.Sum(nil)
	return hex.Dump(hashBytes) == hash
}

func (s *signUpService) VerifyEmail(ctx context.Context, cmd VerifyEmailCommand) error {
	queries := store.New(s.db)
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}
	queries = queries.WithTx(tx)

	user, err := queries.User_Get(ctx, uuidDomainToDb(cmd.UserID))
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}

	if user.Email == "" {
		return SignUpEmailVerificationNA.New("this user does not have email and cannot be verified")
	}

	verification, err := queries.EmailVerification_Get(ctx, user.Email)
	if err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}

	if verification.ValidThrough.Time.Before(time.Now()) {
		return SignUpEmailVerificationTimedOut.New("email verification code is no longer valid")
	}

	if verifyEmailCode(verification.VerificationCodeHash, cmd.Code) {
		err = queries.User_SetEmailVerified(ctx, store.User_SetEmailVerifiedParams{
			ID:            uuidDomainToDb(cmd.UserID),
			EmailVerified: true,
		})
		if err != nil {
			return apperror.WrapUnexpectedDBError(err)
		}
		err = tx.Commit(ctx)
		if err != nil {
			return apperror.WrapUnexpectedDBError(err)
		}
		return nil
	} else {
		rollbackTx(ctx, tx)
		return SignUpEmailVerificationInvalidCode.New("invalid verification code")
	}

}

func (s *signUpService) SendEmailVerification(ctx context.Context, cmd SendEmailVerificationCommand) (SendEmailVerificationResult, error) {
	queries := store.New(s.db)

	user, err := queries.User_Get(ctx, uuidDomainToDb(cmd.UserID))
	if err != nil {
		return SendEmailVerificationResult{}, apperror.WrapUnexpectedDBError(err)
	}

	if user.Email == "" {
		return SendEmailVerificationResult{}, SignUpEmailVerificationNA.New("this user does not have email and cannot be verified")
	}

	verification, err := queries.EmailVerification_Get(ctx, user.Email)
	if err != nil {
		return SendEmailVerificationResult{}, apperror.WrapUnexpectedDBError(err)
	}

	if verification.UserID == user.ID && time.Now().Sub(verification.CreatedAt.Time) < s.emailCodeResendAfter && !cmd.BypassRateLimit {
		return SendEmailVerificationResult{}, SignUpEmailVerificationRateLimit.New("request for verification was rate limited, please wait")
	}

	canResendAfter, err := s.verifyEmailRequest(ctx, user.Email, cmd.UserID)
	if err != nil {
		return SendEmailVerificationResult{}, err
	}

	return SendEmailVerificationResult{
		CanResendAfter: canResendAfter,
	}, nil
}

func (s *signUpService) GetEmailVerificationStatus(ctx context.Context, userID uuid.UUID) (EmailVerificationStatus, error) {
	queries := store.New(s.db)

	user, err := queries.User_Get(ctx, uuidDomainToDb(userID))
	if err != nil {
		return EmailVerificationStatus{}, apperror.WrapUnexpectedDBError(err)
	}
	verification, err := queries.EmailVerification_Get(ctx, user.Email)
	if err != nil && err != store.ErrNoRows {
		return EmailVerificationStatus{}, apperror.WrapUnexpectedDBError(err)
	}

	status := EmailVerificationStatus{}

	if err != store.ErrNoRows {
		status.CanSendAgainAfter = Value(verification.CreatedAt.Time.Add(s.emailCodeResendAfter))
		status.WasSent = time.Now().Before(verification.ValidThrough.Time)
	}

	return status, nil
}
