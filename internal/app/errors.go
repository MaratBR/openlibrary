package app

import (
	"github.com/joomcode/errorx"
)

var (
	ErrTraitForbidden = errorx.RegisterTrait("forbidden")
	AppErrors         = errorx.NewNamespace("app")

	ErrGenericForbidden = AppErrors.NewType("generic_forbidden", ErrTraitForbidden).New("access to this resource is forbidden")
)
