package app

import (
	"context"

	"github.com/MaratBR/openlibrary/internal/store"
)

type signUpService struct {
	db DB
}

func NewSignUpService(db DB) SignUpService {
	return &signUpService{
		db: db,
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

	return SignUpResult{CreatedUserID: userID, Created: true}, nil
}
