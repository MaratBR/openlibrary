package app

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/MaratBR/openlibrary/internal/app/dal"
	"github.com/MaratBR/openlibrary/internal/app/siteconfig"
	"github.com/MaratBR/openlibrary/internal/store"
)

type setupService struct {
	db     DB
	config siteconfig.Store
	mu     sync.Mutex
	cached *bool
}

func NewSetupService(db DB, config siteconfig.Store) SetupService {
	return &setupService{db: db, config: config}
}

func (s *setupService) getSetupState(ctx context.Context) (GetSetupDataResult, error) {
	users, err := store.New(s.db).User_GetSetupState(ctx)
	if err != nil {
		return GetSetupDataResult{}, fmt.Errorf("get setup state: %w", err)
	}
	return GetSetupDataResult{
		AdminUserExist: users,
		IsSetup:        users && s.config.Get().SiteSetup.IsSetup,
	}, nil
}

// GetSetupState always checks current state, independently of the IsSetup cache.
func (s *setupService) GetSetupState() (GetSetupDataResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.getSetupState(context.Background())
}

func (s *setupService) IsSetup() (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cached != nil {
		return *s.cached, nil
	}
	state, err := s.getSetupState(context.Background())
	if err != nil {
		return false, err
	}
	s.cached = &state.IsSetup
	return state.IsSetup, nil
}

func (s *setupService) Setup(ctx context.Context, command SetupCommand) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	state, err := s.getSetupState(ctx)
	if err != nil {
		return err
	}
	if state.IsSetup {
		return errors.New("site is already set up")
	}
	name := strings.TrimSpace(command.SiteName)
	if name == "" {
		return fmt.Errorf("%w: site name is required", ErrSetupInvalidInput)
	}
	if !state.AdminUserExist {
		if err := ValidateUserName(command.AdminUsername); err != nil {
			return fmt.Errorf("%w: %v", ErrSetupInvalidInput, err)
		}
		if command.AdminPassword == "" {
			return fmt.Errorf("%w: admin password is required", ErrSetupInvalidInput)
		}
		if err := siteconfig.ValidatePassword(command.AdminPassword, s.config.Get().PasswordRequirements); err != nil {
			return fmt.Errorf("%w: %v", ErrSetupInvalidInput, err)
		}
		tx, err := s.db.Begin(ctx)
		if err != nil {
			return err
		}
		defer dal.RollbackTx(ctx, tx)
		queries := store.New(tx)
		exists, err := queries.User_ExistsByUsername(ctx, command.AdminUsername)
		if err != nil {
			return err
		}
		if exists {
			return ErrUsernameTaken
		}
		if _, err := createUser(ctx, queries, command.AdminUsername, "", command.AdminPassword, RoleAdmin, true); err != nil {
			return err
		}
		if err := tx.Commit(ctx); err != nil {
			return err
		}
	}
	// A failed config save can be retried without creating another administrator.
	s.cached = nil
	previous := s.config.Get().SiteSetup
	s.config.Get().SiteSetup = siteconfig.SiteSetup{IsSetup: true, SiteName: name}
	if err := s.config.Save(ctx, false); err != nil {
		s.config.Get().SiteSetup = previous
		return err
	}
	return nil
}
