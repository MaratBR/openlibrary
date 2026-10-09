package app

import (
	"context"

	"github.com/MaratBR/openlibrary/internal/app/apperror"
	"github.com/MaratBR/openlibrary/internal/store"
	"uuid"
)

type readerPreferencesService struct {
	queries *store.Queries
	fonts   ReaderFontService
}

func NewReaderPreferencesService(db DB, fonts ReaderFontService) ReaderPreferencesService {
	return &readerPreferencesService{queries: store.New(db), fonts: fonts}
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
		ContentWidth: row.ContentWidth,
		FontSize:     row.FontSize,
		FontFamily:   row.FontFamily,
		PageColor:    row.PageColor,
		Theme:        row.Theme,
	}), nil
}

func (s *readerPreferencesService) Save(ctx context.Context, userID uuid.UUID, preferences ReaderPreferences) error {
	if err := preferences.ValidateWithFonts(s.fonts.List()); err != nil {
		return err
	}
	if preferences.ContentWidth == 0 {
		preferences.ContentWidth = 72
	}
	if err := s.queries.ReaderPreferences_Upsert(ctx, store.ReaderPreferences_UpsertParams{
		UserID:       uuidDomainToDb(userID),
		ContentWidth: preferences.ContentWidth,
		FontSize:     preferences.FontSize,
		FontFamily:   preferences.FontFamily,
		PageColor:    preferences.PageColor,
		Theme:        preferences.Theme,
	}); err != nil {
		return apperror.WrapUnexpectedDBError(err)
	}
	return nil
}
