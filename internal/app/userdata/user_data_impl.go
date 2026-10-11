package userdata

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/app/dal"
	"github.com/MaratBR/openlibrary/internal/app/value"
	"github.com/MaratBR/openlibrary/internal/store"
)

var userDataKeyWhitelist = map[string]struct{}{
	"bm:fonts": {},
	"bm:tools": {},
}

type userDataService struct {
	queries *store.Queries
}

func NewUserDataService(db store.DBTX) UserDataService {
	return &userDataService{queries: store.New(db)}
}

func (s *userDataService) Get(ctx context.Context, query GetUserDataQuery) (value.Nullable[json.RawMessage], error) {
	if err := validateUserDataKey(query.Key, query.IgnoreWhitelist); err != nil {
		return value.Null[json.RawMessage](), err
	}

	data, err := s.queries.UserData_Get(ctx, store.UserData_GetParams{
		UserID: dal.UUIDToDB(query.UserID),
		Key:    query.Key,
	})
	if err == store.ErrNoRows {
		return value.Null[json.RawMessage](), nil
	}
	if err != nil {
		return value.Null[json.RawMessage](), apperror.WrapUnexpectedDBError(err)
	}
	return value.Value(json.RawMessage(data)), nil
}

func (s *userDataService) Set(ctx context.Context, query SetUserDataQuery) error {
	if err := validateUserDataKey(query.Key, query.IgnoreWhitelist); err != nil {
		return err
	}
	if len(query.Data) > UserDataMaxSize {
		return ErrUserDataTooLarge
	}
	if !json.Valid(query.Data) {
		return errors.New("user data must be valid JSON")
	}

	if err := s.queries.UserData_Set(ctx, store.UserData_SetParams{
		UserID: dal.UUIDToDB(query.UserID),
		Key:    query.Key,
		Data:   query.Data,
	}); err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}
	return nil
}

func validateUserDataKey(key string, ignoreWhitelist bool) error {
	if ignoreWhitelist {
		return nil
	}
	if _, ok := userDataKeyWhitelist[key]; !ok {
		return ErrUserDataKeyNotAllowed
	}
	return nil
}
