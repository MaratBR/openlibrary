package app

import (
	"github.com/joomcode/errorx"
)

var (
	ErrTraitForbidden          = errorx.RegisterTrait("forbidden")
	ErrTraitAuthorizationIssue = errorx.RegisterTrait("auth")
	ErrTraitEntityNotFound     = errorx.RegisterTrait("not_found")
	ErrTraitValidationError    = errorx.RegisterTrait("validation")

	AppErrors       = errorx.NewNamespace("app")
	UnexpectedError = AppErrors.NewType("unexpected")

	ErrGenericForbidden = AppErrors.NewType("generic_forbidden", ErrTraitForbidden).New("access to this resource is forbidden")
	ErrDbError          = AppErrors.NewType("db_error")

	ValidationError = AppErrors.NewType("validation", ErrTraitValidationError)
)

func wrapUnexpectedAppError(err error) error {
	return UnexpectedError.New(err.Error())
}

func wrapUnexpectedDBError(err error) error {
	return ErrDbError.Wrap(err, "unexpected database error")
}
