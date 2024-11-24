package app

import (
	"github.com/joomcode/errorx"
)

var (
	ErrTraitForbidden          = errorx.RegisterTrait("forbidden")
	ErrTraitAuthorizationIssue = errorx.RegisterTrait("auth")
	ErrTraitEntityNotFound     = errorx.RegisterTrait("not_found")

	AppErrors = errorx.NewNamespace("app")

	ErrGenericForbidden = AppErrors.NewType("generic_forbidden", ErrTraitForbidden).New("access to this resource is forbidden")
	ErrDbError          = AppErrors.NewType("db_error")
)

func wrapUnexpectedDBError(err error) error {
	return ErrDbError.Wrap(err, "unexpected database error")
}
