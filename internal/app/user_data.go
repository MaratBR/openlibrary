package app

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/gofrs/uuid"
)

const UserDataMaxSize = 20 * 1024

var (
	ErrUserDataKeyNotAllowed = errors.New("user data key is not allowed")
	ErrUserDataTooLarge      = errors.New("user data exceeds the 20KB size limit")
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
