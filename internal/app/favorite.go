package app

import (
	"context"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/gofrs/uuid"
)

type FavoriteService struct {
	db                     DB
	queries                *store.Queries
	favoritesRecalculation *FavoriteRecalculationBackgroundService
}

func NewFavoriteService(db DB, favoritesRecalculation *FavoriteRecalculationBackgroundService) *FavoriteService {
	return &FavoriteService{db: db, queries: store.New(db), favoritesRecalculation: favoritesRecalculation}
}

type SetFavoriteCommand struct {
	UserID     uuid.UUID
	BookID     int64
	IsFavorite bool
}

func (s *FavoriteService) SetFavorite(ctx context.Context, input SetFavoriteCommand) error {
	// TODO check if book exists
	err := s.queries.SetUserFavourite(ctx, store.SetUserFavouriteParams{
		UserID:     uuidDomainToDb(input.UserID),
		BookID:     input.BookID,
		IsFavorite: input.IsFavorite,
	})
	if err != nil {
		return err
	}
	s.favoritesRecalculation.ScheduleRecalculation(input.BookID)

	return nil
}
