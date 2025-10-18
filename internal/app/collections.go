package app

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
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
	BooksCount    int
	LastUpdatedAt Nullable[time.Time]
	UserID        uuid.UUID
}

type Collection2Dto struct {
	ID            int64
	Name          string
	BooksCount    int
	LastUpdatedAt Nullable[time.Time]
	UserID        uuid.UUID
	UserName      string
}

type CreateCollectionCommand struct {
	Name   string
	UserID uuid.UUID
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
	Collection Collection2Dto
	Page       int32
	TotalPages int32
	PageSize   int32
}

type CollectionService interface {
	GetUserCollections(ctx context.Context, query GetUserCollectionsQuery) (GetUserCollectionsResult, error)
	GetRecentUserCollections(ctx context.Context, query GetRecentCollectionsQuery) ([]CollectionDto, error)
	CreateCollection(ctx context.Context, cmd CreateCollectionCommand) (int64, error)
	AddToCollections(ctx context.Context, cmd AddToCollectionsCommand) error
	GetBookCollections(ctx context.Context, query GetBookCollectionsQuery) ([]CollectionDto, error)

	GetCollectionBooks(ctx context.Context, query GetCollectionBooksQuery) (GetCollectionBooksResult, error)
	GetCollectionBooksMap(ctx context.Context, collections []CollectionDto) (map[int64][]CollectionBookDto, error)
}
