package apperror

import "github.com/joomcode/errorx"

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

func WrapUnexpectedAppError(err error) error {
	return UnexpectedError.New(err.Error())
}

func WrapUnexpectedDBError(err error) error {
	return ErrDbError.Wrap(err, "unexpected database error")
}

func IsNotFoundError(err error) bool {
	return errorx.HasTrait(err, ErrTraitEntityNotFound)
}
