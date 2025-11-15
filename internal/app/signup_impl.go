package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/MaratBR/openlibrary/internal/app/email"
	"github.com/MaratBR/openlibrary/internal/commonutil"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
	"github.com/knadh/koanf/v2"
)

type signUpService struct {
	db           DB
	cfg          *koanf.Koanf
	siteConfig   *SiteConfig
	emailService email.Service
}

func NewSignUpService(db DB, cfg *koanf.Koanf, siteConfig *SiteConfig, emailService email.Service) SignUpService {
	return &signUpService{
		db:           db,
		cfg:          cfg,
		siteConfig:   siteConfig,
		emailService: emailService,
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
	if s.cfg.Bool("auth.requireEmail") && input.Email == "" {
		return SignUpResult{}, SignUpInvalidInput.New("email is required")
	}
	if err := ValidateUserName(input.Username); err != nil {
		return SignUpResult{}, SignUpInvalidInput.Wrap(err, "invalid username")
	}
	passwordRequirements := s.siteConfig.Get().PasswordRequirements
	if err := ValidatePassword(input.Password, passwordRequirements); err != nil {
		return SignUpResult{}, SignUpInvalidInput.Wrap(err, "invalid password")
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
		userWithSameEmailExists, err := queries.UserExistsByEmail(ctx, input.Username)
		if err != nil {
			rollbackTx(ctx, tx)
			return SignUpResult{}, wrapUnexpectedDBError(err)
		}

		if userWithSameEmailExists {
			return SignUpResult{EmailTaken: true}, nil
		}
	}

	userWithSameNameExists, err := queries.UserExistsByUsername(ctx, input.Username)
	if err != nil {
		rollbackTx(ctx, tx)
		return SignUpResult{}, wrapUnexpectedDBError(err)
	}

	if userWithSameNameExists {
		rollbackTx(ctx, tx)
		return SignUpResult{}, ErrUsernameTaken
	}

	userID, err := createUser(ctx, queries, input.Username, input.Password, RoleUser)
	if err != nil {
		rollbackTx(ctx, tx)
		return SignUpResult{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return SignUpResult{}, err
	}

	if emailVerificationRequired {
		err := s.verifyEmailRequest(ctx, input.Email, userID)
		if err != nil {
			return SignUpResult{}, err
		}
	}

	return SignUpResult{CreatedUserID: userID, EmailVerificationRequired: emailVerificationRequired, Created: true}, nil
}

func (s *signUpService) verifyEmailRequest(ctx context.Context, email string, userID uuid.UUID) error {
	hash, code, err := newEmailVerificationCode()
	if err != nil {
		return err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return wrapUnexpectedDBError(err)
	}
	queries := store.New(s.db).WithTx(tx)
	err = queries.EmailVerification_Delete(ctx, email)
	if err != nil {
		rollbackTx(ctx, tx)
		return wrapUnexpectedDBError(err)
	}
	err = queries.EmailVerification_Insert(ctx, store.EmailVerification_InsertParams{
		UserID:               uuidDomainToDb(userID),
		Email:                email,
		ValidThrough:         timeToTimestamptz(time.Now().Add(time.Hour * 6)),
		VerificationCodeHash: hash,
	})
	if err != nil {
		rollbackTx(ctx, tx)
		return wrapUnexpectedDBError(err)
	}

	// TODO some good looking message or something
	err = s.emailService.Send(ctx, email, "OL verification code", fmt.Sprintf("Your verification code is %s\n\nIf you did not sign up, please just ignore this message and apologies for disturbance :)", code))
	if err != nil {
		rollbackTx(ctx, tx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return wrapUnexpectedDBError(err)
	}

	return nil
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
