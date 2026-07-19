package app

import (
	"context"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
)

type readerPreferencesService struct {
	queries *store.Queries
}

func NewReaderPreferencesService(db DB) ReaderPreferencesService {
	return &readerPreferencesService{queries: store.New(db)}
}

func (s *readerPreferencesService) Get(ctx context.Context, userID uuid.UUID) (Nullable[ReaderPreferences], error) {
	row, err := s.queries.ReaderPreferences_Get(ctx, uuidDomainToDb(userID))
	if err == store.ErrNoRows {
		return Null[ReaderPreferences](), nil
	}
	if err != nil {
		return Null[ReaderPreferences](), apperror.WrapUnexpectedDBError(err)
	}
	return Value(ReaderPreferences{
		FontSize:   row.FontSize,
		FontFamily: row.FontFamily,
		PageColor:  row.PageColor,
		Theme:      row.Theme,
	}), nil
}

func (s *readerPreferencesService) Save(ctx context.Context, userID uuid.UUID, preferences ReaderPreferences) error {
	if err := preferences.Validate(); err != nil {
		return err
	}
	if err := s.queries.ReaderPreferences_Upsert(ctx, store.ReaderPreferences_UpsertParams{
		UserID:     uuidDomainToDb(userID),
		FontSize:   preferences.FontSize,
		FontFamily: preferences.FontFamily,
		PageColor:  preferences.PageColor,
		Theme:      preferences.Theme,
	}); err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}
	return nil
}
