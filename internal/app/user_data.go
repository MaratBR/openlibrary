package app

import (
	"context"
	"encoding/json"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/gofrs/uuid"
)

const UserDataMaxSize = 20 * 1024

var (
	UserDataErrors               = apperror.AppErrors.NewSubNamespace("user_data")
	ErrTypeUserDataKeyNotAllowed = UserDataErrors.NewType("key_not_allowed")
	ErrUserDataTooLarge          = UserDataErrors.NewType("size_limit").New("user data exceeds the 20KB size limit")
)

type GetUserDataQuery struct {
	UserID          uuid.UUID
	Key             string
	IgnoreWhitelist bool
}

type SetUserDataQuery struct {
	UserID          uuid.UUID
	Key             string
	Data            json.RawMessage
	IgnoreWhitelist bool
}

type UserDataService interface {
	Get(ctx context.Context, query GetUserDataQuery) (Nullable[json.RawMessage], error)
	Set(ctx context.Context, query SetUserDataQuery) error
}
