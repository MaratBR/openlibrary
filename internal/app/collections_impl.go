package app

import (
	"context"
	"log/slog"
	"math"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/MaratBR/openlibrary/lib/gset"
	"github.com/gofrs/uuid"
)

type collectionService struct {
	db            DB
	queries       *store.Queries
	uploadService *UploadService
	tagsService   TagsService
}

func collectionBelongsTo(col *store.Collection, userID uuid.UUID) bool {
	return uuidDbToDomain(col.UserID) == userID
}

func storeCollectionToCollectionDto(row store.Collection) CollectionDto {
	return CollectionDto{
		ID:            row.ID,
		Name:          row.Name,
		LastUpdatedAt: timeNullableDbToDomain(row.LastUpdatedAt),
		BooksCount:    int(row.BooksCount),
	}
}

func (c *collectionService) GetUserCollections(ctx context.Context, query GetUserCollectionsQuery) (GetUserCollectionsResult, error) {
	rows, err := c.queries.Collection_GetByUser(ctx, store.Collection_GetByUserParams{
		Offset: query.PageSize * (query.Page - 1),
		Limit:  query.PageSize,
		UserID: uuidDomainToDb(query.UserID),
	})
	if err != nil {
		return GetUserCollectionsResult{}, wrapUnexpectedDBError(err)
	}

	collections := MapSlice(rows, storeCollectionToCollectionDto)

	count, err := c.queries.Collections_CountByUser(ctx, uuidDomainToDb(query.UserID))
	if err != nil {
		return GetUserCollectionsResult{}, wrapUnexpectedDBError(err)
	}

	return GetUserCollectionsResult{
		Collections: collections,
		TotalPages:  int32(math.Ceil(float64(count) / float64(query.PageSize))),
		Page:        query.Page,
	}, nil
}

// GetRecentUserCollections implements CollectionService.
func (c *collectionService) GetRecentUserCollections(ctx context.Context, query GetRecentCollectionsQuery) ([]CollectionDto, error) {
	rows, err := c.queries.Collection_GetRecentByUser(ctx, store.Collection_GetRecentByUserParams{
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
	collections, err := queries.Collections_ListByID(ctx, cmd.CollectionID)
	if err != nil {
		rollbackTx(ctx, tx)
		return wrapUnexpectedDBError(err)
	}
	bookCollections, err := queries.Collection_GetByBook(ctx, store.Collection_GetByBookParams{
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
		err = queries.Collection_AddBookToCollection(ctx, store.Collection_AddBookToCollectionParams{
			BookID:       cmd.BookID,
			CollectionID: colID,
		})
		if err != nil {
			rollbackTx(ctx, tx)
			return wrapUnexpectedDBError(err)
		}
	}

	for _, colID := range removeCollections {
		err = queries.Collection_DeleteBookFromCollection(ctx, store.Collection_DeleteBookFromCollectionParams{
			BookID:       cmd.BookID,
			CollectionID: colID,
		})
		if err != nil {
			rollbackTx(ctx, tx)
			return wrapUnexpectedDBError(err)
		}

	}

	// recalculate collection counter
	for _, colID := range removeCollections {
		err = queries.Collection_RecalculateCounter(ctx, colID)
		if err != nil {
			rollbackTx(ctx, tx)
			return wrapUnexpectedDBError(err)
		}
	}

	for _, colID := range addCollections {
		err = queries.Collection_RecalculateCounter(ctx, colID)
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
	err := c.queries.Collection_Insert(ctx, store.Collection_InsertParams{
		ID:     id,
		Name:   cmd.Name,
		Slug:   makeSlug(cmd.Name),
		UserID: uuidDomainToDb(cmd.UserID),
	})
	if err != nil {
		return 0, err
	} else {
		return id, nil
	}
}

func (c *collectionService) GetBookCollections(ctx context.Context, query GetBookCollectionsQuery) ([]CollectionDto, error) {
	rows, err := c.queries.Collection_GetByBook(ctx, store.Collection_GetByBookParams{
		UserID: uuidDomainToDb(query.ActorUserID),
		BookID: query.BookID,
	})
	if err != nil {
		return nil, wrapUnexpectedDBError(err)
	}

	return MapSlice(rows, storeCollectionToCollectionDto), nil

}

func (c *collectionService) GetCollectionBooks(ctx context.Context, query GetCollectionBooksQuery) (GetCollectionBooksResult, error) {
	collection, err := c.queries.Collection_Get(ctx, query.CollectionID)
	if err != nil {
		if err == store.ErrNoRows {
			return GetCollectionBooksResult{}, ErrCollectionNotExists
		}
		return GetCollectionBooksResult{}, wrapUnexpectedDBError(err)
	}

	rows, err := c.queries.Collection_GetBooks(ctx, store.Collection_GetBooksParams{
		Limit:        query.PageSize,
		Offset:       (query.Page - 1) * query.PageSize,
		CollectionID: query.CollectionID,
	})
	if err != nil {
		return GetCollectionBooksResult{}, wrapUnexpectedDBError(err)
	}

	books := make([]CollectionBook2Dto, 0, len(rows))
	tagsAgg := newTagsAggregator(c.tagsService)

	for _, row := range rows {
		tagsAgg.Add(row.ID, row.TagIds)

		book := CollectionBook2Dto{
			ID:         row.ID,
			Name:       row.Name,
			Slug:       row.Slug,
			Summary:    row.Summary,
			Cover:      getBookCoverURL(c.uploadService, row.ID, row.HasCover),
			AuthorID:   uuidDbToDomain(row.AuthorUserID),
			AuthorName: row.AuthorName,
			Tags:       nil,
		}

		books = append(books, book)
	}

	tags, err := tagsAgg.Fetch(ctx)
	if err != nil {
		return GetCollectionBooksResult{}, err
	}

	for i, row := range rows {
		tagsList := make([]DefinedTagDto, 0, len(row.TagIds))
		for _, tagId := range row.TagIds {
			tag, ok := tags[tagId]
			if ok {
				tagsList = append(tagsList, tag)
			}
		}
		books[i].Tags = tagsList
	}

	return GetCollectionBooksResult{
		Books: books,
		Collection: Collection2Dto{
			ID:            collection.ID,
			Name:          collection.Name,
			BooksCount:    int(collection.BooksCount),
			LastUpdatedAt: timeNullableDbToDomain(collection.LastUpdatedAt),
			UserID:        uuidDbToDomain(collection.UserID),
			UserName:      collection.UserName,
		},
		Page:       query.Page,
		PageSize:   query.PageSize,
		TotalPages: int32(math.Ceil(float64(collection.BooksCount) / float64(query.PageSize))),
	}, nil
}

func (c *collectionService) getCollectionBooksList(ctx context.Context, collectionID int64, page, pageSize int32) ([]CollectionBookDto, error) {
	rows, err := c.queries.Collection_GetBooks(ctx, store.Collection_GetBooksParams{
		CollectionID: collectionID,
		Limit:        page,
		Offset:       pageSize * (page - 1),
	})
	if err != nil {
		return nil, err
	}

	books := MapSlice(rows, func(row store.Collection_GetBooksRow) CollectionBookDto {
		return CollectionBookDto{
			ID:    row.ID,
			Name:  row.Name,
			Cover: getBookCoverURL(c.uploadService, row.ID, row.HasCover),
		}
	})

	return books, nil
}

func (c *collectionService) GetCollectionBooksMap(ctx context.Context, collections []CollectionDto) (map[int64][]CollectionBookDto, error) {
	m := make(map[int64][]CollectionBookDto, len(collections))

	for _, col := range collections {
		books, err := c.getCollectionBooksList(ctx, col.ID, 1, 6)
		if err != nil {
			slog.Error("failed to get collection books", "collectionID", col.ID, "err", err)
			continue
		}
		m[col.ID] = books
	}

	return m, nil
}

func NewCollectionsService(db DB, tagsService TagsService, uploadService *UploadService) CollectionService {
	return &collectionService{
		db:            db,
		queries:       store.New(db),
		uploadService: uploadService,
		tagsService:   tagsService,
	}
}
