package app

import (
	"context"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/MaratBR/openlibrary/lib/gset"
	"github.com/gofrs/uuid"
)

type collectionService struct {
	db      DB
	queries *store.Queries
}

func collectionBelongsTo(col *store.Collection, userID uuid.UUID) bool {
	return uuidDbToDomain(col.UserID) == userID
}

func storeCollectionToCollectionDto(row store.Collection) CollectionDto {
	return CollectionDto{
		ID:            row.ID,
		Name:          row.Name,
		LastUpdatedAt: timeNullableDbToDomain(row.LastUpdatedAt),
	}
}

// GetRecentUserCollections implements CollectionService.
func (c *collectionService) GetRecentUserCollections(ctx context.Context, query GetRecentCollections) ([]CollectionDto, error) {
	rows, err := c.queries.GetLatestUserCollections(ctx, store.GetLatestUserCollectionsParams{
		UserID: uuidDomainToDb(query.UserID),
		Limit:  query.Limit,
	})
	if err != nil {
		return nil, wrapUnexpectedDBError(err)
	}

	return MapSlice(rows, storeCollectionToCollectionDto), nil
}

func (c *collectionService) AddToCollections(ctx context.Context, cmd AddToCollectionsCommand) error {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return wrapUnexpectedDBError(err)
	}

	queries := c.queries.WithTx(tx)

	// get all collections, then figure out which ones have to be added or removed
	collections, err := queries.GetCollections(ctx, cmd.CollectionID)
	if err != nil {
		rollbackTx(ctx, tx)
		return wrapUnexpectedDBError(err)
	}
	bookCollections, err := queries.GetBookCollections(ctx, store.GetBookCollectionsParams{
		UserID: uuidDomainToDb(cmd.ActorUserID),
		BookID: cmd.BookID,
	})
	if err != nil {
		rollbackTx(ctx, tx)
		return wrapUnexpectedDBError(err)
	}

	addCollections := make([]int64, 0)
	removeCollections := make([]int64, 0)
	{
		oldCollections := gset.New[int64]()
		for _, col := range bookCollections {
			oldCollections.Add(col.ID)
		}
		newCollections := gset.New[int64]()
		for _, col := range collections {
			if !collectionBelongsTo(&col, cmd.ActorUserID) {
				continue
			}

			newCollections.Add(col.ID)
			if !oldCollections.Contains(col.ID) {
				// this collection must be added
				addCollections = append(addCollections, col.ID)
			}
		}

		for _, col := range bookCollections {
			if !newCollections.Contains(col.ID) {
				removeCollections = append(removeCollections, col.ID)
			}
		}
	}

	for _, colID := range addCollections {
		err = queries.AddBookToCollection(ctx, store.AddBookToCollectionParams{
			BookID:       cmd.BookID,
			CollectionID: colID,
		})
		if err != nil {
			rollbackTx(ctx, tx)
			return wrapUnexpectedDBError(err)
		}
	}

	for _, colID := range removeCollections {
		err = queries.DeleteBookFromCollection(ctx, store.DeleteBookFromCollectionParams{
			BookID:       cmd.BookID,
			CollectionID: colID,
		})
		if err != nil {
			rollbackTx(ctx, tx)
			return wrapUnexpectedDBError(err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return wrapUnexpectedDBError(err)
	}

	return nil
}

func (c *collectionService) CreateCollection(ctx context.Context, cmd CreateCollectionCommand) (int64, error) {
	id := GenID()
	err := c.queries.InsertCollection(ctx, store.InsertCollectionParams{
		ID:     id,
		Name:   cmd.Name,
		UserID: uuidDomainToDb(cmd.UserID),
	})
	if err != nil {
		return 0, err
	} else {
		return id, nil
	}
}

func (c *collectionService) GetBookCollections(ctx context.Context, query GetBookCollectionsQuery) ([]CollectionDto, error) {
	rows, err := c.queries.GetBookCollections(ctx, store.GetBookCollectionsParams{
		UserID: uuidDomainToDb(query.ActorUserID),
		BookID: query.BookID,
	})
	if err != nil {
		return nil, wrapUnexpectedDBError(err)
	}

	return MapSlice(rows, storeCollectionToCollectionDto), nil

}

func NewCollectionsService(db DB) CollectionService {
	return &collectionService{
		db:      db,
		queries: store.New(db),
	}
}
