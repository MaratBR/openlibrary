package app

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	CollectionErrors       = AppErrors.NewSubNamespace("collection")
	ErrCollectionNotExists = CollectionErrors.NewType("404", ErrTraitEntityNotFound).New("collection not found")
)

type GetUserCollectionsQuery struct {
	UserID   uuid.UUID
	Page     int32
	PageSize int32
}

type GetUserCollectionsResult struct {
	Collections []CollectionDto
	TotalPages  int32
	Page        int32
}

type GetRecentCollectionsQuery struct {
	UserID uuid.UUID
	Limit  int32
}

type CollectionDto struct {
	ID            int64
	Name          string
	Slug          string
	BooksCount    int
	LastUpdatedAt Nullable[time.Time]
	UserID        uuid.UUID
	UserName      string
	IsPublic      bool
	Summary       string
}

func newCollectionDto(
	id int64,
	name string,
	booksCount int,
	lastUpdatedAt pgtype.Timestamptz,
	userID pgtype.UUID,
	userName string,
	isPublic bool,
	summary string,
	slug string,
) CollectionDto {
	return CollectionDto{
		ID:            id,
		Name:          name,
		BooksCount:    booksCount,
		LastUpdatedAt: timeNullableDbToDomain(lastUpdatedAt),
		UserID:        uuidDbToDomain(userID),
		UserName:      userName,
		IsPublic:      isPublic,
		Summary:       summary,
		Slug:          slug,
	}
}

type CreateCollectionCommand struct {
	Name        string
	Description string
	UserID      uuid.UUID
}

type UpdateCollectionCommand struct {
	ID          int64
	Name        string
	Public      bool
	Summary     string
	ActorUserID uuid.UUID
}

type AddToCollectionsCommand struct {
	ActorUserID  uuid.UUID
	CollectionID []int64
	BookID       int64
}

type GetBookCollectionsQuery struct {
	ActorUserID uuid.UUID
	BookID      int64
}

type CollectionBookDto struct {
	ID    int64
	Name  string
	Cover string
}

type CollectionBook2Dto struct {
	ID         int64
	Name       string
	Slug       string
	Summary    string
	AuthorName string
	AuthorID   uuid.UUID
	Cover      string
	Tags       []DefinedTagDto
}

type GetCollectionBooksQuery struct {
	CollectionID int64
	Page         int32
	PageSize     int32
}

type GetCollectionBooksResult struct {
	Books      []CollectionBook2Dto
	Page       int32
	TotalPages int32
	PageSize   int32
}

type RemoveFromCollectionCommand struct {
	CollectionID int64
	BookID       int64
	UserID       uuid.UUID
}

type DeleteCollectionCommand struct {
	ActorUserID  uuid.UUID
	CollectionID int64
}

type CollectionService interface {
	GetUserCollections(ctx context.Context, query GetUserCollectionsQuery) (GetUserCollectionsResult, error)
	GetRecentUserCollections(ctx context.Context, query GetRecentCollectionsQuery) ([]CollectionDto, error)
	CreateCollection(ctx context.Context, cmd CreateCollectionCommand) (int64, error)
	UpdateCollection(ctx context.Context, cmd UpdateCollectionCommand) error
	AddToCollections(ctx context.Context, cmd AddToCollectionsCommand) error
	RemoveFromCollection(ctx context.Context, cmd RemoveFromCollectionCommand) error
	GetBookCollections(ctx context.Context, query GetBookCollectionsQuery) ([]CollectionDto, error)

	GetCollectionBooks(ctx context.Context, query GetCollectionBooksQuery) (GetCollectionBooksResult, error)
	GetCollectionBooksMap(ctx context.Context, collections []CollectionDto) (map[int64][]CollectionBookDto, error)
	GetCollection(ctx context.Context, id int64) (Nullable[CollectionDto], error)
	DeleteCollection(ctx context.Context, cmd DeleteCollectionCommand) error
}
