package app

import (
	"context"
	"errors"
)

var ErrSetupInvalidInput = errors.New("invalid setup input")

type SetupCommand struct {
	AdminUsername string
	AdminPassword string
	SiteName      string
}

type GetSetupDataResult struct {
	AdminUserExist bool
	IsSetup        bool
}

type SetupService interface {
	IsSetup() (bool, error)
	Setup(ctx context.Context, command SetupCommand) error
	GetSetupState() (GetSetupDataResult, error)
}
